package engine

import (
	"github.com/dop251/goja"
	"github.com/matyik/brisk/pkg/stdlib"
)

type VM struct {
	runtime *goja.Runtime
}

// New initializes a fresh VM and binds the standard libraries
func New(baseDir string) (*VM, error) {
	vm := goja.New()

	config := stdlib.RuntimeConfig{
		SandboxDir: baseDir,
	}

	stdlib.LazyInject(vm, config)

	return &VM{runtime: vm}, nil
}

// Run executes a string of JavaScript code
func (v *VM) Run(code string) error {
	_, err := v.runtime.RunString(code)
	return err
}