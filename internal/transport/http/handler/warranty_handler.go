package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type WarrantyHandler struct {
    *GenericHandler
}

func NewWarrantyHandler(service *service.WarrantyService) *WarrantyHandler {
    return &WarrantyHandler{GenericHandler: NewGenericHandler(service, "warranty")}
}

func (h *WarrantyHandler) Routes(router chi.Router) {
    router.Route("/warranties", func(resource chi.Router) {
        h.Register(resource, "warranty:read", "warranty:write")
    })
}