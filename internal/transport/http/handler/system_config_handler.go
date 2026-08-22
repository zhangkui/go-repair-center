package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type SystemConfigHandler struct {
    *GenericHandler
}

func NewSystemConfigHandler(service *service.SystemConfigService) *SystemConfigHandler {
    return &SystemConfigHandler{GenericHandler: NewGenericHandler(service, "system_config")}
}

func (h *SystemConfigHandler) Routes(router chi.Router) {
    router.Route("/system-configs", func(resource chi.Router) {
        h.Register(resource, "config:read", "config:write")
    })
}