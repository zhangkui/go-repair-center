package domain

import "time"

type RepairExecution struct {
    Base
    RepairOrderID int64 `json:"repair_order_id,omitempty"`
    TechnicianID int64 `json:"technician_id,omitempty"`
    FaultCodeID *int64 `json:"fault_code_id,omitempty"`
    Diagnosis string `json:"diagnosis,omitempty"`
    StartedAt *time.Time `json:"started_at,omitempty"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    Status string `json:"status,omitempty"`
}

func (m RepairExecution) EntityID() int64 { return m.ID }
func (m RepairExecution) IsDeleted() bool { return m.DeletedAt != nil }
func (m RepairExecution) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m RepairExecution) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }