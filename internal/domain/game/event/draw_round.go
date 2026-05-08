package event

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"

type DrawRound struct {
	reason  string
	tenpais *[common.NumPlayers]bool
	deltas  *[common.NumPlayers]int
	scores  *[common.NumPlayers]int
}

func NewDrawRound(
	reason string,
	tenpais *[common.NumPlayers]bool,
	deltas *[common.NumPlayers]int,
	scores *[common.NumPlayers]int,
) *DrawRound {
	return &DrawRound{
		reason:  reason,
		tenpais: tenpais,
		deltas:  deltas,
		scores:  scores,
	}
}

func (*DrawRound) isEvent() {}

func (d *DrawRound) Reason() string {
	return d.reason
}

func (d *DrawRound) Tenpais() *[common.NumPlayers]bool {
	return d.tenpais
}

func (d *DrawRound) Deltas() *[common.NumPlayers]int {
	return d.deltas
}

func (d *DrawRound) Scores() *[common.NumPlayers]int {
	return d.scores
}
