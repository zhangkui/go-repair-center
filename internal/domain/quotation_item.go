package domain

type QuotationItem struct {
    Base
    QuotationID int64 `json:"quotation_id,omitempty"`
    ItemType string `json:"item_type,omitempty"`
    PartID *int64 `json:"part_id,omitempty"`
    Name string `json:"name,omitempty"`
    Quantity float64 `json:"quantity,omitempty"`
    UnitPrice float64 `json:"unit_price,omitempty"`
    Amount float64 `json:"amount,omitempty"`
    Description string `json:"description,omitempty"`
}

func (m QuotationItem) EntityID() int64 { return m.ID }
func (m QuotationItem) IsDeleted() bool { return m.DeletedAt != nil }
func (m QuotationItem) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m QuotationItem) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }