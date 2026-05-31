package transpiler

import (
	"strings"
	"testing"
)

func TestProcess_JavaScript(t *testing.T) {
	code := "const x = 5;"
	result, err := Process("test.js", code)
	
	if err != nil {
		t.Fatalf("Expected no error for JS file, got: %v", err)
	}
	if result != code {
		t.Errorf("Expected JS code to remain unchanged. Got: %s", result)
	}
}

func TestProcess_TypeScript(t *testing.T) {
	code := "const greeting: string = 'hello';"
	result, err := Process("test.ts", code)
	
	if err != nil {
		t.Fatalf("Expected no error for valid TS file, got: %v", err)
	}
	
	// If the transpiler worked, the word 'string' (the type definition) should be gone
	if strings.Contains(result, "string") {
		t.Errorf("Expected TypeScript types to be stripped. Result still contained types: %s", result)
	}
}

func TestProcess_TypeScript_SyntaxError(t *testing.T) {
	badCode := "const a: string = ;" // Invalid syntax
	_, err := Process("test.ts", badCode)
	
	if err == nil {
		t.Error("Expected an error for malformed TypeScript, but got none")
	}
}