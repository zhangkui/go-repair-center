package domain

import "time"

type RepairTestRecord struct {
    Base
    RepairExecutionID int64 `json:"repair_execution_id,omitempty"`
    TestItem string `json:"test_item,omitempty"`
    Result string `json:"result,omitempty"`
    TesterID int64 `json:"tester_id,omitempty"`
    TestedAt time.Time `json:"tested_at,omitempty"`
    Notes string `json:"notes,omitempty"`
}

func (m RepairTestRecord) EntityID() int64 { return m.ID }
func (m RepairTestRecord) IsDeleted() bool { return m.DeletedAt != nil }
func (m RepairTestRecord) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m RepairTestRecord) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }