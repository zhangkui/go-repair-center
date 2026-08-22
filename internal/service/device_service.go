package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type DeviceService struct {
    *ResourceService
}

func NewDeviceService(repo repository.CRUDRepository) *DeviceService {
    return &DeviceService{ResourceService: NewResourceService(repo, "device", "customer_id","serial_number","brand")}
}

func (s *DeviceService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *DeviceService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *DeviceService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}

var deviceTransitions = map[string][]string{"NORMAL":{"REPORTED","SCRAPPED"},"REPORTED":{"REPAIRING","NORMAL"},"REPAIRING":{"READY_FOR_PICKUP","SCRAPPED"},"READY_FOR_PICKUP":{"PICKED_UP"}}

func (s *DeviceService) ChangeStatus(ctx context.Context, id int64, current, next, reason string) (repository.Record, error) {
    if err := ValidateTransition(current, next, deviceTransitions); err != nil { return nil, err }
    return s.Update(ctx, id, repository.Record{"status": next, "status_reason": reason})
}