package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type RepairExecutionHandler struct {
    *GenericHandler
}

func NewRepairExecutionHandler(service *service.RepairExecutionService) *RepairExecutionHandler {
    return &RepairExecutionHandler{GenericHandler: NewGenericHandler(service, "repair_execution")}
}

func (h *RepairExecutionHandler) Routes(router chi.Router) {
    router.Route("/repair-executions", func(resource chi.Router) {
        h.Register(resource, "repair_execution:read", "repair_execution:write")
    })
}