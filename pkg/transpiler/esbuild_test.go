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

func TestNodePolyfillPlugin(t *testing.T) {
	tempDir := t.TempDir()
	entryPoint := filepath.Join(tempDir, "legacy.ts")

	code := `
		import fs from 'node:fs';
		import { randomUUID } from 'crypto';
		import process from 'process';

		const id = randomUUID();
		const fileData = fs.readFileSync('test.txt');
		console.log(process.cwd());
	`
	
	err := os.WriteFile(entryPoint, []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	output, err := Process(entryPoint)
	if err != nil {
		t.Fatalf("Transpiler crashed! Error: %v", err)
	}

	
	if !strings.Contains(output, "globalThis.fs.readFileSync") {
		t.Errorf("Expected output to contain polyfilled 'globalThis.fs.readFileSync', got:\n%s", output)
	}

	if !strings.Contains(output, "globalThis.crypto.randomUUID") {
		t.Errorf("Expected output to contain polyfilled 'globalThis.crypto.randomUUID', got:\n%s", output)
	}

	if !strings.Contains(output, "globalThis.process") {
		t.Errorf("Expected output to contain polyfilled 'globalThis.process', got:\n%s", output)
	}

	if strings.Contains(output, "require('node:fs')") || strings.Contains(output, "require(\"node:fs\")") {
		t.Errorf("Security Failure: The plugin failed to intercept the Node import. A native 'require' leaked into the bundle!")
	}
}