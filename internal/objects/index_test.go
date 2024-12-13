package objects

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadIndex_NotExist(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index.json")

	index, err := LoadIndex(indexPath)

	assert.NoError(t, err)
	assert.NotNil(t, index)
	assert.Empty(t, index.Entries)
}

func TestSaveAndLoadIndex(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index.json")

	index := &Index{
		Entries: map[string]IndexEntry{
			"file1.txt": {
				Path:       "file1.txt",
				BlobHash:   "hash1",
				FileMode:   0644,
				FileSize:   123,
				ModifiedAt: 1623840123,
			},
		},
	}

	err := SaveIndex(indexPath, index)
	assert.NoError(t, err)

	loadedIndex, err := LoadIndex(indexPath)
	assert.NoError(t, err)
	assert.Equal(t, index, loadedIndex)
}

func TestAddToIndex(t *testing.T) {
	tempDir := t.TempDir()
	objectsDir := filepath.Join(tempDir, "objects")
	indexPath := filepath.Join(tempDir, "index.json")

	err := os.Mkdir(objectsDir, 0755)
	assert.NoError(t, err)

	file1Path := filepath.Join(tempDir, "file1.txt")
	err = os.WriteFile(file1Path, []byte("content1"), 0644)
	assert.NoError(t, err)

	err = AddToIndex(objectsDir, indexPath, []string{file1Path})
	assert.NoError(t, err)

	index, err := LoadIndex(indexPath)
	assert.NoError(t, err)
	assert.Contains(t, index.Entries, file1Path)
	assert.NotEmpty(t, index.Entries[file1Path].BlobHash)
}

func TestShouldIgnoreFile_AdvancedPatterns(t *testing.T) {
	patterns := []string{
		"*.log",         
		"node_modules/", 
		"# This is a comment", 
		"!important.log",
		"/logs/",
	}

	assert.True(t, ShouldIgnoreFile("./debug.log", patterns))
	assert.True(t, ShouldIgnoreFile("./node_modules/file.js", patterns))
	assert.True(t, ShouldIgnoreFile("./node_modules/subdir/file.js", patterns))
	assert.False(t, ShouldIgnoreFile("./important.log", patterns))
	assert.False(t, ShouldIgnoreFile("./random/logs/file.log", patterns))
	assert.True(t, ShouldIgnoreFile("./logs/file.log", patterns))
	assert.False(t, ShouldIgnoreFile("./not_logs/debug.log", patterns))
}



func TestLoadIgnorePatterns(t *testing.T) {
	tempDir := t.TempDir()
	ignoreFilePath := filepath.Join(tempDir, ".nignore")

	ignoreContent := "*.tmp\nignore_me\n"
	err := os.WriteFile(ignoreFilePath, []byte(ignoreContent), 0644)
	assert.NoError(t, err)

	patterns, err := LoadIgnorePatterns(ignoreFilePath)
	assert.NoError(t, err)
	assert.Equal(t, []string{"*.tmp", "ignore_me"}, patterns)
}

func TestAddToIndex_IgnoreFile(t *testing.T) {
	tempDir := t.TempDir()
	objectsDir := filepath.Join(tempDir, "objects")
	indexPath := filepath.Join(tempDir, "index.json")
	ignoreFilePath := filepath.Join(tempDir, ".nignore")

	err := os.Mkdir(objectsDir, 0755)
	assert.NoError(t, err)

	ignoreContent := "*.tmp\n"
	err = os.WriteFile(ignoreFilePath, []byte(ignoreContent), 0644)
	assert.NoError(t, err)

	file1Path := filepath.Join(tempDir, "file1.tmp")
	err = os.WriteFile(file1Path, []byte("temporary content"), 0644)
	assert.NoError(t, err)

	err = AddToIndex(objectsDir, indexPath, []string{file1Path})
	assert.NoError(t, err)

	index, err := LoadIndex(indexPath)
	assert.NoError(t, err)
	assert.Empty(t, index.Entries)
}
