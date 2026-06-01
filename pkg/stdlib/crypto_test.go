package stdlib

import (
	"testing"

	"github.com/dop251/goja"
)

func TestCrypto_RandomUUID(t *testing.T) {
	vm := goja.New()
	LazyInject(vm, RuntimeConfig{})

	jsCode := `
		const id1 = crypto.randomUUID();
		const id2 = crypto.randomUUID();
		
		({ id1, id2 });
	`

	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	data := result.Export().(map[string]interface{})
	id1 := data["id1"].(string)
	id2 := data["id2"].(string)

	if id1 == id2 {
		t.Errorf("Expected UUIDs to be unique, got duplicate: %s", id1)
	}

	if len(id1) != 36 {
		t.Errorf("Expected UUID length of 36, got %d", len(id1))
	}
}

func TestCrypto_GetRandomValues(t *testing.T) {
	vm := goja.New()
	LazyInject(vm, RuntimeConfig{})

	jsCode := `
		// Create an empty array of 10 zeros
		const arr = new Uint8Array(10);
		
		// Fill it with entropy
		const returnedArr = crypto.getRandomValues(arr);
		
		// Sum the array. If it's still 0, it failed to inject entropy.
		let sum = 0;
		for (let i = 0; i < arr.length; i++) {
			sum += arr[i];
		}
		
		// Did it return the exact same array instance as per the spec?
		const isSameInstance = (arr === returnedArr);
		
		({ sum, isSameInstance });
	`

	result, err := vm.RunString(jsCode)
	if err != nil {
		t.Fatalf("JS execution failed: %v", err)
	}

	data := result.Export().(map[string]interface{})

	sum := int(data["sum"].(int64))
	if sum == 0 {
		t.Errorf("Expected array sum to be > 0, got 0. Entropy was not injected.")
	}

	if data["isSameInstance"] != true {
		t.Errorf("Expected getRandomValues to return the exact same array instance")
	}
}

func TestCrypto_GetRandomValuesQuota(t *testing.T) {
	vm := goja.New()
	LazyInject(vm, RuntimeConfig{})

	jsCode := `
		let caught = false;
		try {
			// Request 70,000 values (over the 65,536 limit)
			const bigArr = new Uint8Array(70000);
			crypto.getRandomValues(bigArr);
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
		t.Error("Expected QuotaExceededError for arrays larger than 65536 bytes")
	}
}