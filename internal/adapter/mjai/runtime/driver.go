package mjairuntime

import (
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/application"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

type Driver struct {
	name  string
	room  string
	agent ai.Agent
	bot   *application.Bot
	ended bool
	log   io.Writer
}

func NewDriver(name string, room string, agent ai.Agent, log io.Writer) *Driver {
	return &Driver{
		name:  name,
		room:  room,
		agent: agent,
		log:   log,
	}
}

func (d *Driver) Handle(msg inbound.Message) (outbound.Message, error) {
	switch msg := msg.(type) {
	case *inbound.Hello:
		return outbound.NewJoin(d.name, d.room), nil
	case *inbound.StartGame:
		self, err := seat.NewSeat(msg.ID)
		if err != nil {
			return nil, err
		}
		d.bot = application.NewBot(self, d.agent, newRoundStateReporter(d.log))
		d.ended = false
		return nil, nil
	case *inbound.EndGame:
		d.bot = nil
		d.ended = true
		return nil, nil
	case *inbound.Error:
		return nil, fmt.Errorf("server error: %s", msg.Message)
	default:
		if d.bot == nil {
			return nil, fmt.Errorf("cannot process %T: game has not started", msg)
		}
		ev, err := inbound.ParseEvent(msg)
		if err != nil {
			return nil, err
		}
		reaction, err := d.bot.Process(ev)
		if err != nil {
			return nil, err
		}
		if reaction.Kind() != application.ReactionAction {
			return nil, nil
		}
		return outbound.ToMessage(reaction.Action(), reaction.Log())
	}
}

func (d *Driver) Ended() bool {
	return d.ended
}
