package handler

import (
	"net/http"
	"strconv"

	"go-repair-center/internal/service"
	"go-repair-center/internal/transport/http/middleware"
	"go-repair-center/internal/transport/http/response"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) Revenue(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	report, err := h.service.RevenueAndPerformance(r.Context(), days)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, report, middleware.RequestID(r.Context()))
}

func (h *ReportHandler) RevenueCSV(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	payload, err := h.service.RevenueCSV(r.Context(), days)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="revenue-report.csv"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(payload))
}
