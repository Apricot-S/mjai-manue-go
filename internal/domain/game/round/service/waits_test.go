package service_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestWaitsFor_NotTenpai(t *testing.T) {
	h := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"), tile.MustTileFromCode("3m"), tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"), tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("E"), tile.MustTileFromCode("S"),
	})

	got := service.WaitsFor(h)

	assertWaitsForTest(t, got)
}

func TestWaitsFor_SevenPairs(t *testing.T) {
	h := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"), tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3p"), tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"), tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("5s"), tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"), tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("E"),
	})

	got := service.WaitsFor(h)

	assertWaitsForTest(t, got, "E")
}

func TestWaitsFor_ThirteenSidedKokushiMusou(t *testing.T) {
	h := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("9m"),
		tile.MustTileFromCode("1p"), tile.MustTileFromCode("9p"),
		tile.MustTileFromCode("1s"), tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("E"), tile.MustTileFromCode("S"), tile.MustTileFromCode("W"), tile.MustTileFromCode("N"),
		tile.MustTileFromCode("P"), tile.MustTileFromCode("F"), tile.MustTileFromCode("C"),
	})

	got := service.WaitsFor(h)

	assertWaitsForTest(t, got, "1m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C")
}

func assertWaitsForTest(t *testing.T, got service.WaitSet, wantCodes ...string) {
	t.Helper()

	want := map[tile.Tile]bool{}
	for _, code := range wantCodes {
		want[tile.MustTileFromCode(code)] = true
	}

	for id := range tile.NumTileType34 {
		waitTile := tile.MustTileFromID(id)
		if got.Has(waitTile) != want[waitTile] {
			t.Errorf("WaitsFor().Has(%s) = %t, want %t", waitTile, got.Has(waitTile), want[waitTile])
		}
	}
}
