package verify_test

import (
	"errors"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug009_BusinessRegression(t *testing.T) {
	if !service.CacheReadFailureIsMiss(nil) {
		t.Fatal("nil cache result should represent a miss")
	}
	decodeErr := errors.New("cache payload invalid")
	if service.CacheReadFailureIsMiss(decodeErr) {
		t.Fatal("cache decode failure must not be treated as a normal cache miss")
	}
	transportErr := errors.New("cache connection refused")
	if service.CacheReadFailureIsMiss(transportErr) {
		t.Fatal("cache transport failure must remain observable")
	}
}
