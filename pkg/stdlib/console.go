package stdlib

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

func RegisterConsole(vm *goja.Runtime) error {
	console := vm.NewObject()
	
	err := console.Set("log", func(call goja.FunctionCall) goja.Value {
		var strs []string
		for _, arg := range call.Arguments {
			strs = append(strs, arg.String())
		}
		fmt.Println(strings.Join(strs, " "))
		return goja.Undefined()
	})
	
	if err != nil {
		return err
	}
	
	return vm.Set("console", console)
}