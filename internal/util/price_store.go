package util

import "sync"

var symbolPriceMap *SymbolPriceMap

type SymbolPriceMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func GetSymbolPriceMap() *SymbolPriceMap {
	if symbolPriceMap == nil {
		symbolPriceMap = &SymbolPriceMap{
			data: make(map[string]string),
		}
	}

	return symbolPriceMap
}

// Set adds or updates a key-value pair in the map.
func (spm *SymbolPriceMap) Set(key string, value string) {
	spm.mu.Lock()         // Acquire write lock
	defer spm.mu.Unlock() // Release write lock when function exits
	spm.data[key] = value
}

// Get retrieves a value from the map.
func (spm *SymbolPriceMap) Get(key string) (string, bool) {
	spm.mu.RLock()         // Acquire read lock
	defer spm.mu.RUnlock() // Release read lock when function exits
	val, ok := spm.data[key]
	return val, ok
}

// GetAll retrieves all key values from the map.
func (spm *SymbolPriceMap) GetAll() map[string]string {
	spm.mu.RLock()         // Acquire read lock
	defer spm.mu.RUnlock() // Release read lock when function exits
	return spm.data
}

// DeleteAll deletes all key values from the map.
func (spm *SymbolPriceMap) DeleteAll() {
	spm.mu.RLock()         // Acquire read lock
	defer spm.mu.RUnlock() // Release read lock when function exits
	symbolPriceMap = nil
}
