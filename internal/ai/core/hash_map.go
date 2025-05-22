package core

import "github.com/go-json-experiment/json"

type HashMapKey interface {
	float64 | []float64
}

type HashMapEntry[T HashMapKey] struct {
	Key   T
	Value float64
}

type HashMap[T HashMapKey] map[string]HashMapEntry[T]

func NewHashMap[T HashMapKey]() HashMap[T] {
	return make(map[string]HashMapEntry[T])
}

func keyToString[T HashMapKey](key T) string {
	jsonKey, err := json.Marshal(&key)
	if err != nil {
		panic(err)
	}
	return string(jsonKey)
}

func (h *HashMap[T]) Set(key T, value float64) {
	keyStr := keyToString(key)
	(*h)[keyStr] = HashMapEntry[T]{
		Key:   key,
		Value: value,
	}
}

func (h *HashMap[T]) Get(key T, def float64) float64 {
	keyStr := keyToString(key)
	if entry, ok := (*h)[keyStr]; ok {
		return entry.Value
	}
	return def
}

func (h *HashMap[T]) HasKey(key T) bool {
	keyStr := keyToString(key)
	_, ok := (*h)[keyStr]
	return ok
}

func (h *HashMap[T]) ForEach(callback func(key T, value float64)) {
	for _, entry := range *h {
		callback(entry.Key, entry.Value)
	}
}
