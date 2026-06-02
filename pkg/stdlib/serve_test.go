package stdlib

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/dop251/goja"
)

func TestBrisk_Serve(t *testing.T) {
	vm := goja.New()
	LazyInject(vm, RuntimeConfig{})

	jsCode := `
		// Boot the server on a unique test port
		Brisk.serve(9999, (req) => {
			if (req.method === "GET" && req.url === "/test") {
				return {
					status: 201,
					body: "Hello from the test suite!"
				};
			}
			return { status: 404, body: "Not Found" };
		});
	`

	go func() {
		_, err := vm.RunString(jsCode)
		if err != nil {
			t.Logf("Server crashed: %v", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://localhost:9999/test")
	if err != nil {
		t.Fatalf("Failed to make HTTP request to Brisk server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	if string(bodyBytes) != "Hello from the test suite!" {
		t.Errorf("Expected 'Hello from the test suite!', got '%s'", string(bodyBytes))
	}
}