package domain

type Permission struct {
    Base
    Code string `json:"code,omitempty"`
    Name string `json:"name,omitempty"`
    Module string `json:"module,omitempty"`
    Action string `json:"action,omitempty"`
    Description string `json:"description,omitempty"`
}

func (m Permission) EntityID() int64 { return m.ID }
func (m Permission) IsDeleted() bool { return m.DeletedAt != nil }
func (m Permission) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m Permission) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }