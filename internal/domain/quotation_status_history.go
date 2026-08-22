package domain

type QuotationStatusHistory struct {
    Base
    QuotationID int64 `json:"quotation_id,omitempty"`
    FromStatus string `json:"from_status,omitempty"`
    ToStatus string `json:"to_status,omitempty"`
    Reason string `json:"reason,omitempty"`
    ChangedBy int64 `json:"changed_by,omitempty"`
}

func (m QuotationStatusHistory) EntityID() int64 { return m.ID }
func (m QuotationStatusHistory) IsDeleted() bool { return m.DeletedAt != nil }
func (m QuotationStatusHistory) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m QuotationStatusHistory) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }