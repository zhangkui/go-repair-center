package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type PasswordPolicy struct {
	MinLength       int  `json:"min_length"`
	RequireUpper    bool `json:"require_upper"`
	RequireDigit    bool `json:"require_digit"`
	RequireSpecial  bool `json:"require_special"`
}

type SecurityService struct {
	policy PasswordPolicy
}

type RiskAssessment struct {
	Score     int      `json:"score"`
	Level     string   `json:"level"`
	Reasons   []string `json:"reasons"`
	IPAddress string   `json:"ip_address"`
}

func NewSecurityService(policy PasswordPolicy) *SecurityService {
	return &SecurityService{policy: policy}
}

func (s *SecurityService) ValidatePassword(value string) error {
	if len(value) < s.policy.MinLength {
		return fmt.Errorf("password must contain at least %d characters", s.policy.MinLength)
	}
	hasUpper := false
	hasDigit := false
	hasSpecial := false
	for _, item := range value {
		switch {
		case item >= 'A' && item <= 'Z':
			hasUpper = true
		case item >= '0' && item <= '9':
			hasDigit = true
		case !(item >= 'a' && item <= 'z') && !(item >= 'A' && item <= 'Z') && !(item >= '0' && item <= '9'):
			hasSpecial = true
		}
	}
	if s.policy.RequireUpper && !hasUpper {
		return fmt.Errorf("password must include an uppercase letter")
	}
	if s.policy.RequireDigit && !hasDigit {
		return fmt.Errorf("password must include a digit")
	}
	if s.policy.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password must include a special character")
	}
	return nil
}

func (s *SecurityService) HashSensitiveValue(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func (s *SecurityService) MaskPhone(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 4 {
		return value
	}
	return value[:3] + strings.Repeat("*", len(value)-5) + value[len(value)-2:]
}

func (s *SecurityService) MaskEmail(value string) string {
	value = strings.TrimSpace(value)
	parts := strings.Split(value, "@")
	if len(parts) != 2 || len(parts[0]) < 2 {
		return value
	}
	return parts[0][:1] + "***@" + parts[1]
}

func (s *SecurityService) NormalizeUsername(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *SecurityService) NormalizeIP(value string) string {
	host, _, err := net.SplitHostPort(value)
	if err == nil {
		value = host
	}
	return strings.TrimSpace(value)
}

func (s *SecurityService) ClientIP(r *http.Request) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP"} {
		value := strings.TrimSpace(r.Header.Get(header))
		if value == "" {
			continue
		}
		if strings.Contains(value, ",") {
			value = strings.TrimSpace(strings.Split(value, ",")[0])
		}
		return s.NormalizeIP(value)
	}
	return s.NormalizeIP(r.RemoteAddr)
}

func (s *SecurityService) AssessRisk(r *http.Request, username string) RiskAssessment {
	reasons := make([]string, 0, 4)
	score := 0
	ip := s.ClientIP(r)

	if username == "" {
		score += 20
		reasons = append(reasons, "empty username")
	}
	if strings.Contains(r.UserAgent(), "curl") {
		score += 10
		reasons = append(reasons, "curl user-agent")
	}
	if r.Header.Get("X-Forwarded-For") == "" {
		score += 5
		reasons = append(reasons, "no proxy forwarding headers")
	}
	if r.Method != http.MethodGet && r.Header.Get("Content-Type") == "" {
		score += 10
		reasons = append(reasons, "missing content-type")
	}
	if strings.HasPrefix(ip, "127.") || ip == "::1" {
		reasons = append(reasons, "local network")
	}

	level := "LOW"
	switch {
	case score >= 30:
		level = "HIGH"
	case score >= 15:
		level = "MEDIUM"
	}

	return RiskAssessment{
		Score:     score,
		Level:     level,
		Reasons:   reasons,
		IPAddress: ip,
	}
}

func (s *SecurityService) BuildLoginAttemptKey(ip string, username string) string {
	return "login:attempt:" + s.NormalizeIP(ip) + ":" + s.NormalizeUsername(username)
}

func (s *SecurityService) BuildLockKey(ip string) string {
	return "login:lock:" + s.NormalizeIP(ip)
}

func (s *SecurityService) BuildAuditFingerprint(username string, when time.Time) string {
	payload := fmt.Sprintf("%s|%s", s.NormalizeUsername(username), when.UTC().Format(time.RFC3339))
	return s.HashSensitiveValue(payload)
}

func (s *SecurityService) RedactMap(values map[string]string) map[string]string {
	redacted := make(map[string]string, len(values))
	for key, value := range values {
		lowerKey := strings.ToLower(key)
		switch {
		case strings.Contains(lowerKey, "password"):
			redacted[key] = "***"
		case strings.Contains(lowerKey, "phone"):
			redacted[key] = s.MaskPhone(value)
		case strings.Contains(lowerKey, "email"):
			redacted[key] = s.MaskEmail(value)
		default:
			redacted[key] = value
		}
	}
	return redacted
}

func (s *SecurityService) PasswordPolicy() PasswordPolicy {
	return s.policy
}
