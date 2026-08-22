package domain

type FaultCode struct {
    Base
    Code string `json:"code,omitempty"`
    Name string `json:"name,omitempty"`
    Category string `json:"category,omitempty"`
    Description string `json:"description,omitempty"`
    Status string `json:"status,omitempty"`
}

func (m FaultCode) EntityID() int64 { return m.ID }
func (m FaultCode) IsDeleted() bool { return m.DeletedAt != nil }
func (m FaultCode) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m FaultCode) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }