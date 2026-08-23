package verify_test

import (
	"reflect"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug004_BusinessRegression(t *testing.T) {
	security := service.NewSecurityService(service.PasswordPolicy{})

	if got := security.RedactMap(nil); got != nil {
		t.Fatalf("nil audit metadata must remain nil, got %#v", got)
	}

	empty := map[string]string{}
	if got := security.RedactMap(empty); got == nil || len(got) != 0 {
		t.Fatalf("present but empty metadata must remain a non-nil empty map, got %#v", got)
	}

	input := map[string]string{
		"password":      "Secret123!",
		"contact_phone": "13800138000",
		"contact_email": "operator@example.com",
		"note":          "keep",
	}
	want := map[string]string{
		"password":      "***",
		"contact_phone": "138******00",
		"contact_email": "o***@example.com",
		"note":          "keep",
	}
	if got := security.RedactMap(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("redacted metadata = %#v, want %#v", got, want)
	}
}
