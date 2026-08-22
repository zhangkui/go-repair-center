package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type FeedbackService struct {
    *ResourceService
}

func NewFeedbackService(repo repository.CRUDRepository) *FeedbackService {
    return &FeedbackService{ResourceService: NewResourceService(repo, "feedback", "repair_order_id","scheduled_at")}
}

func (s *FeedbackService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *FeedbackService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *FeedbackService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}

var feedbackTransitions = map[string][]string{"PENDING":{"VISITED","CLOSED"},"VISITED":{"CLOSED"},"CLOSED":{}}

func (s *FeedbackService) ChangeStatus(ctx context.Context, id int64, current, next, reason string) (repository.Record, error) {
    if err := ValidateTransition(current, next, feedbackTransitions); err != nil { return nil, err }
    return s.Update(ctx, id, repository.Record{"status": next, "status_reason": reason})
}