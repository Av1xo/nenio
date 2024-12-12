package objects

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadIgnorePatterns(ignoreFilePath string) ([]string, error) {
	data, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read .nignore: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var patterns []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			patterns = append(patterns, trimmed)
		}
	}

	return patterns, nil
}

func ShouldIgnoreFile(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}