package core

type HashMapKey interface {
	float64 | [4]float64
}

type HashMap[T HashMapKey] map[T]float64

func NewHashMap[T HashMapKey]() HashMap[T] {
	return make(map[T]float64)
}

func (h *HashMap[T]) Set(key T, value float64) {
	(*h)[key] = value
}

func (h *HashMap[T]) Get(key T, def float64) float64 {
	if v, ok := (*h)[key]; ok {
		return v
	}
	return def
}

func (h *HashMap[T]) HasKey(key T) bool {
	_, ok := (*h)[key]
	return ok
}

func (h *HashMap[T]) ForEach(callback func(key T, value float64)) {
	for k, v := range *h {
		callback(k, v)
	}
}
