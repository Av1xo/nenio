package objects

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The function `LoadIgnorePatterns` reads ignore patterns from a file and returns them as a slice of
// strings.
func LoadIgnorePatterns(ignoreFilePath string) ([]string, error) {
	data, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{".nenio/"}, nil
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
	patterns = append(patterns, ".nenio/")

	return patterns, nil
}

// The function `ShouldIgnoreFile` determines whether a file path should be ignored based on a list of
// patterns.
func ShouldIgnoreFile(path string, patterns []string) bool {
	relPath, err := filepath.Rel(".", path)
	if err != nil {
		relPath = path
	}
	relPath = filepath.ToSlash(relPath)

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if strings.HasPrefix(pattern, "!") {
			excludePattern := strings.TrimPrefix(pattern, "!")
			if matched, _ := filepath.Match(excludePattern, filepath.Base(path)); matched {
				return false
			}
		}
	}

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}

		if strings.HasPrefix(pattern, "/") {
			absolutePattern := strings.TrimPrefix(pattern, "/")
			if strings.HasPrefix(relPath, absolutePattern) {
				return true
			}
		}

		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(relPath, dirPattern+"/") || relPath == dirPattern {
				return true
			}
		}

		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}

	return false
}
