package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-repair-center/internal/repository"
)

type InventoryControlService struct {
	db    *sql.DB
	parts *PartService
}

type RestockItem struct {
	PartID             int64   `json:"part_id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	StockQuantity      int     `json:"stock_quantity"`
	SafetyStock        int     `json:"safety_stock"`
	SuggestedQuantity  int     `json:"suggested_quantity"`
	UnitPrice          float64 `json:"unit_price"`
	EstimatedCost      float64 `json:"estimated_cost"`
}

func NewInventoryControlService(db *sql.DB, parts *PartService) *InventoryControlService {
	return &InventoryControlService{db: db, parts: parts}
}

func (s *InventoryControlService) LowStock(ctx context.Context) ([]repository.Record, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, code, name, specification, stock_quantity, safety_stock, unit_price, status
FROM parts
WHERE deleted_at IS NULL AND stock_quantity <= safety_stock
ORDER BY stock_quantity ASC, id ASC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]repository.Record, 0, 32)
	for rows.Next() {
		var id int64
		var code string
		var name string
		var specification string
		var stockQuantity int
		var safetyStock int
		var unitPrice float64
		var status string
		if err := rows.Scan(&id, &code, &name, &specification, &stockQuantity, &safetyStock, &unitPrice, &status); err != nil {
			return nil, err
		}
		items = append(items, repository.Record{
			"id":             id,
			"code":           code,
			"name":           name,
			"specification":  specification,
			"stock_quantity": stockQuantity,
			"safety_stock":   safetyStock,
			"unit_price":     unitPrice,
			"status":         status,
		})
	}
	return items, rows.Err()
}

func (s *InventoryControlService) RestockSuggestions(ctx context.Context) ([]RestockItem, error) {
	rows, err := s.LowStock(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]RestockItem, 0, len(rows))
	for _, row := range rows {
		stock := recordInt(row["stock_quantity"])
		safety := recordInt(row["safety_stock"])
		unitPrice := recordFloatValue(row["unit_price"])
		suggestedQuantity := safety*2 - stock
		if suggestedQuantity < 1 {
			suggestedQuantity = 1
		}
		items = append(items, RestockItem{
			PartID:            recordInt64ValueFromAny(row["id"]),
			Code:              fmt.Sprint(row["code"]),
			Name:              fmt.Sprint(row["name"]),
			StockQuantity:     stock,
			SafetyStock:       safety,
			SuggestedQuantity: suggestedQuantity,
			UnitPrice:         unitPrice,
			EstimatedCost:     float64(suggestedQuantity) * unitPrice,
		})
	}
	return items, nil
}

func (s *InventoryControlService) Restock(ctx context.Context, partID int64, quantity int, reason string, userID int64) (repository.Record, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("restock quantity must be positive")
	}
	return s.parts.Adjust(ctx, partID, quantity, reason, userID)
}

func (s *InventoryControlService) Timeline(ctx context.Context, partID int64, limit int) ([]repository.Record, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT id, change_type, quantity_before, change_quantity, quantity_after, reference_type, reference_id, operator_id, remark, created_at
FROM part_stock_logs
WHERE part_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ?
`, partID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]repository.Record, 0, limit)
	for rows.Next() {
		var id int64
		var changeType string
		var quantityBefore int
		var changeQuantity int
		var quantityAfter int
		var referenceType string
		var referenceID sql.NullInt64
		var operatorID int64
		var remark string
		var createdAt time.Time
		if err := rows.Scan(&id, &changeType, &quantityBefore, &changeQuantity, &quantityAfter, &referenceType, &referenceID, &operatorID, &remark, &createdAt); err != nil {
			return nil, err
		}
		items = append(items, repository.Record{
			"id":              id,
			"change_type":     changeType,
			"quantity_before": quantityBefore,
			"change_quantity": changeQuantity,
			"quantity_after":  quantityAfter,
			"reference_type":  referenceType,
			"reference_id":    referenceID.Int64,
			"operator_id":     operatorID,
			"remark":          remark,
			"created_at":      createdAt,
		})
	}
	return items, rows.Err()
}

func (s *InventoryControlService) Summary(ctx context.Context) (map[string]any, error) {
	lowStock, err := s.LowStock(ctx)
	if err != nil {
		return nil, err
	}
	suggestions, err := s.RestockSuggestions(ctx)
	if err != nil {
		return nil, err
	}
	totalEstimatedCost := 0.0
	for _, item := range suggestions {
		totalEstimatedCost += item.EstimatedCost
	}
	return map[string]any{
		"low_stock_count":      len(lowStock),
		"restock_item_count":   len(suggestions),
		"estimated_restock_cost": totalEstimatedCost,
		"suggestions":          suggestions,
	}, nil
}

func recordInt(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case []byte:
		var parsed int
		fmt.Sscanf(string(typed), "%d", &parsed)
		return parsed
	}
	return 0
}

func recordFloatValue(value any) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case int64:
		return float64(typed)
	case int:
		return float64(typed)
	case []byte:
		var parsed float64
		fmt.Sscanf(string(typed), "%f", &parsed)
		return parsed
	}
	return 0
}

func recordInt64ValueFromAny(value any) int64 {
	switch typed := value.(type) {
	case int64:
		return typed
	case int:
		return int64(typed)
	case float64:
		return int64(typed)
	case []byte:
		var parsed int64
		fmt.Sscanf(string(typed), "%d", &parsed)
		return parsed
	}
	return 0
}
