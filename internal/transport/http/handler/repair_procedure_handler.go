package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type RepairProcedureHandler struct {
    *GenericHandler
}

func NewRepairProcedureHandler(service *service.RepairProcedureService) *RepairProcedureHandler {
    return &RepairProcedureHandler{GenericHandler: NewGenericHandler(service, "repair_procedure")}
}

func (h *RepairProcedureHandler) Routes(router chi.Router) {
    router.Route("/repair-procedures", func(resource chi.Router) {
        h.Register(resource, "repair_execution:read", "repair_execution:write")
    })
}