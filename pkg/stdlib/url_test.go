package stdlib

import (
	"testing"

	"github.com/dop251/goja"
)

func TestURL_ConstructorAndProperties(t *testing.T) {
	vm := goja.New()
	err := RegisterURL(vm)
	if err != nil {
		t.Fatalf("Failed to register URL: %v", err)
	}

	jsCode := `
		const u = new URL('https://api.github.com/users/brisk?sort=desc');
		
		// Export an object back to Go so we can assert all values cleanly
		({
			protocol: u.protocol,
			hostname: u.hostname,
			pathname: u.pathname,
			search: u.search,
			href: u.href
		});
	`

	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	data := result.Export().(map[string]interface{})

	if data["protocol"] != "https:" {
		t.Errorf("Expected protocol 'https:', got %v", data["protocol"])
	}
	if data["hostname"] != "api.github.com" {
		t.Errorf("Expected hostname 'api.github.com', got %v", data["hostname"])
	}
	if data["pathname"] != "/users/brisk" {
		t.Errorf("Expected pathname '/users/brisk', got %v", data["pathname"])
	}
	if data["search"] != "?sort=desc" {
		t.Errorf("Expected search '?sort=desc', got %v", data["search"])
	}
}

func TestURL_SearchParamsMutation(t *testing.T) {
	vm := goja.New()
	RegisterURL(vm)

	jsCode := `
		const u = new URL('https://example.com/api?limit=10');
		
		// 1. Read existing param
		const initialLimit = u.searchParams.get('limit');
		
		// 2. Set new param (should overwrite)
		u.searchParams.set('limit', '50');
		
		// 3. Append param (adds new key)
		u.searchParams.append('sort', 'asc');
		
		// 4. Delete a param
		u.searchParams.set('temp', 'delete-me');
		u.searchParams.delete('temp');
		
		({
			initialLimit: initialLimit,
			finalHref: u.href,
			finalSearch: u.search
		});
	`

	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	data := result.Export().(map[string]interface{})

	if data["initialLimit"] != "10" {
		t.Errorf("Expected initial limit '10', got %v", data["initialLimit"])
	}

	expectedHref := "https://example.com/api?limit=50&sort=asc"
	if data["finalHref"] != expectedHref {
		t.Errorf("Expected href '%s', got %v", expectedHref, data["finalHref"])
	}

	expectedSearch := "?limit=50&sort=asc"
	if data["finalSearch"] != expectedSearch {
		t.Errorf("Expected search '%s', got %v", expectedSearch, data["finalSearch"])
	}
}

func TestURL_ThrowsOnInvalidURL(t *testing.T) {
	vm := goja.New()
	RegisterURL(vm)

	jsCode := `
		let caught = false;
		try {
			// Provide an invalid URL (missing scheme, unparseable by net/url)
			new URL(':::not-a-url:::');
		} catch (e) {
			caught = true;
		}
		caught;
	`

	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	if result.Export().(bool) != true {
		t.Error("Expected an invalid URL to throw a JavaScript error")
	}
}