package player

import "fmt"

type PlayerID struct {
	value int
}

func NewPlayerID(id int) (*PlayerID, error) {
	if id < 0 || 3 < id {
		return nil, fmt.Errorf("invalid player id: %d", id)
	}

	return &PlayerID{value: id}, nil
}

func MustPlayerID(id int) *PlayerID {
	pid, err := NewPlayerID(id)
	if err != nil {
		panic(err)
	}
	return pid
}

func (id *PlayerID) Index() int {
	return id.value
}
