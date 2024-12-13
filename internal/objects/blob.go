package objects

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"lukechampine.com/blake3"
)

func CreateBlob(objectsDir string, content []byte) (string, []byte, error) {
	hash := blake3.Sum256(content)
	hashHex := fmt.Sprintf("%x", hash)

	subDir := filepath.Join(objectsDir, hashHex[:2])
	if err := os.MkdirAll(subDir, 0755); err != nil {
		return "", nil, fmt.Errorf("failed to create objects subdirectory: %v", err)
	}

	if exists, err := BlobExists(subDir, hashHex); err != nil {
		return "", nil, err
	} else if exists {
		return hashHex, nil, nil
	}

	compressedContent, err := CompressBlob(content)
	if err != nil {
		return "", nil, err
	}


	return hashHex, compressedContent, nil
}

func ReadBlob(objectsDir, hashHex string) ([]byte, error) {
	file, err := os.Open(blobPath(objectsDir, hashHex))
	if err != nil {
		return nil, fmt.Errorf("failed to open blob: %v", err)
	}
	defer file.Close()

	return DecompressBlob(file)
}

func CompressBlob(content []byte) ([]byte, error) {
	var compressedBuffer bytes.Buffer

	gzipWriter := gzip.NewWriter(&compressedBuffer)

	if _, err := gzipWriter.Write(content); err != nil {
		return nil, fmt.Errorf("failed to compress content: %v", err)
	}

	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}

	return compressedBuffer.Bytes(), nil
}

func DecompressBlob(file io.Reader) ([]byte, error) {
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	var decompressedBuffer bytes.Buffer
	if _, err := io.Copy(&decompressedBuffer, gzipReader); err != nil {
		return nil, fmt.Errorf("failed to decompress blob: %v", err)
	}

	return decompressedBuffer.Bytes(), nil
}

func BlobExists(objectsDir, hashHex string) (bool, error) {
	_, err := os.Stat(blobPath(objectsDir, hashHex))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("failed to check blob existence: %v", err)
}

func DeleteBlob(objectsDir, hashHex string) error {
	if err := os.Remove(blobPath(objectsDir, hashHex)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("blob %s does not exist", hashHex)
		}
		return fmt.Errorf("failed to delete blob: %w", err)
	}
	return nil
}

func blobPath(objectsDir, hashHex string) string {
	return filepath.Join(objectsDir, hashHex[:2], hashHex[2:])
}
