package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go-repair-center/internal/repository"
	"go-repair-center/internal/service"
	"go-repair-center/internal/transport/http/dto"
	"go-repair-center/internal/transport/http/middleware"
	"go-repair-center/internal/transport/http/response"
)

type DataToolsHandler struct {
	importer     *service.ImportService
	exporter     *service.ExportService
	customers    CRUDService
	devices      CRUDService
	parts        CRUDService
	repairOrders CRUDService
	feedbacks    CRUDService
	audits       CRUDService
}

func NewDataToolsHandler(
	importer *service.ImportService,
	exporter *service.ExportService,
	customers CRUDService,
	devices CRUDService,
	parts CRUDService,
	repairOrders CRUDService,
	feedbacks CRUDService,
	audits CRUDService,
) *DataToolsHandler {
	return &DataToolsHandler{
		importer:     importer,
		exporter:     exporter,
		customers:    customers,
		devices:      devices,
		parts:        parts,
		repairOrders: repairOrders,
		feedbacks:    feedbacks,
		audits:       audits,
	}
}

func (h *DataToolsHandler) Routes(mux interface {
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
}) {
	mux.Get("/tools/import/templates/customers", h.CustomerTemplate)
	mux.Get("/tools/import/templates/parts", h.PartTemplate)
	mux.Get("/tools/import/templates/devices", h.DeviceTemplate)
	mux.Post("/tools/import/customers/preview", h.PreviewCustomers)
	mux.Post("/tools/import/customers/apply", h.ApplyCustomers)
	mux.Post("/tools/import/parts/preview", h.PreviewParts)
	mux.Post("/tools/import/parts/apply", h.ApplyParts)
	mux.Post("/tools/import/devices/preview", h.PreviewDevices)
	mux.Post("/tools/import/devices/apply", h.ApplyDevices)
	mux.Get("/tools/export/repair-orders.csv", h.ExportRepairOrders)
	mux.Get("/tools/export/parts.csv", h.ExportParts)
	mux.Get("/tools/export/feedbacks.csv", h.ExportFeedbacks)
	mux.Get("/tools/export/audit-logs.csv", h.ExportAuditLogs)
}

func (h *DataToolsHandler) CustomerTemplate(w http.ResponseWriter, r *http.Request) {
	h.writeCSVTemplate(w, "customers-template.csv", []string{"name", "phone", "address", "level", "remark"})
}

func (h *DataToolsHandler) PartTemplate(w http.ResponseWriter, r *http.Request) {
	h.writeCSVTemplate(w, "parts-template.csv", []string{"code", "name", "specification", "unit", "unit_price", "stock_quantity", "safety_stock", "status"})
}

func (h *DataToolsHandler) DeviceTemplate(w http.ResponseWriter, r *http.Request) {
	h.writeCSVTemplate(w, "devices-template.csv", []string{"customer_id", "brand", "model", "serial_number", "category", "appearance_photos", "purchase_channel", "warranty_status", "warranty_voucher", "status"})
}

func (h *DataToolsHandler) PreviewCustomers(w http.ResponseWriter, r *http.Request) {
	h.previewImport(w, r, h.importer.ValidateCustomers)
}

func (h *DataToolsHandler) ApplyCustomers(w http.ResponseWriter, r *http.Request) {
	h.applyImport(w, r, h.importer.ValidateCustomers, h.customers.Create)
}

func (h *DataToolsHandler) PreviewParts(w http.ResponseWriter, r *http.Request) {
	h.previewImport(w, r, h.importer.ValidateParts)
}

func (h *DataToolsHandler) ApplyParts(w http.ResponseWriter, r *http.Request) {
	h.applyImport(w, r, h.importer.ValidateParts, h.parts.Create)
}

func (h *DataToolsHandler) PreviewDevices(w http.ResponseWriter, r *http.Request) {
	h.previewImport(w, r, h.importer.ValidateDevices)
}

func (h *DataToolsHandler) ApplyDevices(w http.ResponseWriter, r *http.Request) {
	h.applyImport(w, r, h.importer.ValidateDevices, h.devices.Create)
}

func (h *DataToolsHandler) ExportRepairOrders(w http.ResponseWriter, r *http.Request) {
	h.exportCSV(w, r, "repair-orders", h.repairOrders, h.exporter.ExportRepairOrders)
}

func (h *DataToolsHandler) ExportParts(w http.ResponseWriter, r *http.Request) {
	h.exportCSV(w, r, "parts", h.parts, h.exporter.ExportParts)
}

func (h *DataToolsHandler) ExportFeedbacks(w http.ResponseWriter, r *http.Request) {
	h.exportCSV(w, r, "feedbacks", h.feedbacks, h.exporter.ExportFeedbacks)
}

func (h *DataToolsHandler) ExportAuditLogs(w http.ResponseWriter, r *http.Request) {
	h.exportCSV(w, r, "audit-logs", h.audits, h.exporter.ExportAuditLogs)
}

func (h *DataToolsHandler) previewImport(
	w http.ResponseWriter,
	r *http.Request,
	validate func(*service.ImportResult) *service.ImportResult,
) {
	result, err := h.parseImportRequest(r.Context(), r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	result = validate(result)
	response.Success(w, map[string]any{
		"summary": h.importer.Summary(result),
		"result":  result,
	}, middleware.RequestID(r.Context()))
}

func (h *DataToolsHandler) applyImport(
	w http.ResponseWriter,
	r *http.Request,
	validate func(*service.ImportResult) *service.ImportResult,
	create func(context.Context, repository.Record) (repository.Record, error),
) {
	result, err := h.parseImportRequest(r.Context(), r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	result = validate(result)
	result, err = h.importer.Apply(r.Context(), result, create)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, map[string]any{
		"summary": h.importer.Summary(result),
		"result":  result,
	}, middleware.RequestID(r.Context()))
}

func (h *DataToolsHandler) parseImportRequest(ctx context.Context, r *http.Request) (*service.ImportResult, error) {
	var payload dto.CSVImportRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return h.importer.ParseCSV(ctx, []byte(payload.CSVContent))
}

func (h *DataToolsHandler) writeCSVTemplate(w http.ResponseWriter, fileName string, headers []string) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(h.importer.Template(headers)))
}

func (h *DataToolsHandler) exportCSV(
	w http.ResponseWriter,
	r *http.Request,
	fallbackName string,
	resource CRUDService,
	exporter func(context.Context, []repository.Record) (*service.ExportFile, error),
) {
	page, err := resource.List(r.Context(), 1, 500, r.URL.Query().Get("keyword"), r.URL.Query().Get("status"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	file, err := exporter(r.Context(), page.Items)
	if err != nil {
		writeError(w, r, err)
		return
	}
	if file.FileName == "" {
		file.FileName = fallbackName + ".csv"
	}
	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+file.FileName+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(file.Content)
}
