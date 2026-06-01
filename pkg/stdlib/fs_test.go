package stdlib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
)

func TestFS_SandboxedReadWrite(t *testing.T) {
	sandboxDir := t.TempDir()

	vm := goja.New()
	config := RuntimeConfig{SandboxDir: sandboxDir}
	LazyInject(vm, config)

	jsCode := `
		fs.writeFileSync('test.txt', 'hello from brisk');
		const content = fs.readFileSync('test.txt');
		content;
	`
	
	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	if result.String() != "hello from brisk" {
		t.Errorf("Expected 'hello from brisk', got '%s'", result.String())
	}

	hostBytes, err := os.ReadFile(filepath.Join(sandboxDir, "test.txt"))
	if err != nil || string(hostBytes) != "hello from brisk" {
		t.Errorf("Host OS verification failed, file was not written correctly")
	}
}

func TestFS_PathTraversalBlock(t *testing.T) {
	sandboxDir := t.TempDir()
	
	parentDir := filepath.Dir(sandboxDir)
	secretFile := filepath.Join(parentDir, "secret.txt")
	os.WriteFile(secretFile, []byte("super_secret_password"), 0644)
	defer os.Remove(secretFile)

	vm := goja.New()
	config := RuntimeConfig{SandboxDir: sandboxDir}
	LazyInject(vm, config)

	jsCode := `
		let caught = false;
		try {
			fs.readFileSync('../secret.txt');
		} catch (e) {
			caught = true; 
		}
		caught;
	`

	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	if result.Export().(bool) != true {
		t.Error("Security failure: Path traversal was not blocked!")
	}
}