package engine

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/matyik/brisk/pkg/stdlib"
)

type VM struct {
	runtime *goja.Runtime
}

// New initializes a fresh VM and binds the standard libraries
func New() (*VM, error) {
	vm := goja.New()

	if err := stdlib.RegisterConsole(vm); err != nil {
		return nil, fmt.Errorf("failed to register console: %w", err)
	}

	if err := stdlib.RegisterFetch(vm); err != nil {
		return nil, fmt.Errorf("failed to register fetch: %w", err)
	}

	

	return &VM{runtime: vm}, nil
}

// Run executes a string of JavaScript code
func (v *VM) Run(code string) error {
	_, err := v.runtime.RunString(code)
	return err
}