package objects

import (
	"bytes"
	"testing"
)

func TestComputeAndApplyDelta(t *testing.T) {
	base := []byte("The quick brown fox jumps over the lazy dog.")
	updated := []byte("The quick red fox jumps over the smart dog.")

	delta, err := ComputeDelta(base, updated)
	if err != nil {
		t.Fatalf("failed to compute delta: %v", err)
	}

	result, err := ApplyDelta(base, delta)
	if err != nil {
		t.Fatalf("failed to apply delta: %v", err)
	}

	if !bytes.Equal(result, updated) {
		t.Errorf("result does not match updated: got %q, want %q", result, updated)
	}
}

func TestDeltaWithIdenticalData(t *testing.T) {
	base := []byte("Identical data should result in no ADD instructions.")
	updated := []byte("Identical data should result in no ADD instructions.")

	delta, err := ComputeDelta(base, updated)
	if err != nil {
		t.Fatalf("failed to compute delta: %v", err)
	}

	result, err := ApplyDelta(base, delta)
	if err != nil {
		t.Fatalf("failed to apply delta: %v", err)
	}

	if !bytes.Equal(result, updated) {
		t.Errorf("result does not match updated: got %q, want %q", result, updated)
	}

	// Ensure there are no ADD instructions
	for _, instruction := range delta {
		if instruction.Type == ADD {
			t.Errorf("unexpected ADD instruction: %v", instruction)
		}
	}
}

func TestDeltaWithEmptyBase(t *testing.T) {
	base := []byte("")
	updated := []byte("New data entirely.")

	delta, err := ComputeDelta(base, updated)
	if err != nil {
		t.Fatalf("failed to compute delta: %v", err)
	}

	result, err := ApplyDelta(base, delta)
	if err != nil {
		t.Fatalf("failed to apply delta: %v", err)
	}

	if !bytes.Equal(result, updated) {
		t.Errorf("result does not match updated: got %q, want %q", result, updated)
	}

	// Ensure all instructions are ADD
	for _, instruction := range delta {
		if instruction.Type != ADD {
			t.Errorf("unexpected instruction type: got %v, want ADD", instruction.Type)
		}
	}
}

func TestDeltaWithEmptyUpdated(t *testing.T) {
	base := []byte("Data to be removed entirely.")
	updated := []byte("")

	delta, err := ComputeDelta(base, updated)
	if err != nil {
		t.Fatalf("failed to compute delta: %v", err)
	}

	result, err := ApplyDelta(base, delta)
	if err != nil {
		t.Fatalf("failed to apply delta: %v", err)
	}

	if !bytes.Equal(result, updated) {
		t.Errorf("result does not match updated: got %q, want %q", result, updated)
	}

	// Ensure there are no COPY instructions
	for _, instruction := range delta {
		if instruction.Type != ADD {
			t.Errorf("unexpected instruction type: got %v, want ADD", instruction.Type)
		}
	}
}

func TestDeltaWithLargeData(t *testing.T) {
	base := bytes.Repeat([]byte("A"), 10*BlockSize)                                                      // 10 blocks of 'A'
	updated := append(bytes.Repeat([]byte("A"), 5*BlockSize), bytes.Repeat([]byte("B"), 5*BlockSize)...) // Half 'A', half 'B'

	delta, err := ComputeDelta(base, updated)
	if err != nil {
		t.Fatalf("failed to compute delta: %v", err)
	}

	result, err := ApplyDelta(base, delta)
	if err != nil {
		t.Fatalf("failed to apply delta: %v", err)
	}

	if !bytes.Equal(result, updated) {
		t.Errorf("result does not match updated: data mismatch")
	}
}
