package stdlib

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dop251/goja"
)

func TestConsole_StandardStreamsAndColors(t *testing.T) {
	vm := goja.New()
	config := RuntimeConfig{}
	LazyInject(vm, config)

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	
	os.Stdout = wOut
	os.Stderr = wErr

	jsCode := `
		console.log("standard log");
		console.info("standard info");
		console.warn("watch out");
		console.error("fatal crash");
	`
	_, err := vm.RunString(jsCode)

	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("JavaScript execution failed: %v", err)
	}

	var bufOut, bufErr bytes.Buffer
	_, _ = io.Copy(&bufOut, rOut)
	_, _ = io.Copy(&bufErr, rErr)

	stdoutStr := strings.TrimSpace(bufOut.String())
	stderrStr := strings.TrimSpace(bufErr.String())

	expectedStdout := "standard log\nstandard info"
	if stdoutStr != expectedStdout {
		t.Errorf("Expected stdout:\n%q\nGot:\n%q", expectedStdout, stdoutStr)
	}

	expectedStderr := ColorYellow + "[WARN] watch out" + ColorReset + "\n" +
		              ColorRed + "[ERROR] fatal crash" + ColorReset
	
	if stderrStr != expectedStderr {
		t.Errorf("Expected stderr:\n%q\nGot:\n%q", expectedStderr, stderrStr)
	}
}