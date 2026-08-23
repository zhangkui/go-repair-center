package service

import (
	"context"
	"go-repair-center/internal/repository"
)

type RepairExecutionService struct {
	*ResourceService
}

func NewRepairExecutionService(repo repository.CRUDRepository) *RepairExecutionService {
	resource := NewResourceService(repo, "repair_execution", "repair_order_id", "technician_id").LimitUpdates(
		"technician_id",
		"fault_code_id",
		"created_at",
		"started_at",
		"completed_at",
		"status",
	)
	return &RepairExecutionService{ResourceService: resource}
}

func (s *RepairExecutionService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
	return s.List(ctx, page, pageSize, keyword, "")
}

func (s *RepairExecutionService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
	return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *RepairExecutionService) Archive(ctx context.Context, id int64) error {
	return s.Delete(ctx, id)
}

var repair_executionTransitions = map[string][]string{"PENDING": {"IN_PROGRESS"}, "IN_PROGRESS": {"COMPLETED", "WAITING_TEST"}, "COMPLETED": {"WAITING_TEST"}, "WAITING_TEST": {"TESTED"}, "TESTED": {"WAITING_PICKUP"}, "WAITING_PICKUP": {"DELIVERED"}}

func (s *RepairExecutionService) ChangeStatus(ctx context.Context, id int64, current, next, reason string) (repository.Record, error) {
	if err := ValidateTransition(current, next, repair_executionTransitions); err != nil {
		return nil, err
	}
	return s.Update(ctx, id, repository.Record{"status": next, "status_reason": reason})
}
