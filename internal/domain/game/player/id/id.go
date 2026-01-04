package id

import "fmt"

type ID struct {
	value int
}

func NewID(index int) (*ID, error) {
	if index < 0 || 3 < index {
		return nil, fmt.Errorf("invalid player id: %d", index)
	}

	return &ID{value: index}, nil
}

func MustID(id int) *ID {
	pid, err := NewID(id)
	if err != nil {
		panic(err)
	}
	return pid
}

func (id *ID) Index() int {
	return id.value
}
