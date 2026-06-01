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
	}
}

// LazyInject mounts the registry into the VM using memoized getters
func LazyInject(vm *goja.Runtime, config RuntimeConfig) {
	for name, factory := range Registry() {
		// Go 1.22 loop variable capture is safe here!
		var cached goja.Value 

		getter := func(call goja.FunctionCall) goja.Value {
			if cached == nil {
				// Boot the heavy module exactly ONCE, only when accessed
				cached = factory(vm, config)
			}
			return cached
		}

		// Inject the lazy getter into the global JavaScript namespace
		vm.GlobalObject().DefineAccessorProperty(
			name,
			vm.ToValue(getter),
			nil,
			goja.FLAG_TRUE, // Configurable
			goja.FLAG_TRUE, // Enumerable
		)
	}
}