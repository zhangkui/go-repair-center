package domain

import "time"

type RepairOrder struct {
    Base
    OrderNumber string `json:"order_number,omitempty"`
    CustomerID int64 `json:"customer_id,omitempty"`
    DeviceID int64 `json:"device_id,omitempty"`
    OriginalOrderID *int64 `json:"original_order_id,omitempty"`
    FaultDescription string `json:"fault_description,omitempty"`
    Urgency string `json:"urgency,omitempty"`
    ServiceMethod string `json:"service_method,omitempty"`
    ExpectedCompletedAt *time.Time `json:"expected_completed_at,omitempty"`
    AppointmentTime *time.Time `json:"appointment_time,omitempty"`
    AppointmentEnd *time.Time `json:"appointment_end,omitempty"`
    TechnicianID *int64 `json:"technician_id,omitempty"`
    DispatchedAt *time.Time `json:"dispatched_at,omitempty"`
    NotificationStatus string `json:"notification_status,omitempty"`
    Status string `json:"status,omitempty"`
}

func (m RepairOrder) EntityID() int64 { return m.ID }
func (m RepairOrder) IsDeleted() bool { return m.DeletedAt != nil }
func (m RepairOrder) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m RepairOrder) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }