package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type RepairOrderHandler struct {
    *GenericHandler
}

func NewRepairOrderHandler(service *service.RepairOrderService) *RepairOrderHandler {
    return &RepairOrderHandler{GenericHandler: NewGenericHandler(service, "repair_order")}
}

func (h *RepairOrderHandler) Routes(router chi.Router) {
    router.Route("/repair-orders", func(resource chi.Router) {
        h.Register(resource, "repair_order:read", "repair_order:write")
    })
}