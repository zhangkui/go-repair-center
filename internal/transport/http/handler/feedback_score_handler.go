package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type FeedbackScoreHandler struct {
    *GenericHandler
}

func NewFeedbackScoreHandler(service *service.FeedbackScoreService) *FeedbackScoreHandler {
    return &FeedbackScoreHandler{GenericHandler: NewGenericHandler(service, "feedback_score")}
}

func (h *FeedbackScoreHandler) Routes(router chi.Router) {
    router.Route("/feedback-scores", func(resource chi.Router) {
        h.Register(resource, "feedback:read", "feedback:write")
    })
}