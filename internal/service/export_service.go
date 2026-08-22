package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"go-repair-center/internal/repository"
)

type ExportService struct{}

type ExportFile struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content"`
}

func NewExportService() *ExportService {
	return &ExportService{}
}

func (s *ExportService) ExportCSV(_ context.Context, fileName string, records []repository.Record, columns []string) (*ExportFile, error) {
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	if len(columns) == 0 {
		columns = inferColumns(records)
	}
	if err := writer.Write(columns); err != nil {
		return nil, err
	}
	for _, record := range records {
		row := make([]string, 0, len(columns))
		for _, column := range columns {
			row = append(row, stringify(record[column]))
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return &ExportFile{
		FileName:    ensureExtension(fileName, ".csv"),
		ContentType: "text/csv; charset=utf-8",
		Content:     buffer.Bytes(),
	}, nil
}

func (s *ExportService) ExportJSON(_ context.Context, fileName string, records []repository.Record) (*ExportFile, error) {
	payload, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return nil, err
	}
	return &ExportFile{
		FileName:    ensureExtension(fileName, ".json"),
		ContentType: "application/json; charset=utf-8",
		Content:     payload,
	}, nil
}

func (s *ExportService) ExportRepairOrders(ctx context.Context, records []repository.Record) (*ExportFile, error) {
	columns := []string{
		"id",
		"order_number",
		"customer_id",
		"device_id",
		"fault_description",
		"urgency",
		"service_method",
		"status",
		"appointment_time",
		"created_at",
	}
	return s.ExportCSV(ctx, "repair-orders.csv", records, columns)
}

func (s *ExportService) ExportParts(ctx context.Context, records []repository.Record) (*ExportFile, error) {
	columns := []string{
		"id",
		"code",
		"name",
		"specification",
		"unit",
		"unit_price",
		"stock_quantity",
		"safety_stock",
		"status",
	}
	return s.ExportCSV(ctx, "parts.csv", records, columns)
}

func (s *ExportService) ExportFeedbacks(ctx context.Context, records []repository.Record) (*ExportFile, error) {
	columns := []string{
		"id",
		"repair_order_id",
		"scheduled_at",
		"completed_at",
		"method",
		"complaint_type",
		"status",
	}
	return s.ExportCSV(ctx, "feedbacks.csv", records, columns)
}

func (s *ExportService) ExportAuditLogs(ctx context.Context, records []repository.Record) (*ExportFile, error) {
	columns := []string{
		"id",
		"user_id",
		"action",
		"resource_type",
		"resource_id",
		"ip_address",
		"request_id",
		"created_at",
	}
	return s.ExportCSV(ctx, "audit-logs.csv", records, columns)
}

func (s *ExportService) SummarizeExport(records []repository.Record) map[string]any {
	summary := map[string]any{
		"total":   len(records),
		"columns": inferColumns(records),
	}
	if len(records) > 0 {
		summary["sample"] = records[0]
	}
	return summary
}

func inferColumns(records []repository.Record) []string {
	keys := make(map[string]struct{})
	for _, record := range records {
		for key := range record {
			keys[key] = struct{}{}
		}
	}
	columns := make([]string, 0, len(keys))
	for key := range keys {
		columns = append(columns, key)
	}
	sort.Strings(columns)
	return columns
}

func stringify(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case []byte:
		return string(typed)
	default:
		return fmt.Sprintf("%v", typed)
	}
}

func ensureExtension(fileName string, ext string) string {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		fileName = "export"
	}
	if strings.HasSuffix(strings.ToLower(fileName), strings.ToLower(ext)) {
		return fileName
	}
	return fileName + ext
}
