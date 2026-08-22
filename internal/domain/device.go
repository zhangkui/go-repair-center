package domain

import "time"

type Device struct {
    Base
    CustomerID int64 `json:"customer_id,omitempty"`
    Brand string `json:"brand,omitempty"`
    Model string `json:"model,omitempty"`
    SerialNumber string `json:"serial_number,omitempty"`
    Category string `json:"category,omitempty"`
    AppearancePhotos string `json:"appearance_photos,omitempty"`
    PurchaseDate *time.Time `json:"purchase_date,omitempty"`
    PurchaseChannel string `json:"purchase_channel,omitempty"`
    PurchasePrice float64 `json:"purchase_price,omitempty"`
    WarrantyStart *time.Time `json:"warranty_start,omitempty"`
    WarrantyEnd *time.Time `json:"warranty_end,omitempty"`
    WarrantyStatus string `json:"warranty_status,omitempty"`
    WarrantyVoucher string `json:"warranty_voucher,omitempty"`
    Status string `json:"status,omitempty"`
}

func (m Device) EntityID() int64 { return m.ID }
func (m Device) IsDeleted() bool { return m.DeletedAt != nil }
func (m Device) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m Device) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }