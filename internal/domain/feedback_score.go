package domain

type FeedbackScore struct {
    Base
    FeedbackID int64 `json:"feedback_id,omitempty"`
    Dimension string `json:"dimension,omitempty"`
    Score int `json:"score,omitempty"`
    Comment string `json:"comment,omitempty"`
}

func (m FeedbackScore) EntityID() int64 { return m.ID }
func (m FeedbackScore) IsDeleted() bool { return m.DeletedAt != nil }
func (m FeedbackScore) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m FeedbackScore) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }