package domain

type Role struct {
    Base
    Code string `json:"code,omitempty"`
    Name string `json:"name,omitempty"`
    Description string `json:"description,omitempty"`
}

func (m Role) EntityID() int64 { return m.ID }
func (m Role) IsDeleted() bool { return m.DeletedAt != nil }
func (m Role) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m Role) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }