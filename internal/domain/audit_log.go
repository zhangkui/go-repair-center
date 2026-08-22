package domain

type AuditLog struct {
    Base
    UserID *int64 `json:"user_id,omitempty"`
    Action string `json:"action,omitempty"`
    ResourceType string `json:"resource_type,omitempty"`
    ResourceID string `json:"resource_id,omitempty"`
    BeforeValue string `json:"before_value,omitempty"`
    AfterValue string `json:"after_value,omitempty"`
    IPAddress string `json:"ipaddress,omitempty"`
    UserAgent string `json:"user_agent,omitempty"`
    RequestID string `json:"request_id,omitempty"`
}

func (m AuditLog) EntityID() int64 { return m.ID }
func (m AuditLog) IsDeleted() bool { return m.DeletedAt != nil }
func (m AuditLog) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m AuditLog) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }