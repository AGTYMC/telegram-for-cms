package storage

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	SessionsDir = "sessions"
)

func SessionFilePath(phone string) string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Dir(exe)

	clean := strings.ReplaceAll(phone, "+", "")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "-", "")
	return filepath.Join(dir, SessionsDir, clean+".session")
}
