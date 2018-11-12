package memorybackend

import (
	"testing"
)

func TestCreatingNewMemoryBackend(t *testing.T) {
	backend, err := NewBackend()
	if err != nil || backend == nil {
		t.Fatalf("Could not get a new Memory Backend: %s", err)
	}
}
