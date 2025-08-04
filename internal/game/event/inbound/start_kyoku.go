package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type StartKyoku struct {
	Bakaze     base.Pai
	Kyoku      int `validate:"min=1,max=4"`
	Honba      int `validate:"min=0"`
	Kyotaku    int `validate:"min=0"`
	Oya        int `validate:"min=0,max=3"`
	DoraMarker base.Pai
	Scores     *[4]int
	Tehais     [4][13]base.Pai
}

func NewStartKyoku(
	bakaze base.Pai,
	kyoku int,
	honba int,
	kyotaku int,
	oya int,
	doraMarker base.Pai,
	scores *[4]int,
	tehais [4][13]base.Pai,
) (*StartKyoku, error) {
	s := &StartKyoku{
		Bakaze:     bakaze,
		Kyoku:      kyoku,
		Honba:      honba,
		Kyotaku:    kyotaku,
		Oya:        oya,
		DoraMarker: doraMarker,
		Scores:     scores,
		Tehais:     tehais,
	}

	isKazehai := s.Bakaze.IsTsupai() && (s.Bakaze.Number() <= 4)
	if !isKazehai {
		return nil, fmt.Errorf("invalid bakaze: %s", s.Bakaze.ToString())
	}

	if s.DoraMarker.IsUnknown() {
		return nil, fmt.Errorf("dora marker cannot be unknown: %s", s.DoraMarker.ToString())
	}

	if err := eventValidator.Struct(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *StartKyoku) isInboundEvent() {}
