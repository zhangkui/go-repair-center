package domain

type Customer struct {
    Base
    Name string `json:"name,omitempty"`
    Phone string `json:"phone,omitempty"`
    Address string `json:"address,omitempty"`
    Level string `json:"level,omitempty"`
    Remark string `json:"remark,omitempty"`
}

func (m Customer) EntityID() int64 { return m.ID }
func (m Customer) IsDeleted() bool { return m.DeletedAt != nil }
func (m Customer) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m Customer) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }