package stdlib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

func InitFS(vm *goja.Runtime, config RuntimeConfig) goja.Value {
	absRoot, err := filepath.Abs(config.SandboxDir)
	if err != nil {
		panic(vm.NewGoError(fmt.Errorf("failed to resolve fs root: %w", err)))
	}

	fs := vm.NewObject()

	checkPath := func(targetPath string) (string, error) {
		joined := filepath.Join(absRoot, targetPath)

		rel, err := filepath.Rel(absRoot, joined)
		
		if err != nil || strings.HasPrefix(rel, "..") {
			return "", fmt.Errorf("EACCES: permission denied, path '%s' is outside the sandbox", targetPath)
		}
		
		return joined, nil
	}

	fs.Set("readFileSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(vm.NewTypeError("readFileSync requires a file path"))
		}
		
		targetPath := call.Argument(0).String()
		safePath, err := checkPath(targetPath)
		if err != nil {
			panic(vm.NewGoError(err))
		}

		data, err := os.ReadFile(safePath)
		if err != nil {
			panic(vm.NewGoError(err))
		}

		return vm.ToValue(string(data))
	})

	fs.Set("writeFileSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.NewTypeError("writeFileSync requires a file path and data"))
		}

		targetPath := call.Argument(0).String()
		data := call.Argument(1).String()

		safePath, err := checkPath(targetPath)
		if err != nil {
			panic(vm.NewGoError(err))
		}

		err = os.WriteFile(safePath, []byte(data), 0644)
		if err != nil {
			panic(vm.NewGoError(err))
		}

		return goja.Undefined()
	})

	return fs
}