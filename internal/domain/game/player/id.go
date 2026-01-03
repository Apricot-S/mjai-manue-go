package player

import "fmt"

type ID struct {
	value int
}

func NewID(id int) (*ID, error) {
	if id < 0 || 3 < id {
		return nil, fmt.Errorf("invalid player id: %d", id)
	}

	return &ID{value: id}, nil
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
