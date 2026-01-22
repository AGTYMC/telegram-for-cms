package storage

import (
	"path/filepath"
	"strings"
)

const (
	SessionsDir = "sessions"
)

func SessionFilePath(phone string) string {
	clean := strings.ReplaceAll(phone, "+", "")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "-", "")
	return filepath.Join(SessionsDir, clean+".session")
}
