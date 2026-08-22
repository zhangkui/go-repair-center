package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type RepairTestRecordService struct {
    *ResourceService
}

func NewRepairTestRecordService(repo repository.CRUDRepository) *RepairTestRecordService {
    return &RepairTestRecordService{ResourceService: NewResourceService(repo, "repair_test_record", "repair_execution_id","test_item")}
}

func (s *RepairTestRecordService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *RepairTestRecordService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *RepairTestRecordService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
