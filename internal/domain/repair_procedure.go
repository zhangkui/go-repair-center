package domain

import "time"

type RepairProcedure struct {
    Base
    RepairExecutionID int64 `json:"repair_execution_id,omitempty"`
    Name string `json:"name,omitempty"`
    Sequence int `json:"sequence,omitempty"`
    Description string `json:"description,omitempty"`
    StandardMinutes int `json:"standard_minutes,omitempty"`
    TechnicianID *int64 `json:"technician_id,omitempty"`
    StartedAt *time.Time `json:"started_at,omitempty"`
    EndedAt *time.Time `json:"ended_at,omitempty"`
    ActualMinutes int `json:"actual_minutes,omitempty"`
    Result string `json:"result,omitempty"`
}

func (m RepairProcedure) EntityID() int64 { return m.ID }
func (m RepairProcedure) IsDeleted() bool { return m.DeletedAt != nil }
func (m RepairProcedure) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m RepairProcedure) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }