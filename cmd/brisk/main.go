package main

import (
	"fmt"
	"os"

	"github.com/matyik/brisk/pkg/engine"
	"github.com/matyik/brisk/pkg/transpiler"
)

var Version = "dev"

func main() {
	if len(os.Args) >= 2 {
		arg := os.Args[1]
		if arg == "--version" || arg == "-v" {
			fmt.Printf("Brisk v%s\n", Version)
			os.Exit(0)
		}
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: brisk <file.js|file.ts>")
		os.Exit(1)
	}
	filePath := os.Args[1]

	finalCode, err := transpiler.Process(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to determine working directory: %v\n", err)
		os.Exit(1)
	}

	vm, err := engine.New(cwd)
	if err != nil {
		fmt.Printf("Engine initialization failed: %v\n", err)
		os.Exit(1)
	}

	if err := vm.Run(finalCode); err != nil {
		fmt.Printf("Brisk Runtime Error:\n%v\n", err)
		os.Exit(1)
	}
}