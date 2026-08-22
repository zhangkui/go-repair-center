package domain

import "time"

type Warranty struct {
    Base
    RepairOrderID int64 `json:"repair_order_id,omitempty"`
    WarrantyType string `json:"warranty_type,omitempty"`
    DurationDays int `json:"duration_days,omitempty"`
    StartDate time.Time `json:"start_date,omitempty"`
    EndDate time.Time `json:"end_date,omitempty"`
    Status string `json:"status,omitempty"`
    ApprovedBy *int64 `json:"approved_by,omitempty"`
}

func (m Warranty) EntityID() int64 { return m.ID }
func (m Warranty) IsDeleted() bool { return m.DeletedAt != nil }
func (m Warranty) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m Warranty) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }