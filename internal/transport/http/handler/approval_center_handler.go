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

type ApprovalCenterHandler struct {
	approvals *service.ApprovalCenterService
}

func NewApprovalCenterHandler(approvals *service.ApprovalCenterService) *ApprovalCenterHandler {
	return &ApprovalCenterHandler{approvals: approvals}
}

func (h *ApprovalCenterHandler) Routes(router chi.Router) {
	router.With(middleware.Require("quotation:write")).Get("/approvals/overview", h.Overview)
	router.With(middleware.Require("quotation:write")).Get("/approvals/quotations", h.PendingQuotations)
	router.With(middleware.Require("quotation:write")).Post("/approvals/quotations/{id}/decision", h.DecideQuotation)
	router.With(middleware.Require("warranty:write")).Get("/approvals/warranties", h.PendingWarranties)
	router.With(middleware.Require("warranty:write")).Post("/approvals/warranties/{id}/decision", h.DecideWarranty)
}

func (h *ApprovalCenterHandler) Overview(w http.ResponseWriter, r *http.Request) {
	result, err := h.approvals.Overview(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}

func (h *ApprovalCenterHandler) PendingQuotations(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	pageSize := queryInt(r, "page_size", 20)
	result, err := h.approvals.PendingQuotations(r.Context(), page, pageSize)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, response.PageData{Items: result.Items, Pagination: response.Pagination{Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}}, middleware.RequestID(r.Context()))
}

func (h *ApprovalCenterHandler) PendingWarranties(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	pageSize := queryInt(r, "page_size", 20)
	result, err := h.approvals.PendingWarranties(r.Context(), page, pageSize)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, response.PageData{Items: result.Items, Pagination: response.Pagination{Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}}, middleware.RequestID(r.Context()))
}

func (h *ApprovalCenterHandler) DecideQuotation(w http.ResponseWriter, r *http.Request) {
	var payload dto.ApprovalDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, r, err)
		return
	}
	quotationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, r, err)
		return
	}
	result, err := h.approvals.DecideQuotation(r.Context(), quotationID, middleware.UserID(r.Context()), payload.Approved, payload.Comment)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}

func (h *ApprovalCenterHandler) DecideWarranty(w http.ResponseWriter, r *http.Request) {
	var payload dto.ApprovalDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, r, err)
		return
	}
	warrantyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, r, err)
		return
	}
	result, err := h.approvals.DecideWarranty(r.Context(), warrantyID, middleware.UserID(r.Context()), payload.Approved, payload.Comment)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, result, middleware.RequestID(r.Context()))
}
