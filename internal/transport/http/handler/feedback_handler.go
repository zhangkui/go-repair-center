package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type FeedbackHandler struct {
    *GenericHandler
}

func NewFeedbackHandler(service *service.FeedbackService) *FeedbackHandler {
    return &FeedbackHandler{GenericHandler: NewGenericHandler(service, "feedback")}
}

func (h *FeedbackHandler) Routes(router chi.Router) {
    router.Route("/feedbacks", func(resource chi.Router) {
        h.Register(resource, "feedback:read", "feedback:write")
    })
}