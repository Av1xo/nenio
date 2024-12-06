package objects

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type IndexEntry struct {
	Mode string
	Hash string
	FilePath string
}

func parseIndex(indexPath string) ([]IndexEntry, error) {
	file, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []IndexEntry{}, nil
		}
		return nil, fmt.Errorf("coul not open index file: %v", err)
	}
	defer file.Close()

	var entries []IndexEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid index entry: %s", line)
		}

		entry := IndexEntry{
			Mode: parts[0],
			Hash: parts[1],
			FilePath: parts[2],
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading index file: %v", err)
	}

	return entries, nil
}

func UpdateIndex(indexPath string, newEntry IndexEntry) error {
	entries, err := parseIndex(indexPath)
	if err != nil {
		return fmt.Errorf("filed to read index: %v", err)
	}

	for _, entry := range entries {
		if entry.FilePath == newEntry.FilePath {
			return fmt.Errorf("file %s is already staged", newEntry.FilePath)
		}
	}

	indexFile, err := os.OpenFile(indexPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open index: %v", err)
	}
	defer indexFile.Close()


	indexEntry := fmt.Sprintf("%s %s %s\n", newEntry.Mode, newEntry.Hash, newEntry.FilePath)
	if _, err := indexFile.WriteString(indexEntry); err != nil {
		return fmt.Errorf("failed to write to index: %v", err)
	}

	return nil
}

