package engine

import (
	"testing"
)

func TestEngine_RunValidJS(t *testing.T) {
	vm, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	code := `
		const a = 10;
		const b = 20;
		const sum = a + b;
	`
	err = vm.Run(code)
	if err != nil {
		t.Errorf("Expected no error when running valid JavaScript, got: %v", err)
	}
}

func TestEngine_RunInvalidJS(t *testing.T) {
	vm, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	// Expected to throw ReferenceError
	code := `
		console.log(undefinedVariable);
	`
	err = vm.Run(code)
	if err == nil {
		t.Error("Expected an error when running invalid JavaScript, but got nil")
	}
}