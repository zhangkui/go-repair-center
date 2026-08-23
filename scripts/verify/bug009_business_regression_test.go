package verify_test

import (
	"errors"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug009_BusinessRegression(t *testing.T) {
	if service.CacheReadFailureIsMiss(errors.New("cache payload invalid")) {
		t.Fatal("cache decode failure must not be treated as a normal cache miss")
	}
}
