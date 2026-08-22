package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go-repair-center/internal/repository"
)

type ImportService struct{}

type ImportRow struct {
	LineNumber int               `json:"line_number"`
	Values     map[string]string `json:"values"`
}

type ImportError struct {
	LineNumber int    `json:"line_number"`
	Message    string `json:"message"`
}

type ImportResult struct {
	Headers  []string        `json:"headers"`
	Rows     []ImportRow     `json:"rows"`
	Valid    []repository.Record `json:"valid"`
	Errors   []ImportError   `json:"errors"`
	Inserted int             `json:"inserted"`
}

func NewImportService() *ImportService {
	return &ImportService{}
}

func (s *ImportService) ParseCSV(_ context.Context, payload []byte) (*ImportResult, error) {
	reader := csv.NewReader(bytes.NewReader(payload))
	reader.TrimLeadingSpace = true
	reader.ReuseRecord = false

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	headers = normalizeHeaders(headers)

	result := &ImportResult{
		Headers: headers,
		Rows:    make([]ImportRow, 0, 64),
		Valid:   make([]repository.Record, 0, 64),
		Errors:  make([]ImportError, 0, 16),
	}

	lineNumber := 1
	for {
		lineNumber++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, ImportError{LineNumber: lineNumber, Message: err.Error()})
			continue
		}
		values := make(map[string]string, len(headers))
		for index, header := range headers {
			if index < len(record) {
				values[header] = strings.TrimSpace(record[index])
			} else {
				values[header] = ""
			}
		}
		result.Rows = append(result.Rows, ImportRow{LineNumber: lineNumber, Values: values})
	}
	return result, nil
}

func (s *ImportService) ValidateCustomers(result *ImportResult) *ImportResult {
	for _, row := range result.Rows {
		name := row.Values["name"]
		level := row.Values["level"]
		if name == "" {
			result.Errors = append(result.Errors, ImportError{LineNumber: row.LineNumber, Message: "name is required"})
			continue
		}
		if level == "" {
			level = "NORMAL"
		}
		if !containsString([]string{"NORMAL", "VIP", "ENTERPRISE"}, level) {
			result.Errors = append(result.Errors, ImportError{LineNumber: row.LineNumber, Message: "invalid customer level"})
			continue
		}
		result.Valid = append(result.Valid, repository.Record{
			"name":    name,
			"phone":   row.Values["phone"],
			"address": row.Values["address"],
			"level":   level,
			"remark":  row.Values["remark"],
		})
	}
	return result
}

func (s *ImportService) ValidateParts(result *ImportResult) *ImportResult {
	for _, row := range result.Rows {
		code := row.Values["code"]
		name := row.Values["name"]
		if code == "" || name == "" {
			result.Errors = append(result.Errors, ImportError{LineNumber: row.LineNumber, Message: "code and name are required"})
			continue
		}
		unitPrice, err := strconv.ParseFloat(defaultString(row.Values["unit_price"], "0"), 64)
		if err != nil {
			result.Errors = append(result.Errors, ImportError{LineNumber: row.LineNumber, Message: "invalid unit_price"})
			continue
		}
		stockQuantity, err := strconv.Atoi(defaultString(row.Values["stock_quantity"], "0"))
		if err != nil {
			result.Errors = append(result.Errors, ImportError{LineNumber: row.LineNumber, Message: "invalid stock_quantity"})
			continue
		}
		safetyStock, err := strconv.Atoi(defaultString(row.Values["safety_stock"], "0"))
		if err != nil {
			result.Errors = append(result.Errors, ImportError{LineNumber: row.LineNumber, Message: "invalid safety_stock"})
			continue
		}
		result.Valid = append(result.Valid, repository.Record{
			"code":           code,
			"name":           name,
			"specification":  row.Values["specification"],
			"unit":           row.Values["unit"],
			"unit_price":     unitPrice,
			"stock_quantity": stockQuantity,
			"safety_stock":   safetyStock,
			"status":         defaultString(row.Values["status"], "ACTIVE"),
		})
	}
	return result
}

func (s *ImportService) ValidateDevices(result *ImportResult) *ImportResult {
	for _, row := range result.Rows {
		customerID, err := strconv.ParseInt(defaultString(row.Values["customer_id"], "0"), 10, 64)
		if err != nil || customerID <= 0 {
			result.Errors = append(result.Errors, ImportError{LineNumber: row.LineNumber, Message: "invalid customer_id"})
			continue
		}
		serialNumber := row.Values["serial_number"]
		if serialNumber == "" {
			result.Errors = append(result.Errors, ImportError{LineNumber: row.LineNumber, Message: "serial_number is required"})
			continue
		}
		result.Valid = append(result.Valid, repository.Record{
			"customer_id":      customerID,
			"brand":            row.Values["brand"],
			"model":            row.Values["model"],
			"serial_number":    serialNumber,
			"category":         row.Values["category"],
			"appearance_photos": row.Values["appearance_photos"],
			"purchase_channel": row.Values["purchase_channel"],
			"warranty_status":  defaultString(row.Values["warranty_status"], "ACTIVE"),
			"warranty_voucher": row.Values["warranty_voucher"],
			"status":           defaultString(row.Values["status"], "NORMAL"),
		})
	}
	return result
}

func (s *ImportService) Apply(ctx context.Context, result *ImportResult, creator func(context.Context, repository.Record) (repository.Record, error)) (*ImportResult, error) {
	for _, record := range result.Valid {
		if _, err := creator(ctx, record); err != nil {
			result.Errors = append(result.Errors, ImportError{Message: err.Error()})
			continue
		}
		result.Inserted++
	}
	return result, nil
}

func (s *ImportService) Template(headers []string) string {
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	_ = writer.Write(headers)
	writer.Flush()
	return buffer.String()
}

func normalizeHeaders(headers []string) []string {
	items := make([]string, 0, len(headers))
	for _, header := range headers {
		header = strings.TrimSpace(strings.ToLower(header))
		header = strings.ReplaceAll(header, " ", "_")
		items = append(items, header)
	}
	return items
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func (s *ImportService) Summary(result *ImportResult) map[string]any {
	return map[string]any{
		"headers":       result.Headers,
		"row_count":     len(result.Rows),
		"valid_count":   len(result.Valid),
		"error_count":   len(result.Errors),
		"inserted_count": result.Inserted,
	}
}

func (s *ImportService) ExplainError(err ImportError) string {
	if err.LineNumber > 0 {
		return fmt.Sprintf("line %d: %s", err.LineNumber, err.Message)
	}
	return err.Message
}
