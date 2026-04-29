package mjairuntime

import (
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
)

type roundStateReporter struct {
	w io.Writer
}

func newRoundStateReporter(w io.Writer) *roundStateReporter {
	if w == nil {
		return nil
	}
	return &roundStateReporter{w: w}
}

func (r *roundStateReporter) ReportRoundState(state round.BoardRenderer) error {
	if r == nil || r.w == nil {
		return nil
	}
	_, err := fmt.Fprint(r.w, state.RenderBoard())
	return err
}
