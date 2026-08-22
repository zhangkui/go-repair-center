package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type RoleHandler struct {
    *GenericHandler
}

func NewRoleHandler(service *service.RoleService) *RoleHandler {
    return &RoleHandler{GenericHandler: NewGenericHandler(service, "role")}
}

func (h *RoleHandler) Routes(router chi.Router) {
    router.Route("/roles", func(resource chi.Router) {
        h.Register(resource, "rbac:read", "rbac:write")
    })
}