package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func containsDiscard(actions []action.Action, tileCode string, tsumogiri bool) bool {
	for _, a := range actions {
		discard, ok := a.(*action.Discard)
		if !ok {
			continue
		}
		if discard.Tile().String() == tileCode && discard.Tsumogiri() == tsumogiri {
			return true
		}
	}
	return false
}

func containsRiichi(actions []action.Action, actor seat.Seat) bool {
	for _, a := range actions {
		riichi, ok := a.(*action.Riichi)
		if !ok {
			continue
		}
		if riichi.Actor() == actor {
			return true
		}
	}
	return false
}

func containsKyushukyuhai(actions []action.Action, actor seat.Seat) bool {
	for _, a := range actions {
		kyushukyuhai, ok := a.(*action.Kyushukyuhai)
		if !ok {
			continue
		}
		if kyushukyuhai.Actor() == actor {
			return true
		}
	}
	return false
}

func containsWin(actions []action.Action, actor seat.Seat, target seat.Seat, winningTileCode string) bool {
	for _, a := range actions {
		win, ok := a.(*action.Win)
		if !ok {
			continue
		}
		if win.Actor() == actor && win.Target() == target && win.WinningTile().String() == winningTileCode {
			return true
		}
	}
	return false
}

func containsPon(actions []action.Action, actor seat.Seat, target seat.Seat, takenCode string, consumedCodes [2]string) bool {
	for _, a := range actions {
		pon, ok := a.(*action.Pon)
		if !ok {
			continue
		}
		if pon.Actor() != actor || pon.Target() != target || pon.Taken().String() != takenCode {
			continue
		}
		consumed := pon.Consumed()
		if consumed[0].String() == consumedCodes[0] &&
			consumed[1].String() == consumedCodes[1] {
			return true
		}
	}
	return false
}

func containsCalledKan(actions []action.Action, actor seat.Seat, target seat.Seat, takenCode string, consumedCodes [3]string) bool {
	for _, a := range actions {
		calledKan, ok := a.(*action.CalledKan)
		if !ok {
			continue
		}
		if calledKan.Actor() != actor || calledKan.Target() != target || calledKan.Taken().String() != takenCode {
			continue
		}
		consumed := calledKan.Consumed()
		if consumed[0].String() == consumedCodes[0] &&
			consumed[1].String() == consumedCodes[1] &&
			consumed[2].String() == consumedCodes[2] {
			return true
		}
	}
	return false
}

func containsPass(actions []action.Action, actor seat.Seat) bool {
	for _, a := range actions {
		pass, ok := a.(*action.Pass)
		if !ok {
			continue
		}
		if pass.Actor() == actor {
			return true
		}
	}
	return false
}

func containsPromotedKan(actions []action.Action, actor seat.Seat, addedCode string, consumedCodes [3]string) bool {
	for _, a := range actions {
		promotedKan, ok := a.(*action.PromotedKan)
		if !ok {
			continue
		}
		if promotedKan.Actor() != actor || promotedKan.Added().String() != addedCode {
			continue
		}
		consumed := promotedKan.Consumed()
		if consumed[0].String() == consumedCodes[0] &&
			consumed[1].String() == consumedCodes[1] &&
			consumed[2].String() == consumedCodes[2] {
			return true
		}
	}
	return false
}

func containsConcealedKan(actions []action.Action, actor seat.Seat, consumedCodes [4]string) bool {
	for _, a := range actions {
		concealedKan, ok := a.(*action.ConcealedKan)
		if !ok {
			continue
		}
		if concealedKan.Actor() != actor {
			continue
		}
		consumed := concealedKan.Consumed()
		if consumed[0].String() == consumedCodes[0] &&
			consumed[1].String() == consumedCodes[1] &&
			consumed[2].String() == consumedCodes[2] &&
			consumed[3].String() == consumedCodes[3] {
			return true
		}
	}
	return false
}

func unknownHandForLegalActionsTest() [common.InitHandSize]tile.Tile {
	var hand [common.InitHandSize]tile.Tile
	for i := range hand {
		hand[i] = tile.MustTileFromCode("?")
	}
	return hand
}

func kyushukyuhaiHandForTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("9m"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("9p"),
		tile.MustTileFromCode("1s"),
		tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("E"),
		tile.MustTileFromCode("S"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("6m"),
	}
}

func menzenTsumoHandForLegalActionsTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("9s"),
	}
}

func openChiiHandForLegalActionsTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("5p"),
		tile.MustTileFromCode("6p"),
		tile.MustTileFromCode("7s"),
		tile.MustTileFromCode("8s"),
	}
}

func stateAfterChiiForLegalActionsTest(t *testing.T, s *State) *State {
	t.Helper()

	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("2m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewChii(actor, target, taken, [2]tile.Tile{
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
	})); err != nil {
		t.Fatalf("Apply(Chii) failed: %v", err)
	}
	return s
}

func openPlayerWithoutYakuTsumoForLegalActionsTest(t *testing.T) player.Player {
	t.Helper()

	handTiles := [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("E"),
	}
	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	pon := meld.MustPon(
		tile.MustTileFromCode("1m"),
		[2]tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m")},
		seat.MustSeat(1),
	)
	if err := p.Pon(*pon); err != nil {
		t.Fatalf("Pon() failed: %v", err)
	}
	if err := p.Discard(tile.MustTileFromCode("E"), false); err != nil {
		t.Fatalf("Discard() failed: %v", err)
	}
	if err := p.Draw(tile.MustTileFromCode("9s")); err != nil {
		t.Fatalf("Draw() failed: %v", err)
	}
	return p
}

func newStateWithOpenNoYakuTsumoForLegalActionsTest(t *testing.T, numLeftTiles int, lastDrawWasReplacement bool) *State {
	t.Helper()

	actor := seat.MustSeat(0)
	players := [common.NumPlayers]player.Player{
		openPlayerWithoutYakuTsumoForLegalActionsTest(t),
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
	}
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		actor,
		actor,
		tile.Tiles{tile.MustTileFromCode("E")},
		numLeftTiles,
		players,
	)
	s.pendingDiscard = &actor
	s.lastDrawWasReplacement = lastDrawWasReplacement
	return &s
}

func newStateBeforeConcealedKanForLegalActionsTest(t *testing.T, numLeftTiles int, numKans int) *State {
	t.Helper()

	hands := newValidHands()
	hands[0] = concealedKanHandForTest()
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		seat.MustSeat(0),
		seat.MustSeat(0),
		tile.Tiles{tile.MustTileFromCode("E")},
		numLeftTiles+1,
		newVisiblePlayersForTest(t, hands),
	)
	actor := seat.MustSeat(0)
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("E"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	s.numKans = numKans
	return &s
}

func newRiichiAcceptedStateBeforeConcealedKanForTest(t *testing.T) *State {
	t.Helper()

	hands := newValidHands()
	hands[0] = [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("7s"),
		tile.MustTileFromCode("8s"),
		tile.MustTileFromCode("E"),
		tile.MustTileFromCode("E"),
		tile.MustTileFromCode("W"),
	}
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("9s"))); err != nil {
		t.Fatalf("Apply(first Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("W"), false)); err != nil {
		t.Fatalf("Apply(first Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichiAccepted(actor, nil, nil)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}
	for i := 1; i < 4; i++ {
		other := seat.MustSeat(i)
		drawnTile := tile.MustTileFromCode("6m")
		if err := s.Apply(event.NewDraw(other, drawnTile)); err != nil {
			t.Fatalf("Apply(other Draw %d) failed: %v", i, err)
		}
		if err := s.Apply(event.NewDiscard(other, drawnTile, true)); err != nil {
			t.Fatalf("Apply(other Discard %d) failed: %v", i, err)
		}
	}
	return s
}

func newRiichiAcceptedStateBeforeConcealedKanChangingOnlyWinningFormForTest(t *testing.T) *State {
	t.Helper()

	hands := newValidHands()
	hands[0] = [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("W"),
	}
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("E"))); err != nil {
		t.Fatalf("Apply(first Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("W"), false)); err != nil {
		t.Fatalf("Apply(first Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichiAccepted(actor, nil, nil)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}
	for i := 1; i < 4; i++ {
		other := seat.MustSeat(i)
		drawnTile := tile.MustTileFromCode("6m")
		if err := s.Apply(event.NewDraw(other, drawnTile)); err != nil {
			t.Fatalf("Apply(other Draw %d) failed: %v", i, err)
		}
		if err := s.Apply(event.NewDiscard(other, drawnTile, true)); err != nil {
			t.Fatalf("Apply(other Discard %d) failed: %v", i, err)
		}
	}
	return s
}
