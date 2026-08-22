package dto

type BatchDispatchRequest struct {
	OrderIDs         []int64 `json:"order_ids"`
	TechnicianID     int64   `json:"technician_id"`
	TechnicianName   string  `json:"technician_name"`
	AppointmentTime  string  `json:"appointment_time"`
	AppointmentEnd   string  `json:"appointment_end"`
	ServiceMethod    string  `json:"service_method"`
	NotifyTechnician bool    `json:"notify_technician"`
}

type ConflictCheckRequest struct {
	TechnicianID    int64  `json:"technician_id"`
	AppointmentTime string `json:"appointment_time"`
	AppointmentEnd  string `json:"appointment_end"`
}

type ApprovalDecisionRequest struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment"`
}

type RestockRequest struct {
	PartID   int64  `json:"part_id"`
	Quantity int    `json:"quantity"`
	Reason   string `json:"reason"`
}
