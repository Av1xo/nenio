package objects

import (
	"bytes"
	"fmt"

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
	Data   []byte
}

func ComputeDelta(base, updated []byte) ([]DeltaInstruction, error) {
	if base == nil || updated == nil {
		return nil, fmt.Errorf("base or updated data is nil")
	}

	var delta []DeltaInstruction
	baseHashes := make(map[[32]byte]int)

	for i := 0; i < len(base); i += BlockSize {
		end := i + BlockSize
		if end > len(base) {
			end = len(base)
		}
		hash := blake3.Sum256(base[i:end])
		baseHashes[hash] = i
	}

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

func ApplyDelta(base []byte, delta []DeltaInstruction) ([]byte, error) {
	if base == nil {
		return nil, fmt.Errorf("base data is nil")
	}

	var result bytes.Buffer
	for _, instruction := range delta {
		switch instruction.Type {
		case COPY:
			start := instruction.Offset
			end := start + BlockSize
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
