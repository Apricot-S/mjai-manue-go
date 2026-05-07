package seat

import "fmt"

type Seat struct {
	index int
}

func NewSeat(index int) (Seat, error) {
	if index < 0 || 3 < index {
		return Seat{}, fmt.Errorf("invalid player seat: %d", index)
	}

	return Seat{index: index}, nil
}

func MustSeat(id int) Seat {
	s, err := NewSeat(id)
	if err != nil {
		panic(err)
	}
	return s
}

func (seat Seat) Index() int {
	return seat.index
}

// IsShimochaOf reports whether seat is the player to the left of target.
func (seat Seat) IsShimochaOf(target Seat) bool {
	return seat.index == (target.index+1)%4
}
