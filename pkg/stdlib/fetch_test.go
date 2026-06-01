package stdlib

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dop251/goja"
)

func TestFetch_AdvancedOptions(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.Header.Get("User-Agent") != "brisk-test-suite" {
			t.Errorf("Expected User-Agent 'brisk-test-suite', got %s", r.Header.Get("User-Agent"))
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		expectedBody := `{"ping":"pong"}`
		if string(bodyBytes) != expectedBody {
			t.Errorf("Expected body '%s', got '%s'", expectedBody, string(bodyBytes))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) 
		w.Write([]byte(`{"success": true, "id": 99}`))
	}))
	defer mockServer.Close() 

	vm := goja.New()
	config := RuntimeConfig{}
	LazyInject(vm, config)

	vm.Set("MOCK_URL", mockServer.URL)

	jsCode := `
		fetch(MOCK_URL, {
			method: 'POST',
			headers: {
				'User-Agent': 'brisk-test-suite',
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ ping: 'pong' })
		});
	`
	
	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	promise, ok := result.Export().(*goja.Promise)
	if !ok {
		t.Fatalf("Expected fetch to return a Promise")
	}

	if promise.State() != goja.PromiseStateFulfilled {
		t.Fatalf("Expected Promise to be fulfilled, got state: %v", promise.State())
	}

	responseObj := promise.Result().ToObject(vm)

	status := responseObj.Get("status").ToInteger()
	if status != 201 {
		t.Errorf("Expected status 201, got %d", status)
	}

	isOk := responseObj.Get("ok").ToBoolean()
	if !isOk {
		t.Errorf("Expected ok to be true")
	}

	jsonFunc, ok := goja.AssertFunction(responseObj.Get("json"))
	if !ok {
		t.Fatalf("Expected res.json to be a function")
	}

	jsonPromiseVal, err := jsonFunc(goja.Undefined())
	if err != nil {
		t.Fatalf("res.json() threw an error: %v", err)
	}

	jsonPromise := jsonPromiseVal.Export().(*goja.Promise)
	jsonData := jsonPromise.Result().Export().(map[string]interface{})

	if jsonData["success"] != true {
		t.Errorf("Expected success=true in JSON response")
	}
	if int(jsonData["id"].(float64)) != 99 {
		t.Errorf("Expected id=99 in JSON response")
	}
}