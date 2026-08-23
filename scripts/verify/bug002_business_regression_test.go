package verify_test

import (
	"context"
	"testing"

	"go-repair-center/internal/repository"
	"go-repair-center/internal/service"
)

type bug002Repository struct {
	updated repository.Record
}

func (r *bug002Repository) Create(context.Context, repository.Record) (repository.Record, error) {
	return nil, nil
}
func (r *bug002Repository) Get(context.Context, int64) (repository.Record, error) { return nil, nil }
func (r *bug002Repository) List(context.Context, int, int, string, string) (repository.Page, error) {
	return repository.Page{}, nil
}
func (r *bug002Repository) Update(_ context.Context, _ int64, input repository.Record) (repository.Record, error) {
	r.updated = input
	return input, nil
}
func (r *bug002Repository) Delete(context.Context, int64) error { return nil }

func TestBug002_BusinessRegression(t *testing.T) {
	repo := &bug002Repository{}
	orders := service.NewRepairOrderService(repo)
	updated, err := orders.ChangeStatus(context.Background(), 1002, "DISPATCHED", "CANCELLED", "customer cancelled")
	if err != nil {
		t.Fatalf("cancelling a dispatched repair order should succeed: %v", err)
	}
	if updated["status"] != "CANCELLED" || repo.updated["status"] != "CANCELLED" {
		t.Fatalf("expected persisted status CANCELLED, got result=%v persisted=%v", updated["status"], repo.updated["status"])
	}
}
