package service

import (
	"context"
	"time"

	"go-repair-center/internal/repository/mysql"
)

type DashboardService struct {
	analytics *mysql.AnalyticsRepo
	cache     *CacheService
}

type DashboardSummary struct {
	TodayRepairCount      int64                       `json:"today_repair_count"`
	InProgressCount       int64                       `json:"in_progress_count"`
	WaitingPickupCount    int64                       `json:"waiting_pickup_count"`
	PendingFeedbackCount  int64                       `json:"pending_feedback_count"`
	CompletedTodayCount   int64                       `json:"completed_today_count"`
	CompletionRate        float64                     `json:"completion_rate"`
	RepairTrend           []mysql.DailyCount          `json:"repair_trend"`
	Satisfaction          []mysql.SatisfactionStat    `json:"satisfaction"`
	TechnicianPerformance []mysql.TechnicianStat      `json:"technician_performance"`
	Revenue               map[string]float64          `json:"revenue"`
	StatusCatalogue       map[string][]string         `json:"status_catalogue"`
}

func NewDashboardService(analytics *mysql.AnalyticsRepo, cache *CacheService) *DashboardService {
	return &DashboardService{analytics: analytics, cache: cache}
}

func (s *DashboardService) Summary(ctx context.Context) (*DashboardSummary, error) {
	cacheKey := "dashboard:summary"
	result := &DashboardSummary{}
	if s.cache != nil {
		if ok, err := s.cache.GetJSON(ctx, cacheKey, result); err == nil && ok {
			return result, nil
		}
	}

	todayRepairCount, err := s.analytics.CreatedToday(ctx)
	if err != nil {
		return nil, err
	}
	inProgressCount, err := s.analytics.InProgressExecutionCount(ctx)
	if err != nil {
		return nil, err
	}
	waitingPickupCount, err := s.analytics.WaitingPickupCount(ctx)
	if err != nil {
		return nil, err
	}
	pendingFeedbackCount, err := s.analytics.PendingFeedbackCount(ctx)
	if err != nil {
		return nil, err
	}
	completedTodayCount, err := s.analytics.CompletedToday(ctx)
	if err != nil {
		return nil, err
	}
	repairTrend, err := s.analytics.DailyRepairOrders(ctx, 14)
	if err != nil {
		return nil, err
	}
	satisfaction, err := s.analytics.SatisfactionByDimension(ctx)
	if err != nil {
		return nil, err
	}
	techStats, err := s.analytics.TechnicianPerformance(ctx, time.Now().AddDate(0, 0, -30), 10)
	if err != nil {
		return nil, err
	}
	revenue, err := s.analytics.RevenueSummary(ctx, time.Now().AddDate(0, -1, 0))
	if err != nil {
		return nil, err
	}

	completionRate := 0.0
	if todayRepairCount > 0 {
		completionRate = float64(completedTodayCount) / float64(todayRepairCount)
	}

	result = &DashboardSummary{
		TodayRepairCount:      todayRepairCount,
		InProgressCount:       inProgressCount,
		WaitingPickupCount:    waitingPickupCount,
		PendingFeedbackCount:  pendingFeedbackCount,
		CompletedTodayCount:   completedTodayCount,
		CompletionRate:        completionRate,
		RepairTrend:           repairTrend,
		Satisfaction:          satisfaction,
		TechnicianPerformance: techStats,
		Revenue:               revenue,
		StatusCatalogue:       statusCatalogue(),
	}
	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, result, 2*time.Minute)
	}
	return result, nil
}

func statusCatalogue() map[string][]string {
	return map[string][]string{
		"device":           {"NORMAL", "REPORTED", "REPAIRING", "READY_FOR_PICKUP", "PICKED_UP", "SCRAPPED"},
		"repair_order":     {"PENDING", "DISPATCHED", "WAITING_VISIT", "WAITING_DELIVERY", "ACCEPTED", "COMPLETED", "CANCELLED"},
		"quotation":        {"DRAFT", "WAITING_CLIENT", "CONFIRMED", "REJECTED", "EXPIRED", "CONVERTED"},
		"repair_execution": {"PENDING", "IN_PROGRESS", "COMPLETED", "WAITING_TEST", "TESTED", "WAITING_PICKUP", "DELIVERED"},
		"warranty":         {"ACTIVE", "EXPIRED", "VOID"},
		"feedback":         {"PENDING", "VISITED", "CLOSED"},
	}
}
