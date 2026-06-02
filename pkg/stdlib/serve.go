package stdlib

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/dop251/goja"
)

// InitBrisk lazily builds the global Brisk namespace (for APIs like Brisk.serve)
func InitBrisk(vm *goja.Runtime, config RuntimeConfig) goja.Value {
	briskObj := vm.NewObject()

	// Mutex guarantees that only ONE HTTP request executes JS at a time
	// prevents memory corruption in single-threaded js vm
	var vmMutex sync.Mutex

	_ = briskObj.Set("serve", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.NewTypeError("Brisk.serve requires 2 arguments: port and a callback function"))
		}

		port := call.Argument(0).ToInteger()
		callback, ok := goja.AssertFunction(call.Argument(1))
		if !ok {
			panic(vm.NewTypeError("Second argument to Brisk.serve must be a function"))
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, _ := io.ReadAll(r.Body)
			defer r.Body.Close()

			vmMutex.Lock()
			defer vmMutex.Unlock()

			jsReq := vm.NewObject()
			_ = jsReq.Set("method", r.Method)
			_ = jsReq.Set("url", r.URL.Path)
			_ = jsReq.Set("body", string(bodyBytes))

			res, err := callback(goja.Undefined(), jsReq)
			if err != nil {
				http.Error(w, fmt.Sprintf("Runtime Error: %v", err), http.StatusInternalServerError)
				return
			}

			if goja.IsNull(res) || goja.IsUndefined(res) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}

			if resObj, ok := res.(*goja.Object); ok {
				statusVal := resObj.Get("status")
				if statusVal != nil && !goja.IsUndefined(statusVal) {
					w.WriteHeader(int(statusVal.ToInteger()))
				} else {
					w.WriteHeader(http.StatusOK) // Default 200
				}

				bodyVal := resObj.Get("body")
				if bodyVal != nil && !goja.IsUndefined(bodyVal) {
					w.Write([]byte(bodyVal.String()))
				}
				return
			}

			// Fallback: Just stringify whatever they returned
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(res.String()))
		})

		err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
		if err != nil {
			panic(vm.NewGoError(err))
		}

		return goja.Undefined()
	})

	return briskObj
}