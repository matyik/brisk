package stdlib

import "github.com/dop251/goja"

// RuntimeConfig holds dynamic variables (like sandbox paths) for modules
type RuntimeConfig struct {
	SandboxDir string
}

// ModuleFactory creates the heavy JS object/function ONLY when requested
type ModuleFactory func(vm *goja.Runtime, config RuntimeConfig) goja.Value

// Registry holds the map of global variable names to their lazy factories
func Registry() map[string]ModuleFactory {
	return map[string]ModuleFactory{
		"console": InitConsole,
		"fetch":   InitFetch,
		"URL":     InitURL,
		"process": InitProcess,
		"fs":      InitFS,
		"crypto":  InitCrypto,
		"Brisk":   InitBrisk,
	}
}

// LazyInject mounts the registry into the VM using memoized getters
func LazyInject(vm *goja.Runtime, config RuntimeConfig) {
	for name, factory := range Registry() {
		var cached goja.Value 

		getter := func(call goja.FunctionCall) goja.Value {
			if cached == nil {
				cached = factory(vm, config)
			}
			return cached
		}

		vm.GlobalObject().DefineAccessorProperty(
			name,
			vm.ToValue(getter),
			nil,
			goja.FLAG_TRUE, // Configurable
			goja.FLAG_TRUE, // Enumerable
		)
	}
}