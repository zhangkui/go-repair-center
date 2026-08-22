package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type PartStockLogHandler struct {
    *GenericHandler
}

func NewPartStockLogHandler(service *service.PartStockLogService) *PartStockLogHandler {
    return &PartStockLogHandler{GenericHandler: NewGenericHandler(service, "part_stock_log")}
}

func (h *PartStockLogHandler) Routes(router chi.Router) {
    router.Route("/part-stock-logs", func(resource chi.Router) {
        h.Register(resource, "part:read", "part:write")
    })
}