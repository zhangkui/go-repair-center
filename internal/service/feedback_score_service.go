package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type FeedbackScoreService struct {
    *ResourceService
}

func NewFeedbackScoreService(repo repository.CRUDRepository) *FeedbackScoreService {
    return &FeedbackScoreService{ResourceService: NewResourceService(repo, "feedback_score", "feedback_id","dimension","score")}
}

func (s *FeedbackScoreService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *FeedbackScoreService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *FeedbackScoreService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
