package stdlib

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

// InitConsole builds the console object
func InitConsole(vm *goja.Runtime, config RuntimeConfig) goja.Value {
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
		panic(vm.NewGoError(err))
	}
	
	return console
}