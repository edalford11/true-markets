package util

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestSymbolPriceMap_ConcurrentSetSameKey(t *testing.T) {
	spm := &SymbolPriceMap{
		data: make(map[string]string),
	}

	const numGoroutines = 100
	const numOperations = 100
	const key = "TESTKEY"

	var wg sync.WaitGroup
	var atomicCounter int64
	writtenValues := make(map[string]bool)
	var mapMutex sync.Mutex

	t.Run("Concurrent Set operations on same key - race condition validation", func(t *testing.T) {
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					value := "goroutine" + string(rune(48+id%10)) + "op" + string(rune(48+j%10))
					
					mapMutex.Lock()
					writtenValues[value] = true
					mapMutex.Unlock()
					
					spm.Set(key, value)
					atomic.AddInt64(&atomicCounter, 1)
				}
			}(i)
		}

		wg.Wait()

		finalValue, exists := spm.Get(key)
		if !exists {
			t.Errorf("Expected key %s to exist after concurrent operations", key)
		}

		mapMutex.Lock()
		isValidValue := writtenValues[finalValue]
		mapMutex.Unlock()

		if !isValidValue {
			t.Errorf("Final value '%s' was not one of the values written by any goroutine", finalValue)
		}

		expectedOperations := int64(numGoroutines * numOperations)
		actualOperations := atomic.LoadInt64(&atomicCounter)
		if actualOperations != expectedOperations {
			t.Errorf("Expected %d operations, but got %d", expectedOperations, actualOperations)
		}

		if len(spm.data) != 1 {
			t.Errorf("Expected exactly 1 key in map, but got %d", len(spm.data))
		}

		t.Logf("Successfully completed %d concurrent operations with final value: %s", actualOperations, finalValue)
	})
	
	t.Run("Concurrent read-write validation", func(t *testing.T) {
		testKey := "READ_WRITE_TEST"
		spm.Set(testKey, "initial")
		
		var readWg sync.WaitGroup
		var writeWg sync.WaitGroup
		const readers = 50
		const writers = 50
		
		readValues := make(chan string, readers*10)
		
		writeWg.Add(writers)
		readWg.Add(readers)
		
		for i := 0; i < writers; i++ {
			go func(id int) {
				defer writeWg.Done()
				value := "writer" + string(rune(48+id%10))
				spm.Set(testKey, value)
			}(i)
		}
		
		for i := 0; i < readers; i++ {
			go func() {
				defer readWg.Done()
				for j := 0; j < 10; j++ {
					if value, exists := spm.Get(testKey); exists {
						readValues <- value
					}
				}
			}()
		}
		
		writeWg.Wait()
		readWg.Wait()
		close(readValues)
		
		finalValue, exists := spm.Get(testKey)
		if !exists {
			t.Error("Expected test key to exist after concurrent read-write operations")
		}
		
		readCount := 0
		for range readValues {
			readCount++
		}
		
		if readCount == 0 {
			t.Error("Expected at least some successful reads during concurrent operations")
		}
		
		t.Logf("Concurrent read-write test completed. Final value: %s, Read operations: %d", finalValue, readCount)
	})
}

func TestGetSymbolPriceMap_Singleton(t *testing.T) {
	symbolPriceMap = nil

	spm1 := GetSymbolPriceMap()
	spm2 := GetSymbolPriceMap()

	if spm1 != spm2 {
		t.Error("GetSymbolPriceMap should return the same instance (singleton)")
	}

	spm1.Set("TEST", "123")
	value, exists := spm2.Get("TEST")
	if !exists || value != "123" {
		t.Error("Changes in one instance should be visible in another")
	}
}
