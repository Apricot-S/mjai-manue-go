package hand_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func TestNewHand(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tiles []tile.Tile
		want  []tile.Tile
	}{
		{
			name:  "empty hand",
			tiles: []tile.Tile{},
			want:  []tile.Tile{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hand.NewHand(tt.tiles)
			if !reflect.DeepEqual(got.ToTiles(), tt.tiles) {
				t.Errorf("NewHand() = %v, want %v", got, tt.want)
			}
		})
	}
}
