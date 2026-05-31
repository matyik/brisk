package stdlib

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dop251/goja"
)

func TestRegisterConsole(t *testing.T) {
	vm := goja.New()
	err := RegisterConsole(vm)
	if err != nil {
		t.Fatalf("Failed to register console standard library: %v", err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	_, err = vm.RunString(`console.log("hello", "brisk", 123)`)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("JavaScript execution failed: %v", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	expected := "hello brisk 123"
	if output != expected {
		t.Errorf("Expected output '%s', got '%s'", expected, output)
	}
}