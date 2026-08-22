package domain

type WarrantyStatusHistory struct {
    Base
    WarrantyID int64 `json:"warranty_id,omitempty"`
    FromStatus string `json:"from_status,omitempty"`
    ToStatus string `json:"to_status,omitempty"`
    Reason string `json:"reason,omitempty"`
    ChangedBy int64 `json:"changed_by,omitempty"`
}

func (m WarrantyStatusHistory) EntityID() int64 { return m.ID }
func (m WarrantyStatusHistory) IsDeleted() bool { return m.DeletedAt != nil }
func (m WarrantyStatusHistory) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m WarrantyStatusHistory) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }