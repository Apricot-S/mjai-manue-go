package player_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewVisiblePlayer(t *testing.T) {
	tests := []struct {
		name      string
		handTiles []tile.Tile
		wantHand  *hand.VisibleHand
		wantErr   bool
	}{
		{
			name:      "valid",
			handTiles: []tile.Tile{},
			wantHand:  hand.CodesToHand([]string{}),
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := player.NewVisiblePlayer(tt.handTiles)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewVisiblePlayer() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewVisiblePlayer() succeeded unexpectedly")
			}

			h, ok := got.Hand()
			if !ok {
				t.Errorf("NewVisiblePlayer() Hand() returned not ok")
			}
			if *h != *tt.wantHand {
				t.Errorf("NewVisiblePlayer().Hand() = %v, want %v", h, tt.wantHand)
			}

			ts := tile.Tiles(tt.wantHand.ToTiles())
			sort.Sort(ts)
			if !reflect.DeepEqual(got.HandTiles(), []tile.Tile(ts)) {
				t.Errorf("NewVisiblePlayer().HandTiles() = %v, want %v", got.HandTiles(), ts)
			}
		})
	}
}
