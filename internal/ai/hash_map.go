package ai

import (
	"github.com/go-json-experiment/json"
)

type HashMapKey interface {
	~float64 | ~[]float64
}

type HashMapEntry[K HashMapKey] struct {
	Key   K
	Value float64
}

type HashMap[K HashMapKey] struct {
	data map[string]HashMapEntry[K]
}

func NewHashMap[K HashMapKey]() *HashMap[K] {
	return &HashMap[K]{
		data: make(map[string]HashMapEntry[K]),
	}
}

func keyToString[K HashMapKey](key K) string {
	jsonKey, err := json.Marshal(&key)
	if err != nil {
		panic(err)
	}
	return string(jsonKey)
}

func (h *HashMap[K]) Set(key K, value float64) {
	keyStr := keyToString(key)
	h.data[keyStr] = HashMapEntry[K]{
		Key:   key,
		Value: value,
	}
}

func (h *HashMap[K]) Get(key K, def float64) float64 {
	keyStr := keyToString(key)
	if entry, ok := h.data[keyStr]; ok {
		return entry.Value
	}
	return def
}

func (h *HashMap[K]) HasKey(key K) bool {
	keyStr := keyToString(key)
	_, ok := h.data[keyStr]
	return ok
}

func (h *HashMap[K]) ForEach(callback func(key K, value float64)) {
	for _, entry := range h.data {
		callback(entry.Key, entry.Value)
	}
}
