package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type CustomerService struct {
    *ResourceService
}

func NewCustomerService(repo repository.CRUDRepository) *CustomerService {
    return &CustomerService{ResourceService: NewResourceService(repo, "customer", "name")}
}

func (s *CustomerService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *CustomerService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *CustomerService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
