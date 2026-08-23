package service

import "strings"

func forwardedClientAddress(value string) string {
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[len(parts)-1])
}
