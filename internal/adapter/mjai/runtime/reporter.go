package mjairuntime

import (
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
)

type reporter struct {
	w io.Writer
}

func newReporter(w io.Writer) *reporter {
	if w == nil {
		return nil
	}
	return &reporter{w: w}
}

func (r *reporter) ReportRoundState(state round.BoardRenderer) error {
	if r == nil || r.w == nil {
		return nil
	}
	_, err := fmt.Fprint(r.w, state.RenderBoard())
	return err
}

func (r *reporter) ReportDecisionTrace(trace string) error {
	if r == nil || r.w == nil || trace == "" {
		return nil
	}
	_, err := fmt.Fprint(r.w, trace)
	return err
}
