package objects

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type IndexEntry struct {
	Path       string	`json:"path"`
	BlobHash   string	`json:"blob_hash"`
	BlobData []byte `json:"blob_data"`
	FileMode   os.FileMode `json:"file_mode"`
	FileSize   int64	`json:"file_size"`
	ModifiedAt int64	`json:"modified_at"`
}

type Index struct {
	Entries map[string]IndexEntry `json:"entries"`
	UpdatedAt int64               `json:"updated_at"`
}

func (index *Index) AddEntry(entry IndexEntry) {
	if index.Entries == nil {
		index.Entries = make(map[string]IndexEntry)
	}
	index.Entries[entry.Path] = entry
}

func (index *Index) RemoveEntry(path string) {
	delete(index.Entries, path)
}

func NewIndex() *Index {
	return &Index{
		Entries: make(map[string]IndexEntry),
		UpdatedAt: time.Now().Unix(),
	}
}

// The LoadIndex function reads and parses an index file, returning the index data or an error.
func LoadIndex(indexPath string) (*Index, error) {
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return NewIndex(), nil
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read index file: %v", err)
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index file: %v", err)
	}

	return &index, nil
}

// The function `SaveIndex` serializes an index to JSON format and saves it to a specified file path.
func SaveIndex(indexPath string, index *Index) error {
	data, err := json.MarshalIndent(index, "", " ")
	if err != nil {
		return fmt.Errorf("failed to serialize index: %v", err)
	}

	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save index: %v", err)
	}

	return nil
}

// The `AddToIndex` function concurrently adds files to an index, handling file content updates and
// creating blobs as needed.
func AddToIndex(objectsDir, indexPath string, ignorePatterns, files []string) error {
	index, err := LoadIndex(indexPath)
	if err != nil {
		return err
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		errChan = make(chan error, len(files))
	)

	for _, file := range files {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			if ShouldIgnoreFile(file, ignorePatterns) {
				return
			}

			info, err := os.Stat(file)
			if err != nil {
				errChan <- fmt.Errorf("failed to stat file %s: %v", file, err)
				return
			}

			if info.IsDir() {
				return
			}

			content, err := os.ReadFile(file)
			if err != nil {
				errChan <- fmt.Errorf("failed to read file %s: %v", file, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			existingEntry, exists := index.Entries[file]
			if exists {
				savedContent, err := ReadBlob(objectsDir, existingEntry.BlobHash)
				if err != nil {
					errChan <- fmt.Errorf("failed to read existing blob for file %s: %v", file, err)
					return
				}

				delta, err := ComputeDelta(savedContent, content)
				if err != nil {
					errChan <- fmt.Errorf("failed to compute delta for file %s: %v", file, err)
					return
				}

				if len(delta) > 0 {
					blobHash, blobData, err := CreateBlob(objectsDir, content)
					if err != nil {
						errChan <- fmt.Errorf("failed to create blob for file %s: %v", file, err)
						return
					}

					existingEntry.BlobHash = blobHash
					existingEntry.BlobData = blobData
					existingEntry.FileMode = info.Mode()
					existingEntry.FileSize = info.Size()
					existingEntry.ModifiedAt = info.ModTime().Unix()
					index.Entries[file] = existingEntry
				}
			} else {
				blobHash, blobData, err := CreateBlob(objectsDir, content)
				if err != nil {
					errChan <- fmt.Errorf("failed to create blob for file %s: %v", file, err)
					return
				}

				entry := IndexEntry{
					Path:       file,
					BlobHash:   blobHash,
					BlobData: blobData,
					FileMode:   info.Mode(),
					FileSize:   info.Size(),
					ModifiedAt: info.ModTime().Unix(),
				}
				index.AddEntry(entry)
			}
		}(file)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	if err := SaveIndex(indexPath, index); err != nil {
		return fmt.Errorf("failed to save index: %v", err)
	}

	return nil
}

