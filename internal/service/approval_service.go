package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-repair-center/internal/repository"
)

type ApprovalService struct {
	quotationThreshold float64
	warrantyDays       int
}

type ApprovalDecision struct {
	ResourceType string            `json:"resource_type"`
	ResourceID   int64             `json:"resource_id"`
	NeedReview   bool              `json:"need_review"`
	Approved     bool              `json:"approved"`
	Reason       string            `json:"reason"`
	ReviewerID   int64             `json:"reviewer_id"`
	DecidedAt    time.Time         `json:"decided_at"`
	Metadata     map[string]string `json:"metadata"`
}

type PendingApprovalSummary struct {
	Total          int64              `json:"total"`
	ByResourceType map[string]int64   `json:"by_resource_type"`
	Decisions      []ApprovalDecision `json:"decisions"`
}

func NewApprovalService(quotationThreshold float64, warrantyDays int) *ApprovalService {
	return &ApprovalService{
		quotationThreshold: quotationThreshold,
		warrantyDays:       warrantyDays,
	}
}

func (s *ApprovalService) ReviewQuotation(resourceID int64, amount float64, reviewerID int64, approve bool, comment string) ApprovalDecision {
	needReview := amount > s.quotationThreshold
	reason := strings.TrimSpace(comment)
	if reason == "" {
		if needReview {
			reason = fmt.Sprintf("quotation amount %.2f exceeds threshold %.2f", amount, s.quotationThreshold)
		} else {
			reason = "quotation below review threshold"
		}
	}
	return ApprovalDecision{
		ResourceType: "quotation",
		ResourceID:   resourceID,
		NeedReview:   needReview,
		Approved:     !needReview || approve,
		Reason:       reason,
		ReviewerID:   reviewerID,
		DecidedAt:    time.Now(),
		Metadata: map[string]string{
			"amount":    fmt.Sprintf("%.2f", amount),
			"threshold": fmt.Sprintf("%.2f", s.quotationThreshold),
		},
	}
}

func (s *ApprovalService) ReviewWarranty(resourceID int64, durationDays int, reviewerID int64, approve bool, comment string) ApprovalDecision {
	needReview := durationDays > s.warrantyDays
	reason := strings.TrimSpace(comment)
	if reason == "" {
		if needReview {
			reason = fmt.Sprintf("warranty duration %d exceeds default %d", durationDays, s.warrantyDays)
		} else {
			reason = "warranty duration within default range"
		}
	}
	return ApprovalDecision{
		ResourceType: "warranty",
		ResourceID:   resourceID,
		NeedReview:   needReview,
		Approved:     !needReview || approve,
		Reason:       reason,
		ReviewerID:   reviewerID,
		DecidedAt:    time.Now(),
		Metadata: map[string]string{
			"duration_days": strconvItoa(durationDays),
			"default_days":  strconvItoa(s.warrantyDays),
		},
	}
}

func (s *ApprovalService) ReviewRework(resourceID int64, reworkCost float64, reviewerID int64, approve bool, comment string) ApprovalDecision {
	reason := strings.TrimSpace(comment)
	if reason == "" {
		reason = fmt.Sprintf("rework cost reviewed: %.2f", reworkCost)
	}
	return ApprovalDecision{
		ResourceType: "rework",
		ResourceID:   resourceID,
		NeedReview:   reworkCost > 0,
		Approved:     approve,
		Reason:       reason,
		ReviewerID:   reviewerID,
		DecidedAt:    time.Now(),
		Metadata: map[string]string{
			"rework_cost": fmt.Sprintf("%.2f", reworkCost),
		},
	}
}

func (s *ApprovalService) Validate(decision ApprovalDecision) error {
	if decision.ResourceType == "" {
		return fmt.Errorf("resource type is required")
	}
	if decision.ResourceID <= 0 {
		return fmt.Errorf("resource id is required")
	}
	if decision.NeedReview && decision.ReviewerID <= 0 {
		return fmt.Errorf("reviewer id is required")
	}
	if decision.NeedReview && strings.TrimSpace(decision.Reason) == "" {
		return fmt.Errorf("review reason is required")
	}
	return nil
}

