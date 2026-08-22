package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type CustomerHandler struct {
    *GenericHandler
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
    return &CustomerHandler{GenericHandler: NewGenericHandler(service, "customer")}
}

func (h *CustomerHandler) Routes(router chi.Router) {
    router.Route("/customers", func(resource chi.Router) {
        h.Register(resource, "customer:read", "customer:write")
    })
}