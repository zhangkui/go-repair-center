package handler

import (
    "github.com/go-chi/chi/v5"
    "go-repair-center/internal/service"
)

type DeviceHandler struct {
    *GenericHandler
}

func NewDeviceHandler(service *service.DeviceService) *DeviceHandler {
    return &DeviceHandler{GenericHandler: NewGenericHandler(service, "device")}
}

func (h *DeviceHandler) Routes(router chi.Router) {
    router.Route("/devices", func(resource chi.Router) {
        h.Register(resource, "device:read", "device:write")
    })
}