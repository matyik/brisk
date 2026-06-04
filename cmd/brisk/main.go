package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/matyik/brisk/pkg/engine"
	"github.com/matyik/brisk/pkg/transpiler"
)

var Version = "dev"

func getBinarySeparator() []byte {
	return []byte(strings.Join([]string{"\n---BRISK", "_EMBED_V1---\n"}, ""))
}

func main() {
	binarySeparator := getBinarySeparator()

	exePath, err := os.Executable()
	if err == nil {
		exeBytes, err := os.ReadFile(exePath)
		if err == nil {
			parts := bytes.Split(exeBytes, binarySeparator)
			if len(parts) > 1 {
				embeddedCode := string(parts[len(parts)-1])
				
				cwd, _ := os.Getwd()
				vm, _ := engine.New(cwd)
				if err := vm.Run(embeddedCode); err != nil {
					fmt.Printf("Brisk Runtime Error:\n%v\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			}
		}
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "--version" || command == "-v" {
		fmt.Printf("Brisk Engine v%s\n", Version)
		os.Exit(0)
	}

	if command == "compile" {
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing file path.\nUsage: brisk compile <file.ts>")
			os.Exit(1)
		}
		compileToBinary(os.Args[2], binarySeparator)
		os.Exit(0)
	}

	runFile(command)
}

func compileToBinary(filePath string, separator []byte) {
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	outName := baseName
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}

	distDir := "./dist"
	outPath := filepath.Join(distDir, outName)

	fmt.Printf("Compiling %s into a standalone executable at %s...\n", filePath, outPath)

	if err := os.MkdirAll(distDir, 0755); err != nil {
		fmt.Printf("Failed to create dist directory: %v\n", err)
		os.Exit(1)
	}

	finalJS, err := transpiler.Process(filePath)
	if err != nil {
		fmt.Printf("Compile Error: %v\n", err)
		os.Exit(1)
	}

	briskExe, _ := os.Executable()
	briskBytes, err := os.ReadFile(briskExe)
	if err != nil {
		fmt.Printf("Failed to read base executable: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	outFile.Write(briskBytes)
	outFile.Write(separator)
	outFile.Write([]byte(finalJS))

	outFile.Chmod(0755)

	fmt.Printf("Success: Standalone binary created: %s\n", outPath)
}

func runFile(filePath string) {
	finalCode, err := transpiler.Process(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cwd, _ := os.Getwd()
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

func printUsage() {
	fmt.Println("Brisk Edge Runtime")
	fmt.Println("\nUsage:")
	fmt.Println("  brisk <file.ts>         Run a JavaScript or TypeScript file")
	fmt.Println("  brisk compile <file.ts> Compile a file into a standalone executable")
	fmt.Println("  brisk --version         Print the current version")
}