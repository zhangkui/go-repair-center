package domain

import "time"

type PaymentRecord struct {
    Base
    RepairOrderID int64 `json:"repair_order_id,omitempty"`
    PaymentMethod string `json:"payment_method,omitempty"`
    PaidAt time.Time `json:"paid_at,omitempty"`
    Amount float64 `json:"amount,omitempty"`
    VoucherNumber string `json:"voucher_number,omitempty"`
    IdempotencyKey string `json:"idempotency_key,omitempty"`
    OperatorID int64 `json:"operator_id,omitempty"`
}

func (m PaymentRecord) EntityID() int64 { return m.ID }
func (m PaymentRecord) IsDeleted() bool { return m.DeletedAt != nil }
func (m PaymentRecord) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m PaymentRecord) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }