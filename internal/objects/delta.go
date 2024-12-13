package objects

import (
	"bytes"
	"fmt"
	"sync"

	"lukechampine.com/blake3"
)

/*
Дельта алгоритм
	Розділити файл на блоки
	Видобути хеш кожного блоку
	Порівняти хеші
		хеші рівні блок не змінювався
		інакше зберегти як частину дельти
	зберегти дельту як набір інструкцій
		COPY(offset, length) скопіювати з базового файлу
		ADD(data) додати нові дані
*/

type instructionType string

const (
	ADD  instructionType = "ADD"
	COPY instructionType = "COPY"
)

const BlockSize = 4096 // 4kb

type DeltaInstruction struct {
	Type   instructionType
	Offset int
	Length int
	Data   []byte
}

// The ComputeDelta function generates a list of delta instructions to transform a base byte slice into
// an updated byte slice using block-based differencing.
func ComputeDelta(base, updated []byte) ([]DeltaInstruction, error) {
	if base == nil || updated == nil {
		return nil, fmt.Errorf("base or updated data is nil")
	}

	var delta []DeltaInstruction
	baseHashes := make(map[[32]byte]int)

	var wg sync.WaitGroup
	mu := sync.Mutex{}

	for i := 0; i < len(base); i += BlockSize {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			end := i + BlockSize
			if end > len(base) {
				end = len(base)
			}
			hash := blake3.Sum256(base[i:end])

			mu.Lock()
			baseHashes[hash] = i
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	var buffer bytes.Buffer
	for i := 0; i < len(updated); i += BlockSize {
		end := i + BlockSize
		if end > len(updated) {
			end = len(updated)
		}
		block := updated[i:end]
		hash := blake3.Sum256(block)

		if offset, found := baseHashes[hash]; found {
			if buffer.Len() > 0 {
				delta = append(delta, DeltaInstruction{
					Type: ADD,
					Data: buffer.Bytes(),
				})
				buffer.Reset()
			}

			delta = append(delta, DeltaInstruction{
				Type:   COPY,
				Offset: offset,
				Length: end - i,
			})
		} else {
			buffer.Write(block)
		}
	}

	if buffer.Len() > 0 {
		delta = append(delta, DeltaInstruction{
			Type: ADD,
			Data: buffer.Bytes(),
		})
	}

	return delta, nil
}

// The ApplyDelta function takes a base byte slice and a list of DeltaInstructions to apply changes and
// return the resulting byte slice.
func ApplyDelta(base []byte, delta []DeltaInstruction) ([]byte, error) {
	if base == nil {
		return nil, fmt.Errorf("base data is nil")
	}

	var result bytes.Buffer
	for _, instruction := range delta {
		switch instruction.Type {
		case COPY:
			start := instruction.Offset
			end := start + instruction.Length
			if start < 0 || start >= len(base) {
				return nil, fmt.Errorf("COPY instruction out of bounds: start=%d, base length=%d", start, len(base))
			}
			if end > len(base) {
				end = len(base)
			}
			result.Write(base[start:end])
		case ADD:
			result.Write(instruction.Data)
		default:
			return nil, fmt.Errorf("unknown delta instruction type: %s", instruction.Type)
		}
	}

	return result.Bytes(), nil
}
