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

func (s Seat) Index() int {
	return s.index
}

// DistanceFrom returns seat's relative position from base as a value in 0..3.
//
// A return value of 0 means the same seat as base, 1 means shimocha, 2 means
// toimen, and 3 means kamicha. If base is the starting dealer seat 0, the
// return value matches seat.Index().
func (s Seat) DistanceFrom(base Seat) int {
	return (s.index - base.index + 4) % 4
}

// IsShimochaOf reports whether seat is the player to the left of target.
func (s Seat) IsShimochaOf(target Seat) bool {
	return s.index == (target.index+1)%4
}
