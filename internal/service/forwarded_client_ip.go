package service

import (
	"net"
	"strings"
)

func forwardedClientAddress(value string) string {
	parts := strings.Split(value, ",")
	for index := len(parts) - 1; index >= 0; index-- {
		candidate := strings.TrimSpace(parts[index])
		if candidate == "" {
			continue
		}
		if net.ParseIP(candidate) == nil {
			continue
		}
		return candidate
	}
	return ""
}
