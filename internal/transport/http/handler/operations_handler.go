package handler

import (
	"net/http"
	"strconv"
	"time"

	"go-repair-center/internal/service"
	"go-repair-center/internal/transport/http/middleware"
	"go-repair-center/internal/transport/http/response"
)

type OperationsHandler struct {
	ops         *service.OpsService
	maintenance *service.MaintenanceService
}

func NewOperationsHandler(ops *service.OpsService, maintenance *service.MaintenanceService) *OperationsHandler {
	return &OperationsHandler{ops: ops, maintenance: maintenance}
}

func (h *OperationsHandler) Routes(mux interface {
	Get(string, http.HandlerFunc)
}) {
	mux.Get("/ops/backlog", h.Backlog)
	mux.Get("/ops/sla-alerts", h.SLAAlerts)
	mux.Get("/ops/warranty-expiring", h.WarrantyExpiring)
	mux.Get("/ops/health", h.HealthSnapshot)
}

func (h *OperationsHandler) Backlog(w http.ResponseWriter, r *http.Request) {
	result, err := h.ops.Backlog(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}

func (h *OperationsHandler) SLAAlerts(w http.ResponseWriter, r *http.Request) {
	result, err := h.ops.SLAAlerts(r.Context(), time.Now())
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}

func (h *OperationsHandler) WarrantyExpiring(w http.ResponseWriter, r *http.Request) {
	withinDays, _ := strconv.Atoi(r.URL.Query().Get("within_days"))
	result, err := h.ops.WarrantyExpiring(r.Context(), withinDays)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}

func (h *OperationsHandler) HealthSnapshot(w http.ResponseWriter, r *http.Request) {
	result, err := h.maintenance.Snapshot(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}
