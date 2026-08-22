package domain

type PartStockLog struct {
    Base
    PartID int64 `json:"part_id,omitempty"`
    ChangeType string `json:"change_type,omitempty"`
    QuantityBefore int `json:"quantity_before,omitempty"`
    ChangeQuantity int `json:"change_quantity,omitempty"`
    QuantityAfter int `json:"quantity_after,omitempty"`
    ReferenceType string `json:"reference_type,omitempty"`
    ReferenceID *int64 `json:"reference_id,omitempty"`
    OperatorID int64 `json:"operator_id,omitempty"`
    Remark string `json:"remark,omitempty"`
}

func (m PartStockLog) EntityID() int64 { return m.ID }
func (m PartStockLog) IsDeleted() bool { return m.DeletedAt != nil }
func (m PartStockLog) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m PartStockLog) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }