package main

import (
	"encoding/json/v2"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol/mjai"
	"github.com/Apricot-S/mjai-manue-go/tools/shared"
)

type LawEvent struct {
	Logs [4]string `json:"logs"`
}

const playerName = "Manue014"

type Verifier struct {
	PlayerID int
	AI       ai.AI
	Decision outbound.Event
}

func NewVerifier(ai ai.AI) *Verifier {
	return &Verifier{AI: ai}
}

// ModifyAction is called before applying the action to the state.
// It can modify the action if needed.
func (v *Verifier) ModifyAction(action inbound.Event, g game.StateAnalyzer) (inbound.Event, error) {
	switch a := action.(type) {
	case *inbound.Error:
		return action, fmt.Errorf("error in the log: %+v", a)
	case *inbound.StartGame:
		v.PlayerID = slices.Index(a.Names[:], playerName)
		if v.PlayerID == -1 {
			return action, fmt.Errorf("player name %q not found in players %v", playerName, a.Names)
		}

		startGame, err := inbound.NewStartGame(v.PlayerID, a.Names)
		if err != nil {
			return action, err
		}
		action = startGame
	}

	decision, err := v.AI.DecideAction(g, v.PlayerID)
	if err != nil {
		return action, fmt.Errorf("failed to decide action: %w", err)
	}
	v.Decision = decision

	return action, nil
}

func (v *Verifier) VerifyAction(action inbound.Event, g game.StateViewer) (string, error) {
	switch a := action.(type) {
	case *inbound.Dahai:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		da, ok := v.Decision.(*outbound.Dahai)
		if !ok {
			return "", fmt.Errorf("expected dahai decision, got %#v", v.Decision)
		}
		if a.Pai != da.Pai || a.Tsumogiri != da.Tsumogiri {
			return fmt.Sprintf("dahai mismatch:\nport:\n%+v\n", da.Log), nil
		}
	case *inbound.Chi:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		ch, ok := v.Decision.(*outbound.Chi)
		if !ok {
			return "", fmt.Errorf("expected chi decision, got %#v", v.Decision)
		}
		if a.Taken != ch.Taken || a.Consumed != ch.Consumed {
			return fmt.Sprintf("chi mismatch:\nport:\n%+v\n", ch.Log), nil
		}
	case *inbound.Pon:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		po, ok := v.Decision.(*outbound.Pon)
		if !ok {
			return "", fmt.Errorf("expected pon decision, got %#v", v.Decision)
		}
		if a.Target != po.Target || a.Taken != po.Taken || a.Consumed != po.Consumed {
			return fmt.Sprintf("pon mismatch:\nport:\n%+v\n", po.Log), nil
		}
	case *inbound.Daiminkan:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		dm, ok := v.Decision.(*outbound.Daiminkan)
		if !ok {
			return "", fmt.Errorf("expected daiminkan decision, got %#v", v.Decision)
		}
		if a.Target != dm.Target || a.Taken != dm.Taken {
			return fmt.Sprintf("daiminkan mismatch:\nport:\n%+v\n", dm.Log), nil
		}
	case *inbound.Reach:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		_, ok := v.Decision.(*outbound.Reach)
		if !ok {
			return "", fmt.Errorf("expected reach decision, got %#v", v.Decision)
		}
	case *inbound.Hora:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		ho, ok := v.Decision.(*outbound.Hora)
		if !ok {
			return "", fmt.Errorf("expected hora decision, got %#v", v.Decision)
		}
		if a.Target != ho.Target || *a.Pai != ho.Pai {
			return fmt.Sprintf("hora mismatch:\nport:\n%+v\n", ho.Log), nil
		}
	}
	return "", nil
}

func run(args []string) error {
	paths, err := shared.GlobAll(args)
	if err != nil {
		return fmt.Errorf("error in glob: %w", err)
	}

	ai, err := ai.NewManueAIDefault()
	if err != nil {
		log.Fatalf("failed to create AI: %v", err)
	}
	verifier := NewVerifier(ai)

	archive := shared.NewArchive(paths, &mjai.MjaiAdapter{})
	var prevLog string
	separator := strings.Repeat("=", 122)

	onAction := func(action inbound.Event) error {
		action, err := verifier.ModifyAction(action, archive.StateAnalyzer())
		if err != nil {
			return err
		}
		if err := archive.StateUpdater().Update(action); err != nil {
			return err
		}
		detail, err := verifier.VerifyAction(action, archive.StateViewer())
		if err != nil {
			return err
		}

		if detail != "" {
			fmt.Print("VerifyAction mismatch:\n\n")
			fmt.Printf("state (after action):\n%s\n", archive.StateViewer().RenderBoard())
			fmt.Printf("%v", detail)
			fmt.Println("original:")
			fmt.Println(prevLog)
			fmt.Println(separator)
		}
		return nil
	}

	onRaw := func(raw []byte) error {
		logs := LawEvent{}
		if err := json.Unmarshal(raw, &logs); err != nil {
			return fmt.Errorf("failed to unmarshal raw json: %w", err)
		}

		for _, l := range logs.Logs {
			if l == "" {
				continue
			}
			prevLog = l
		}
		return nil
	}

	if err := archive.PlayLight(onAction, onRaw); err != nil {
		return fmt.Errorf("error in processing log: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <LOG_GLOB_PATTERNS>...\n", os.Args[0])
		os.Exit(2)
	}

	if err := run(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}
