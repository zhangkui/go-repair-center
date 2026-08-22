package domain

import "time"

type Quotation struct {
    Base
    RepairOrderID int64 `json:"repair_order_id,omitempty"`
    Version int `json:"version,omitempty"`
    LaborFee float64 `json:"labor_fee,omitempty"`
    PartsFee float64 `json:"parts_fee,omitempty"`
    InspectionFee float64 `json:"inspection_fee,omitempty"`
    OtherFee float64 `json:"other_fee,omitempty"`
    TotalAmount float64 `json:"total_amount,omitempty"`
    ValidUntil time.Time `json:"valid_until,omitempty"`
    CustomerStatus string `json:"customer_status,omitempty"`
    ApprovalStatus string `json:"approval_status,omitempty"`
    ApprovedBy *int64 `json:"approved_by,omitempty"`
    ApprovedAt *time.Time `json:"approved_at,omitempty"`
    Status string `json:"status,omitempty"`
}

func (m Quotation) EntityID() int64 { return m.ID }
func (m Quotation) IsDeleted() bool { return m.DeletedAt != nil }
func (m Quotation) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m Quotation) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }