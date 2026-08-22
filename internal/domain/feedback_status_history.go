package domain

type FeedbackStatusHistory struct {
    Base
    FeedbackID int64 `json:"feedback_id,omitempty"`
    FromStatus string `json:"from_status,omitempty"`
    ToStatus string `json:"to_status,omitempty"`
    Reason string `json:"reason,omitempty"`
    ChangedBy int64 `json:"changed_by,omitempty"`
}

func (m FeedbackStatusHistory) EntityID() int64 { return m.ID }
func (m FeedbackStatusHistory) IsDeleted() bool { return m.DeletedAt != nil }
func (m FeedbackStatusHistory) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m FeedbackStatusHistory) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }