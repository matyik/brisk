package stdlib

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/dop251/goja"
)

// InitFetch builds the global fetch function
func InitFetch(vm *goja.Runtime, config RuntimeConfig) goja.Value {
	
	fetchFn := func(call goja.FunctionCall) goja.Value {
		promise, resolve, reject := vm.NewPromise()

		if len(call.Arguments) == 0 {
			reject(vm.ToValue("TypeError: fetch requires at least 1 argument"))
			return vm.ToValue(promise)
		}
		
		url := call.Argument(0).String()
		method := "GET"
		var bodyReader io.Reader = nil
		var headers *goja.Object

		if len(call.Arguments) > 1 && !goja.IsUndefined(call.Argument(1)) {
			options := call.Argument(1).ToObject(vm)

			if m := options.Get("method"); m != nil && !goja.IsUndefined(m) {
				method = strings.ToUpper(m.String())
			}
			if b := options.Get("body"); b != nil && !goja.IsUndefined(b) {
				bodyReader = strings.NewReader(b.String())
			}
			if h := options.Get("headers"); h != nil && !goja.IsUndefined(h) {
				headers = h.ToObject(vm)
			}
		}

		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			reject(vm.ToValue("NetworkError: " + err.Error()))
			return vm.ToValue(promise)
		}

		if headers != nil {
			for _, key := range headers.Keys() {
				req.Header.Set(key, headers.Get(key).String())
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			reject(vm.ToValue("NetworkError: " + err.Error()))
			return vm.ToValue(promise)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			reject(vm.ToValue("IOError: " + err.Error()))
			return vm.ToValue(promise)
		}

		responseObj := vm.NewObject()
		responseObj.Set("ok", resp.StatusCode >= 200 && resp.StatusCode < 300)
		responseObj.Set("status", resp.StatusCode)

		responseObj.Set("text", func(goja.FunctionCall) goja.Value {
			textPromise, textResolve, _ := vm.NewPromise()
			textResolve(string(bodyBytes))
			return vm.ToValue(textPromise)
		})

		responseObj.Set("json", func(goja.FunctionCall) goja.Value {
			jsonPromise, jsonResolve, jsonReject := vm.NewPromise()
			var jsonData interface{}
			if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
				jsonReject(vm.ToValue("SyntaxError: Failed to parse JSON"))
			} else {
				jsonResolve(vm.ToValue(jsonData))
			}
			return vm.ToValue(jsonPromise)
		})

		resolve(responseObj)
		return vm.ToValue(promise)
	}

	return vm.ToValue(fetchFn)
}