package wind

import (
	"fmt"
)

type Wind int

const (
	East = iota + 1
	South
	West
	North
)

func NewWind(w string) (Wind, error) {
	switch w {
	case "E":
		return East, nil
	case "S":
		return South, nil
	case "W":
		return West, nil
	case "N":
		return North, nil
	default:
		return 0, fmt.Errorf("invalid wind string %q", w)
	}
}

func MustWind(w string) Wind {
	ret, err := NewWind(w)
	if err != nil {
		panic(err)
	}
	return ret
}

func (w Wind) String() string {
	switch w {
	case East:
		return "E"
	case South:
		return "S"
	case West:
		return "W"
	case North:
		return "N"
	default:
		panic(fmt.Sprintf("wind: invalid value %d", w))
	}
}
