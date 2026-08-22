package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type RepairTestRecordHandler struct {
    *GenericHandler
}

func NewRepairTestRecordHandler(service *service.RepairTestRecordService) *RepairTestRecordHandler {
    return &RepairTestRecordHandler{GenericHandler: NewGenericHandler(service, "repair_test_record")}
}

func (h *RepairTestRecordHandler) Routes(router chi.Router) {
    router.Route("/repair-tests", func(resource chi.Router) {
        h.Register(resource, "repair_execution:read", "repair_execution:write")
    })
}