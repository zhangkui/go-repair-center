package service

import "strings"

func normalizeIdempotencyKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "-")
	return strings.Trim(value, "-")
}
