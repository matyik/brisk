package stdlib

import (
	"os"
	"strings"

	"github.com/dop251/goja"
)

func InitProcess(vm *goja.Runtime, config RuntimeConfig) goja.Value {
	process := vm.NewObject()
	
	env := vm.NewObject()
	for _, e := range os.Environ() {
		// SplitN in case value contains '='
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			env.Set(pair[0], pair[1])
		}
	}
	process.Set("env", env)

	process.Set("argv", vm.ToValue(os.Args))

	process.Set("cwd", func(goja.FunctionCall) goja.Value {
		dir, err := os.Getwd()
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(dir)
	})

	process.Set("exit", func(call goja.FunctionCall) goja.Value {
		code := 0
		if len(call.Arguments) > 0 {
			code = int(call.Argument(0).ToInteger())
		}
		os.Exit(code)
		return goja.Undefined()
	})

	return process
}