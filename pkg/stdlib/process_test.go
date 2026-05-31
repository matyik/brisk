package stdlib

import (
	"os"
	"testing"

	"github.com/dop251/goja"
)

func TestProcess_StandardAPIs(t *testing.T) {
	testKey := "BRISK_TEST_VAR"
	testVal := "super_secret"
	os.Setenv(testKey, testVal)
	defer os.Unsetenv(testKey)

	vm := goja.New()
	err := RegisterProcess(vm)
	if err != nil {
		t.Fatalf("Failed to register process: %v", err)
	}

	jsCode := `
		const testVar = process.env.BRISK_TEST_VAR;
		const isArgvArray = Array.isArray(process.argv);
		const argCount = process.argv.length;
		const currentDir = process.cwd();
		
		({ testVar, isArgvArray, argCount, currentDir });
	`

	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	data := result.Export().(map[string]interface{})

	if data["testVar"] != testVal {
		t.Errorf("Expected env var to be '%s', got '%v'", testVal, data["testVar"])
	}

	if data["isArgvArray"] != true {
		t.Errorf("Expected process.argv to be a JavaScript Array")
	}

	if int(data["argCount"].(int64)) == 0 {
		t.Errorf("Expected process.argv to have length > 0")
	}

	expectedDir, _ := os.Getwd()
	if data["currentDir"] != expectedDir {
		t.Errorf("Expected cwd to be '%s', got '%v'", expectedDir, data["currentDir"])
	}
}