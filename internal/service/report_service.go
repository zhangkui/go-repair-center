package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-repair-center/internal/repository/mysql"
)

type ReportService struct {
	analytics *mysql.AnalyticsRepo
	cache     *CacheService
}

type RevenueReport struct {
	PeriodStart       time.Time                `json:"period_start"`
	PeriodEnd         time.Time                `json:"period_end"`
	Revenue           map[string]float64       `json:"revenue"`
	RepairTrend       []mysql.DailyCount       `json:"repair_trend"`
	Satisfaction      []mysql.SatisfactionStat `json:"satisfaction"`
	TechnicianRanking []mysql.TechnicianStat   `json:"technician_ranking"`
}

func NewReportService(analytics *mysql.AnalyticsRepo, cache *CacheService) *ReportService {
	return &ReportService{analytics: analytics, cache: cache}
}

func (s *ReportService) RevenueAndPerformance(ctx context.Context, days int) (*RevenueReport, error) {
	if days <= 0 {
		days = 30
	}
	cacheKey := fmt.Sprintf("report:revenue:%d", days)
	result := &RevenueReport{}
	if s.cache != nil {
		if ok, err := s.cache.GetJSON(ctx, cacheKey, result); err == nil && ok {
			return result, nil
		}
	}

	start := time.Now().AddDate(0, 0, -days)
	revenue, err := s.analytics.RevenueSummary(ctx, start)
	if err != nil {
		return nil, err
	}
	trend, err := s.analytics.DailyRepairOrders(ctx, days)
	if err != nil {
		return nil, err
	}
	satisfaction, err := s.analytics.SatisfactionByDimension(ctx)
	if err != nil {
		return nil, err
	}
	ranking, err := s.analytics.TechnicianPerformance(ctx, start, 20)
	if err != nil {
		return nil, err
	}

	result = &RevenueReport{
		PeriodStart:       start,
		PeriodEnd:         time.Now(),
		Revenue:           revenue,
		RepairTrend:       trend,
		Satisfaction:      satisfaction,
		TechnicianRanking: ranking,
	}
	NormalizeReportSlices(result)
	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, result, 5*time.Minute)
	}
	return result, nil
}

func (s *ReportService) RevenueCSV(ctx context.Context, days int) (string, error) {
	report, err := s.RevenueAndPerformance(ctx, days)
	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}
	writer := csv.NewWriter(builder)

	rows := [][]string{
		{"period_start", report.PeriodStart.Format(time.RFC3339)},
		{"period_end", report.PeriodEnd.Format(time.RFC3339)},
		{"total_amount", fmt.Sprintf("%.2f", report.Revenue["total_amount"])},
		{"period_amount", fmt.Sprintf("%.2f", report.Revenue["period_amount"])},
		{},
		{"repair_day", "repair_count"},
	}
	for _, item := range report.RepairTrend {
		rows = append(rows, []string{item.Day, strconv.FormatInt(item.Count, 10)})
	}
	rows = append(rows, []string{}, []string{"dimension", "avg_score"})
	for _, item := range report.Satisfaction {
		rows = append(rows, []string{item.Dimension, fmt.Sprintf("%.2f", item.Average)})
	}
	rows = append(rows, []string{}, []string{"technician_id", "username", "order_count", "rework_count", "avg_score"})
	for _, item := range report.TechnicianRanking {
		rows = append(rows, []string{
			strconv.FormatInt(item.TechnicianID, 10),
			item.Username,
			strconv.FormatInt(item.OrderCount, 10),
			strconv.FormatInt(item.ReworkCount, 10),
			fmt.Sprintf("%.2f", item.AvgScore),
		})
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}
