package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go-repair-center/internal/service"
	"go-repair-center/internal/transport/http/dto"
	"go-repair-center/internal/transport/http/middleware"
	"go-repair-center/internal/transport/http/response"
)

type DispatchCenterHandler struct {
	dispatch *service.DispatchService
}

func NewDispatchCenterHandler(dispatch *service.DispatchService) *DispatchCenterHandler {
	return &DispatchCenterHandler{dispatch: dispatch}
}

func (h *DispatchCenterHandler) Routes(router chi.Router) {
	router.With(middleware.Require("repair_order:write")).Post("/dispatch/conflicts", h.CheckConflicts)
	router.With(middleware.Require("repair_order:write")).Post("/dispatch/batch", h.BatchDispatch)
	router.With(middleware.Require("repair_order:read")).Get("/dispatch/rework-chain/{id}", h.ReworkChain)
}

func (h *DispatchCenterHandler) CheckConflicts(w http.ResponseWriter, r *http.Request) {
	var payload dto.ConflictCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, r, err)
		return
	}
	start, end, err := parseWindow(payload.AppointmentTime, payload.AppointmentEnd)
	if err != nil {
		writeError(w, r, err)
		return
	}
	items, err := h.dispatch.CheckConflicts(r.Context(), payload.TechnicianID, start, end)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, items, middleware.RequestID(r.Context()))
}

func (h *DispatchCenterHandler) BatchDispatch(w http.ResponseWriter, r *http.Request) {
	var payload dto.BatchDispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, r, err)
		return
	}
	start, end, err := parseWindow(payload.AppointmentTime, payload.AppointmentEnd)
	if err != nil {
		writeError(w, r, err)
		return
	}
	result, err := h.dispatch.BatchDispatch(r.Context(), service.BatchDispatchRequest{
		OrderIDs:         payload.OrderIDs,
		TechnicianID:     payload.TechnicianID,
		TechnicianName:   payload.TechnicianName,
		AppointmentTime:  start,
		AppointmentEnd:   end,
		ServiceMethod:    payload.ServiceMethod,
		NotifyTechnician: payload.NotifyTechnician,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}

func (h *DispatchCenterHandler) ReworkChain(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, r, err)
		return
	}
	items, err := h.dispatch.ReworkChain(r.Context(), orderID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, items, middleware.RequestID(r.Context()))
}

func parseWindow(start string, end string) (time.Time, time.Time, error) {
	parsedStart, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	parsedEnd, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return parsedStart, parsedEnd, nil
}
