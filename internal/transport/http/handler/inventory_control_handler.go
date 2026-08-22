package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go-repair-center/internal/service"
	"go-repair-center/internal/transport/http/dto"
	"go-repair-center/internal/transport/http/middleware"
	"go-repair-center/internal/transport/http/response"
)

type InventoryControlHandler struct {
	inventory *service.InventoryControlService
}

func NewInventoryControlHandler(inventory *service.InventoryControlService) *InventoryControlHandler {
	return &InventoryControlHandler{inventory: inventory}
}

func (h *InventoryControlHandler) Routes(router chi.Router) {
	router.With(middleware.Require("part:read")).Get("/inventory/low-stock", h.LowStock)
	router.With(middleware.Require("part:read")).Get("/inventory/restock-suggestions", h.RestockSuggestions)
	router.With(middleware.Require("part:read")).Get("/inventory/summary", h.Summary)
	router.With(middleware.Require("part:read")).Get("/inventory/parts/{id}/timeline", h.Timeline)
	router.With(middleware.Require("part:write")).Post("/inventory/restock", h.Restock)
}

func (h *InventoryControlHandler) LowStock(w http.ResponseWriter, r *http.Request) {
	items, err := h.inventory.LowStock(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, items, middleware.RequestID(r.Context()))
}

func (h *InventoryControlHandler) RestockSuggestions(w http.ResponseWriter, r *http.Request) {
	items, err := h.inventory.RestockSuggestions(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, items, middleware.RequestID(r.Context()))
}

func (h *InventoryControlHandler) Summary(w http.ResponseWriter, r *http.Request) {
	result, err := h.inventory.Summary(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}

func (h *InventoryControlHandler) Timeline(w http.ResponseWriter, r *http.Request) {
	partID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, r, err)
		return
	}
	limit := queryInt(r, "limit", 20)
	items, err := h.inventory.Timeline(r.Context(), partID, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, items, middleware.RequestID(r.Context()))
}

func (h *InventoryControlHandler) Restock(w http.ResponseWriter, r *http.Request) {
	var payload dto.RestockRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, r, err)
		return
	}
	result, err := h.inventory.Restock(r.Context(), payload.PartID, payload.Quantity, payload.Reason, middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}