func (s *ApprovalService) BuildAuditRecord(decision ApprovalDecision) repository.Record {
	return repository.Record{
		"action":        "APPROVAL_DECISION",
		"resource_type": decision.ResourceType,
		"resource_id":   fmt.Sprintf("%d", decision.ResourceID),
		"after_value": fmt.Sprintf(
			`{"need_review":%t,"approved":%t,"reason":%q,"reviewer_id":%d}`,
			decision.NeedReview,
			decision.Approved,
			decision.Reason,
			decision.ReviewerID,
		),
	}
}

func (s *ApprovalService) Summarize(decisions []ApprovalDecision) PendingApprovalSummary {
	summary := PendingApprovalSummary{
		ByResourceType: make(map[string]int64),
		Decisions:      make([]ApprovalDecision, 0, len(decisions)),
	}
	for _, decision := range decisions {
		summary.Total++
		summary.ByResourceType[decision.ResourceType]++
		summary.Decisions = append(summary.Decisions, decision)
	}
	return summary
}

func (s *ApprovalService) FilterPending(decisions []ApprovalDecision) []ApprovalDecision {
	items := make([]ApprovalDecision, 0, len(decisions))
	for _, decision := range decisions {
		if isPendingApproval(decision) {
			items = append(items, decision)
		}
	}
	return items
}

func (s *ApprovalService) CollectFromQuotationRecords(records []repository.Record) []ApprovalDecision {
	items := make([]ApprovalDecision, 0, len(records))
	for _, record := range records {
		amount := recordFloat(record, "total_amount")
		resourceID := recordInt64Value(record, "id")
		decision := s.ReviewQuotation(resourceID, amount, 0, false, "")
		if decision.NeedReview {
			items = append(items, decision)
		}
	}
	return items
}

func (s *ApprovalService) ReviewBatch(ctx context.Context, resourceType string, records []repository.Record, reviewerID int64, approve bool, comment string) ([]ApprovalDecision, error) {
	_ = ctx
	decisions := make([]ApprovalDecision, 0, len(records))
	for _, record := range records {
		resourceID := recordInt64Value(record, "id")
		switch resourceType {
		case "quotation":
			decisions = append(decisions, s.ReviewQuotation(resourceID, recordFloat(record, "total_amount"), reviewerID, approve, comment))
		case "warranty":
			decisions = append(decisions, s.ReviewWarranty(resourceID, int(recordInt64Value(record, "duration_days")), reviewerID, approve, comment))
		default:
			decisions = append(decisions, ApprovalDecision{
				ResourceType: resourceType,
				ResourceID:   resourceID,
				NeedReview:   true,
				Approved:     approve,
				Reason:       strings.TrimSpace(comment),
				ReviewerID:   reviewerID,
				DecidedAt:    time.Now(),
			})
		}
	}
	for _, decision := range decisions {
		if err := s.Validate(decision); err != nil {
			return nil, err
		}
	}
	return decisions, nil
}

func (s *ApprovalService) ReviewerWorkload(decisions []ApprovalDecision) map[int64]int64 {
	workload := make(map[int64]int64)
	for _, decision := range decisions {
		if decision.ReviewerID <= 0 {
			continue
		}
		workload[decision.ReviewerID]++
	}
	return workload
}

func (s *ApprovalService) NeedsReview(resourceType string, numericValue float64) bool {
	switch resourceType {
	case "quotation":
		return numericValue > s.quotationThreshold
	case "warranty":
		return int(numericValue) > s.warrantyDays
	case "rework":
		return numericValue > 0
	default:
		return false
	}
}

func recordFloat(record repository.Record, key string) float64 {
	switch value := record[key].(type) {
	case float64:
		return value
	case int64:
		return float64(value)
	case int:
		return float64(value)
	case string:
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return 0
}

func recordInt64Value(record repository.Record, key string) int64 {
	switch value := record[key].(type) {
	case int64:
		return value
	case int:
		return int64(value)
	case float64:
		return int64(value)
	case string:
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return 0
}

func strconvItoa(value int) string {
	return fmt.Sprintf("%d", value)
}
