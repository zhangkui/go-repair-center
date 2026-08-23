package verify_test

import (
	"context"
	"errors"
	"testing"

	"go-repair-center/internal/repository"
	"go-repair-center/internal/service"
)

type bug003Repository struct {
	record      repository.Record
	updateCalls int
	lastUpdate  repository.Record
}

func newBug003Repository() *bug003Repository {
	return &bug003Repository{record: repository.Record{"id": int64(42), "repair_order_id": int64(7), "technician_id": int64(9), "diagnosis": "initial diagnosis"}}
}
func (r *bug003Repository) Create(context.Context, repository.Record) (repository.Record, error) { return nil, errors.New("unexpected create") }
func (r *bug003Repository) Get(_ context.Context, id int64) (repository.Record, error) { if id != 42 { return nil, errors.New("record not found") }; return cloneBug003Record(r.record), nil }
func (r *bug003Repository) List(context.Context, int, int, string, string) (repository.Page, error) { return repository.Page{}, errors.New("unexpected list") }
func (r *bug003Repository) Update(_ context.Context, id int64, input repository.Record) (repository.Record, error) { if id != 42 { return nil, errors.New("record not found") }; r.updateCalls++; r.lastUpdate = cloneBug003Record(input); for field, value := range input { r.record[field] = value }; return cloneBug003Record(r.record), nil }
func (r *bug003Repository) Delete(context.Context, int64) error { return errors.New("unexpected delete") }
func cloneBug003Record(input repository.Record) repository.Record { output := make(repository.Record, len(input)); for field, value := range input { output[field] = value }; return output }

func TestBug003_BusinessRegression(t *testing.T) {
	t.Run("rejects an update with no writable business fields", func(t *testing.T) {
		repo := newBug003Repository()
		_, err := service.NewRepairExecutionService(repo).Update(context.Background(), 42, repository.Record{"created_at": "2026-08-23 10:00:00"})
		if !errors.Is(err, service.ErrValidation) { t.Fatalf("expected validation error for non-writable fields, got %v", err) }
		if repo.updateCalls != 0 { t.Fatalf("repository update must not run for rejected input, got %d calls", repo.updateCalls) }
	})
	t.Run("keeps diagnosis as a supported update", func(t *testing.T) {
		repo := newBug003Repository()
		updated, err := service.NewRepairExecutionService(repo).Update(context.Background(), 42, repository.Record{"diagnosis": "power module failure confirmed"})
		if err != nil { t.Fatalf("expected diagnosis update to succeed, got %v", err) }
		if repo.updateCalls != 1 || repo.lastUpdate["diagnosis"] != "power module failure confirmed" || updated["diagnosis"] != "power module failure confirmed" { t.Fatalf("diagnosis update was not preserved: calls=%d update=%v result=%v", repo.updateCalls, repo.lastUpdate, updated) }
	})
}