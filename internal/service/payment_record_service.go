package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type PaymentRecordService struct {
    *ResourceService
}

func NewPaymentRecordService(repo repository.CRUDRepository) *PaymentRecordService {
    return &PaymentRecordService{ResourceService: NewResourceService(repo, "payment_record", "repair_order_id","amount","payment_method")}
}

func (s *PaymentRecordService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *PaymentRecordService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *PaymentRecordService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
