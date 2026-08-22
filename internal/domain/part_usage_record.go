package domain

type PartUsageRecord struct {
    Base
    RepairExecutionID int64 `json:"repair_execution_id,omitempty"`
    PartID int64 `json:"part_id,omitempty"`
    Quantity int `json:"quantity,omitempty"`
    UnitPrice float64 `json:"unit_price,omitempty"`
    UsedBy int64 `json:"used_by,omitempty"`
    IdempotencyKey string `json:"idempotency_key,omitempty"`
}

func (m PartUsageRecord) EntityID() int64 { return m.ID }
func (m PartUsageRecord) IsDeleted() bool { return m.DeletedAt != nil }
func (m PartUsageRecord) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m PartUsageRecord) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }