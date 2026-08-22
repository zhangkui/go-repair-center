package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type PaymentRecordHandler struct {
    *GenericHandler
}

func NewPaymentRecordHandler(service *service.PaymentRecordService) *PaymentRecordHandler {
    return &PaymentRecordHandler{GenericHandler: NewGenericHandler(service, "payment_record")}
}

func (h *PaymentRecordHandler) Routes(router chi.Router) {
    router.Route("/payments", func(resource chi.Router) {
        h.Register(resource, "payment:read", "payment:write")
    })
}