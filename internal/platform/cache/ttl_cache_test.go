package cache_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/recco-app/recco/internal/platform/cache"
)

func TestTTLCache_BasicOperations(t *testing.T) {
	c := cache.NewTTLCache[string, string](100 * time.Millisecond)

	// Set and Get
	c.Set("key1", "val1")
	val, ok := c.Get("key1")
	if !ok || val != "val1" {
		t.Fatalf("expected val1, got %s (ok=%v)", val, ok)
	}

	// Wait for TTL expiration
	time.Sleep(120 * time.Millisecond)
	_, ok = c.Get("key1")
	if ok {
		t.Fatalf("expected key1 to expire")
	}

	// Delete
	c.Set("key2", "val2")
	c.Delete("key2")
	_, ok = c.Get("key2")
	if ok {
		t.Fatalf("expected key2 to be deleted")
	}
}

func TestTTLCache_ConcurrentAccess(t *testing.T) {
	c := cache.NewTTLCache[string, int](500 * time.Millisecond)
	var wg sync.WaitGroup

	// Run 100 concurrent writers and readers
	for i := 0; i < 100; i++ {
		wg.Add(2)
		val := i
		key := fmt.Sprintf("k-%d", i%10)

		go func() {
			defer wg.Done()
			c.Set(key, val)
		}()

		go func() {
			defer wg.Done()
			_, _ = c.Get(key)
		}()
	}

	wg.Wait()
}
