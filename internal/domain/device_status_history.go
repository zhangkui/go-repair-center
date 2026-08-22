package domain

type DeviceStatusHistory struct {
    Base
    DeviceID int64 `json:"device_id,omitempty"`
    FromStatus string `json:"from_status,omitempty"`
    ToStatus string `json:"to_status,omitempty"`
    Reason string `json:"reason,omitempty"`
    ChangedBy int64 `json:"changed_by,omitempty"`
}

func (m DeviceStatusHistory) EntityID() int64 { return m.ID }
func (m DeviceStatusHistory) IsDeleted() bool { return m.DeletedAt != nil }
func (m DeviceStatusHistory) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m DeviceStatusHistory) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }