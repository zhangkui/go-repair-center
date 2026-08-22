package domain

type QuotationVersionHistory struct {
    Base
    QuotationID int64 `json:"quotation_id,omitempty"`
    Version int `json:"version,omitempty"`
    Snapshot string `json:"snapshot,omitempty"`
    CreatedBy int64 `json:"created_by,omitempty"`
}

func (m QuotationVersionHistory) EntityID() int64 { return m.ID }
func (m QuotationVersionHistory) IsDeleted() bool { return m.DeletedAt != nil }
func (m QuotationVersionHistory) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m QuotationVersionHistory) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }