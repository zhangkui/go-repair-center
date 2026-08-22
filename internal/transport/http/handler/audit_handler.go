package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type AuditHandler struct {
    *GenericHandler
}

func NewAuditHandler(service *service.AuditService) *AuditHandler {
    return &AuditHandler{GenericHandler: NewGenericHandler(service, "audit")}
}

func (h *AuditHandler) Routes(router chi.Router) {
    router.Route("/audit-logs", func(resource chi.Router) {
        h.Register(resource, "audit:read", "audit:write")
    })
}