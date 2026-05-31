package transpiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcess_BundlesImportsAndStripsTypes(t *testing.T) {
	tempDir := t.TempDir()

	mathCode := `
		export function add(a: number, b: number): number {
			return a + b;
		}
	`
	err := os.WriteFile(filepath.Join(tempDir, "math.ts"), []byte(mathCode), 0644)
	if err != nil {
		t.Fatalf("Failed to write math.ts: %v", err)
	}

	indexCode := `
		import { add } from './math';
		const result: number = add(5, 10);
		console.log(result);
	`
	indexPath := filepath.Join(tempDir, "index.ts")
	err = os.WriteFile(indexPath, []byte(indexCode), 0644)
	if err != nil {
		t.Fatalf("Failed to write index.ts: %v", err)
	}

	result, err := Process(indexPath)

	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if strings.Contains(result, "number") {
		t.Errorf("Expected types to be stripped, but 'number' was found in output:\n%s", result)
	}

	if !strings.Contains(result, "a + b") {
		t.Errorf("Expected imported math module to be bundled, but logic was missing:\n%s", result)
	}
}

func TestProcess_MissingImport(t *testing.T) {
	tempDir := t.TempDir()

	badCode := `
		import { nothing } from './does-not-exist';
		console.log(nothing);
	`
	
	badPath := filepath.Join(tempDir, "bad.ts")
	_ = os.WriteFile(badPath, []byte(badCode), 0644)

	_, err := Process(badPath)

	if err == nil {
		t.Error("Expected an error for a missing import, but got none")
	}
}