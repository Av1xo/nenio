package objects

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type IndexEntry struct {
	Path       string
	BlobHash   string
	FileMode   os.FileMode
	FileSize   int64
	ModifiedAt int64
}

type Index struct {
	Entries map[string]IndexEntry
}

func (index *Index) AddEntry(entry IndexEntry) {
	index.Entries[entry.Path] = entry
}

func (index *Index) RemoveEntry(path string) {
	delete(index.Entries, path)
}

func LoadIndex(indexPath string) (*Index, error) {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Index{
				Entries: make(map[string]IndexEntry),
			}, nil
		}
		return nil, fmt.Errorf("failed to load index: %v", err)
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %v", err)
	}

	return &index, nil
}

func SaveIndex(indexPath string, index *Index) error {
	data, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to serialize index: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return fmt.Errorf("failed to create index directory: %v", err)
	}

	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save index: %v", err)
	}

	return nil
}

func AddToIndex(objectsDir, indexPath string, files []string) error {
	index, err := LoadIndex(indexPath)
	if err != nil {
		return err
	}

	ignorePatterns, _ := LoadIgnorePatterns(".nignore")

	for _, file := range files {
		if ShouldIgnoreFile(file, ignorePatterns) {
			continue
		}

		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %v", file, err)
		}

		if info.IsDir() {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", file, err)
		}

		existingEntry, exists := index.Entries[file]
		if exists {
			savedContent, err := ReadBlob(objectsDir, existingEntry.BlobHash)
			if err != nil {
				return fmt.Errorf("failed to  read existing blob for file %s: %v", file, err)
			}

			delta, err := ComputeDelta(savedContent, content)
			if err != nil {
				return fmt.Errorf("failed to compute delta for file %s: %v", file, err)
			}
			if len(delta) > 0 {
				blobHash, err := CreateBlob(objectsDir, content)
				if err != nil {
					return fmt.Errorf("failed to create blob for file %s: %v", file, err)
				}

				existingEntry.BlobHash = blobHash
				existingEntry.FileMode = info.Mode()
				existingEntry.FileSize = info.Size()
				existingEntry.ModifiedAt = info.ModTime().Unix()
				index.Entries[file] = existingEntry
			}

		} else {
			blobHash, err := CreateBlob(objectsDir, content)
			if err != nil {
				return fmt.Errorf("failed to create blob for file %s: %v", objectsDir, err)
			}

			entry := IndexEntry{
				Path:       file,
				BlobHash:   blobHash,
				FileMode:   info.Mode(),
				FileSize:   info.Size(),
				ModifiedAt: info.ModTime().Unix(),
			}

			index.AddEntry(entry)
		}

	}

	return SaveIndex(indexPath, index)
}

func LoadIgnorePatterns(ignoreFilePath string) ([]string, error) {
	data, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read .nignore: %v", err)
	}

	patterns := strings.Split(string(data), "\n")
	for i, p := range patterns {
		patterns[i] = strings.TrimSpace(p)
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
