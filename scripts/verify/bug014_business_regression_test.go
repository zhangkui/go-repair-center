package verify_test

import (
	"errors"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug014_BusinessRegression(t *testing.T) {
	releaseErr := errors.New("release failed")
	if got := service.CombineLockErrors(nil, releaseErr); !errors.Is(got, releaseErr) {
		t.Fatalf("successful operation must return release error, got %v", got)
	}
	operationErr := errors.New("operation failed")
	if got := service.CombineLockErrors(operationErr, releaseErr); !errors.Is(got, operationErr) {
		t.Fatalf("operation error must take precedence, got %v", got)
	}
}
