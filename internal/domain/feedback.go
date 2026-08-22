package domain

import "time"

type Feedback struct {
    Base
    RepairOrderID int64 `json:"repair_order_id,omitempty"`
    ScheduledAt time.Time `json:"scheduled_at,omitempty"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    Method string `json:"method,omitempty"`
    ComplaintType string `json:"complaint_type,omitempty"`
    Complaint string `json:"complaint,omitempty"`
    Resolution string `json:"resolution,omitempty"`
    ReworkOrderID *int64 `json:"rework_order_id,omitempty"`
    ReworkReason string `json:"rework_reason,omitempty"`
    ReworkCost float64 `json:"rework_cost,omitempty"`
    Status string `json:"status,omitempty"`
}

func (m Feedback) EntityID() int64 { return m.ID }
func (m Feedback) IsDeleted() bool { return m.DeletedAt != nil }
func (m Feedback) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m Feedback) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }