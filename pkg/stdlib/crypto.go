package stdlib

import (
	"crypto/rand"
	"fmt"

	"github.com/dop251/goja"
)

// InitCrypto lazily builds the global Web Standard crypto object
func InitCrypto(vm *goja.Runtime, config RuntimeConfig) goja.Value {
	crypto := vm.NewObject()

	_ = crypto.Set("randomUUID", func(call goja.FunctionCall) goja.Value {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			panic(vm.NewGoError(fmt.Errorf("failed to generate random bytes: %w", err)))
		}

		b[6] = (b[6] & 0x0f) | 0x40 // Version 4
		b[8] = (b[8] & 0x3f) | 0x80 // Variant 10

		uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

		return vm.ToValue(uuid)
	})

	_ = crypto.Set("getRandomValues", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(vm.NewTypeError("getRandomValues requires 1 argument"))
		}

		obj := call.Argument(0).ToObject(vm)
		lengthVal := obj.Get("length")

		if lengthVal == nil || goja.IsUndefined(lengthVal) {
			panic(vm.NewTypeError("Argument must be an ArrayBufferView"))
		}

		length := lengthVal.ToInteger()

		if length > 65536 {
			panic(vm.NewGoError(fmt.Errorf("QuotaExceededError: requested length exceeds 65,536 limit")))
		}

		for i := range length {
			b := make([]byte, 4)
			rand.Read(b)
			
			val := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
			
			obj.Set(fmt.Sprintf("%d", i), val)
		}

		return obj 
	})

	return crypto
}