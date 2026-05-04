package application

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

type Bot struct {
	self         seat.Seat
	agent        ai.Agent
	gameState    *game.State
	currentRound *round.State
	reporter     RoundStateReporter
}

type RoundStateReporter interface {
	ReportRoundState(state round.BoardRenderer) error
}

func NewBot(self seat.Seat, agent ai.Agent, reporter RoundStateReporter) *Bot {
	return &Bot{
		self:      self,
		agent:     agent,
		reporter:  reporter,
		gameState: game.NewDefaultState(),
	}
}

func (b *Bot) Process(ev event.Event) (Reaction, error) {
	switch ev := ev.(type) {
	case *event.StartRound:
		return b.processStartRound(ev)
	case *event.EndRound:
		return b.processEndRound()
	default:
		return b.processRoundEvent(ev)
	}
}

func (b *Bot) processStartRound(ev *event.StartRound) (Reaction, error) {
	currentRound, err := round.NewState(ev, b.gameState.Scores())
	if err != nil {
		return Reaction{}, err
	}
	b.currentRound = currentRound
	b.gameState.UpdateScores(currentRound.Scores())
	if err := b.reportRoundState(); err != nil {
		return Reaction{}, err
	}
	return NewNoReaction(), nil
}

func (b *Bot) processRoundEvent(ev event.Event) (Reaction, error) {
	if b.currentRound == nil {
		return Reaction{}, fmt.Errorf("cannot process %T: round has not started", ev)
	}
	if err := b.currentRound.Apply(ev); err != nil {
		return Reaction{}, err
	}
	if err := b.reportRoundState(); err != nil {
		return Reaction{}, err
	}

	legalActions, err := b.currentRound.LegalActions(b.self)
	if err != nil {
		return Reaction{}, err
	}
	if len(legalActions) == 0 {
		return NewNoReaction(), nil
	}

	decision, err := b.agent.Decide(ai.Request{
		Self:  b.self,
		Round: b.currentRound,
	})
	if err != nil {
		return Reaction{}, err
	}
	return NewActionReaction(decision.Action, decision.Log), nil
}

func (b *Bot) processEndRound() (Reaction, error) {
	if b.currentRound != nil {
		if err := b.reportRoundState(); err != nil {
			return Reaction{}, err
		}
		b.gameState.UpdateScores(b.currentRound.Scores())
	}
	b.currentRound = nil
	return NewNoReaction(), nil
}

func (b *Bot) reportRoundState() error {
	if b.reporter == nil || b.currentRound == nil {
		return nil
	}
	return b.reporter.ReportRoundState(b.currentRound)
}
