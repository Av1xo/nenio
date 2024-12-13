package objects

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateBlob(t *testing.T) {
	tempDir := t.TempDir()
	content := []byte("This is a test blob")

	hash, _, err := CreateBlob(tempDir, content)
	if err != nil {
		t.Fatalf("CreateBlob failed: %v", err)
	}

	blobPath := filepath.Join(tempDir, hash[:2], hash[2:])
	if _, err := os.Stat(blobPath); os.IsNotExist(err) {
		t.Errorf("Blob file does not exist at %s", blobPath)
	}
}

func TestReadBlob(t *testing.T) {
	tempDir := t.TempDir()
	content := []byte("Another test blob")

	hash, _, err := CreateBlob(tempDir, content)
	if err != nil {
		t.Fatalf("CreateBlob failed: %v", err)
	}

	readContent, err := ReadBlob(tempDir, hash)
	if err != nil {
		t.Fatalf("ReadBlob failed: %v", err)
	}

	if !bytes.Equal(content, readContent) {
		t.Errorf("ReadBlob content mismatch. Got %s, want %s", string(readContent), string(content))
	}
}

func TestCompressAndDecompressBlob(t *testing.T) {
	content := []byte("Test data for compression and decompression")

	compressed, err := CompressBlob(content)
	if err != nil {
		t.Fatalf("CompressBlob failed: %v", err)
	}

	reader := bytes.NewReader(compressed)
	decompressed, err := DecompressBlob(reader)
	if err != nil {
		t.Fatalf("DecompressBlob failed: %v", err)
	}

	if !bytes.Equal(content, decompressed) {
		t.Errorf("Decompressed content mismatch. Got %s, want %s", string(decompressed), string(content))
	}
}

func TestBlobExists(t *testing.T) {
	tempDir := t.TempDir()
	content := []byte("Check if blob exists")

	hash, _, err := CreateBlob(tempDir, content)
	if err != nil {
		t.Fatalf("CreateBlob failed: %v", err)
	}

	exists, err := BlobExists(tempDir, hash)
	if err != nil {
		t.Fatalf("BlobExists failed: %v", err)
	}

	if !exists {
		t.Errorf("BlobExists returned false for an existing blob")
	}

	nonExistentHash := "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	exists, err = BlobExists(tempDir, nonExistentHash)
	if err != nil {
		t.Fatalf("BlobExists failed: %v", err)
	}

	if exists {
		t.Errorf("BlobExists returned true for a non-existent blob")
	}
}
