package verify_test

import (
	"net/http"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug007_BusinessRegression(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("X-Forwarded-For", "198.51.100.10, 10.0.0.2")
	request.RemoteAddr = "10.0.0.3:1234"
	if got := service.NewSecurityService(service.PasswordPolicy{}).ClientIP(request); got != "198.51.100.10" {
		t.Fatalf("client IP = %q, want 198.51.100.10", got)
	}
}
