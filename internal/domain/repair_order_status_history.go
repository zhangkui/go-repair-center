package domain

type RepairOrderStatusHistory struct {
    Base
    RepairOrderID int64 `json:"repair_order_id,omitempty"`
    FromStatus string `json:"from_status,omitempty"`
    ToStatus string `json:"to_status,omitempty"`
    Reason string `json:"reason,omitempty"`
    ChangedBy int64 `json:"changed_by,omitempty"`
}

func (m RepairOrderStatusHistory) EntityID() int64 { return m.ID }
func (m RepairOrderStatusHistory) IsDeleted() bool { return m.DeletedAt != nil }
func (m RepairOrderStatusHistory) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m RepairOrderStatusHistory) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }