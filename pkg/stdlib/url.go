package stdlib

import (
	"net/url"

	"github.com/dop251/goja"
)

func RegisterURL(vm *goja.Runtime) error {
	err := vm.Set("URL", func(call goja.ConstructorCall) *goja.Object {
		if len(call.Arguments) == 0 {
			panic(vm.NewTypeError("Failed to construct 'URL': 1 argument required, but only 0 present."))
		}

		rawURL := call.Argument(0).String()
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			panic(vm.NewTypeError("Invalid URL: " + err.Error()))
		}

		instance := call.This

		// Initialize static properties
		instance.Set("href", parsedURL.String())
		instance.Set("protocol", parsedURL.Scheme+":")
		instance.Set("hostname", parsedURL.Hostname())
		instance.Set("pathname", parsedURL.Path)
		instance.Set("search", "?"+parsedURL.RawQuery)

		queryParams := parsedURL.Query()
		searchParams := vm.NewObject()

		syncState := func() {
			parsedURL.RawQuery = queryParams.Encode()
			
			searchString := ""
			if parsedURL.RawQuery != "" {
				searchString = "?" + parsedURL.RawQuery
			}
			
			instance.Set("search", searchString)
			instance.Set("href", parsedURL.String())
		}

		// searchParams.get()
		searchParams.Set("get", func(c goja.FunctionCall) goja.Value {
			if len(c.Arguments) == 0 {
				return goja.Undefined()
			}
			key := c.Argument(0).String()
			val := queryParams.Get(key)
			if val == "" && !queryParams.Has(key) {
				return goja.Null() // Spec says return null if not found
			}
			return vm.ToValue(val)
		})

		// searchParams.set()
		searchParams.Set("set", func(c goja.FunctionCall) goja.Value {
			if len(c.Arguments) < 2 {
				return goja.Undefined()
			}
			key := c.Argument(0).String()
			val := c.Argument(1).String()
			queryParams.Set(key, val)
			syncState()
			return goja.Undefined()
		})

		// searchParams.append()
		searchParams.Set("append", func(c goja.FunctionCall) goja.Value {
			key := c.Argument(0).String()
			val := c.Argument(1).String()
			queryParams.Add(key, val)
			syncState()
			return goja.Undefined()
		})

		// searchParams.delete()
		searchParams.Set("delete", func(c goja.FunctionCall) goja.Value {
			key := c.Argument(0).String()
			queryParams.Del(key)
			syncState()
			return goja.Undefined()
		})

		instance.Set("searchParams", searchParams)

		// URL.toString() method
		instance.Set("toString", func(goja.FunctionCall) goja.Value {
			return vm.ToValue(parsedURL.String())
		})

		return instance
	})

	return err
}