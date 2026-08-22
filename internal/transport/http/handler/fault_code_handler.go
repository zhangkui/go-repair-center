package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type FaultCodeHandler struct {
    *GenericHandler
}

func NewFaultCodeHandler(service *service.FaultCodeService) *FaultCodeHandler {
    return &FaultCodeHandler{GenericHandler: NewGenericHandler(service, "fault_code")}
}

func (h *FaultCodeHandler) Routes(router chi.Router) {
    router.Route("/fault-codes", func(resource chi.Router) {
        h.Register(resource, "config:read", "config:write")
    })
}