package handler

import (
	"strings"
)

func extractID(path string, prefix string) (string, bool) {
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	id := strings.TrimPrefix(path, prefix)
	id = strings.Trim(id, "/")
	if id == "" {
		return "", false
	}
	// não aceitar barras dentro do id
	if strings.Contains(id, "/") {
		return "", false
	}
	return id, true
}
