package objects

import (
	"fmt"
	"os"
)

func UpdateIndex(indexPath, hash, filePath string) error {
	indexFile, err := os.OpenFile(indexPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open index: %v", err)
	}
	defer indexFile.Close()


	indexEntry := fmt.Sprintf("%s %s\n", hash, filePath)
	if _, err := indexFile.WriteString(indexEntry); err != nil {
		return fmt.Errorf("failed to write to index: %v", err)
	}

	return nil
}

