package domain

type Part struct {
    Base
    Code string `json:"code,omitempty"`
    Name string `json:"name,omitempty"`
    Specification string `json:"specification,omitempty"`
    Unit string `json:"unit,omitempty"`
    UnitPrice float64 `json:"unit_price,omitempty"`
    StockQuantity int `json:"stock_quantity,omitempty"`
    SafetyStock int `json:"safety_stock,omitempty"`
    Status string `json:"status,omitempty"`
}

func (m Part) EntityID() int64 { return m.ID }
func (m Part) IsDeleted() bool { return m.DeletedAt != nil }
func (m Part) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m Part) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }