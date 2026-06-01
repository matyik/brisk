package stdlib

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dop251/goja"
)

const (
	ColorReset  = "\033[0m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
)

// InitConsole builds the global console object
func InitConsole(vm *goja.Runtime, config RuntimeConfig) goja.Value {
	console := vm.NewObject()

	createPrinter := func(out io.Writer, prefix string) func(goja.FunctionCall) goja.Value {
		return func(call goja.FunctionCall) goja.Value {
			var strs []string
			
			if prefix != "" {
				strs = append(strs, prefix)
			}

			for _, arg := range call.Arguments {
				strs = append(strs, arg.String())
			}

			finalString := strings.Join(strs, " ")
			if prefix != "" {
				finalString += ColorReset
			}
			
			fmt.Fprintln(out, finalString)
			return goja.Undefined()
		}
	}

	_ = console.Set("log", createPrinter(os.Stdout, ""))
	_ = console.Set("info", createPrinter(os.Stdout, ""))

	_ = console.Set("warn", createPrinter(os.Stderr, ColorYellow+"[WARN]"))

	_ = console.Set("error", createPrinter(os.Stderr, ColorRed+"[ERROR]"))

	return console
}