package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: brisk <file.js|file.ts>")
		os.Exit(1)
	}
	filePath := os.Args[1]

	codeBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	code := string(codeBytes)

	if strings.HasSuffix(filePath, ".ts") {
		result := api.Transform(code, api.TransformOptions{
			Loader: api.LoaderTS,
		})

		if len(result.Errors) > 0 {
			fmt.Println("TypeScript compilation failed:")
			for _, err := range result.Errors {
				fmt.Printf("- %s\n", err.Text)
			}
			os.Exit(1)
		}
		code = string(result.Code)
	}

	vm := goja.New()

	console := vm.NewObject()
	err = console.Set("log", func(call goja.FunctionCall) goja.Value {
		var strs []string
		for _, arg := range call.Arguments {
			strs = append(strs, arg.String())
		}
		fmt.Println(strings.Join(strs, " "))
		return goja.Undefined()
	})
	if err != nil {
		panic("Failed to initialize console.log")
	}
	vm.Set("console", console)

	_, err = vm.RunString(code)
	if err != nil {
		fmt.Printf("Brisk Runtime Error:\n%v\n", err)
		os.Exit(1)
	}
}