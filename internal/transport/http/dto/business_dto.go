package dto

type StockUsageRequest struct { PartID int64 `json:"part_id"`; Quantity int `json:"quantity"`; RepairExecutionID int64 `json:"repair_execution_id"` }
type StockAdjustmentRequest struct { Delta int `json:"delta"`; Reason string `json:"reason"` }
type QuotationVersionRequest struct { RepairOrderID int64 `json:"repair_order_id"`; LaborFee float64 `json:"labor_fee"`; PartsFee float64 `json:"parts_fee"`; InspectionFee float64 `json:"inspection_fee"`; OtherFee float64 `json:"other_fee"`; Items []map[string]any `json:"items"` }
type ApprovalRequest struct { Approved bool `json:"approved"`; Comment string `json:"comment"` }
type StatusChangeRequest struct { CurrentStatus string `json:"current_status"`; NextStatus string `json:"next_status"`; Reason string `json:"reason"` }