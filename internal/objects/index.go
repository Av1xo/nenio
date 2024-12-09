package objects

import (
	"os"
)

type IndexEntry struct {
	Path string
	BlobHash string
	FileMode os.FileMode
	FileSize int64
	ModifiedAt int64
}

type Index struct {
	Entries map[string]IndexEntry
}
