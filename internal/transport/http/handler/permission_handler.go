package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type PermissionHandler struct {
    *GenericHandler
}

func NewPermissionHandler(service *service.PermissionService) *PermissionHandler {
    return &PermissionHandler{GenericHandler: NewGenericHandler(service, "permission")}
}

func (h *PermissionHandler) Routes(router chi.Router) {
    router.Route("/permissions", func(resource chi.Router) {
        h.Register(resource, "rbac:read", "rbac:write")
    })
}