package ai

import (
	"slices"

	"github.com/go-json-experiment/json"
)

type HashMapEntry struct {
	Key   []float64
	Value float64
}

type HashMap map[string]HashMapEntry

func NewHashMap() HashMap {
	return make(map[string]HashMapEntry)
}

func keyToString(key []float64) string {
	jsonKey, err := json.Marshal(&key)
	if err != nil {
		panic(err)
	}
	return string(jsonKey)
}

func (h *HashMap) Set(key []float64, value float64) {
	keyStr := keyToString(key)
	(*h)[keyStr] = HashMapEntry{
		Key:   slices.Clone(key),
		Value: value,
	}
}

func (h *HashMap) Get(key []float64, def float64) float64 {
	keyStr := keyToString(key)
	if entry, ok := (*h)[keyStr]; ok {
		return entry.Value
	}
	return def
}

func (h *HashMap) HasKey(key []float64) bool {
	keyStr := keyToString(key)
	_, ok := (*h)[keyStr]
	return ok
}

func (h *HashMap) ForEach(callback func(key []float64, value float64)) {
	for _, entry := range *h {
		callback(entry.Key, entry.Value)
	}
}
