package block_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/service/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewBlockFromMeld(t *testing.T) {
	tests := []struct {
		name string
		m    meld.Meld
		want block.Block
	}{
		{
			name: "1-23m",
			m: meld.MustChii(
				*tile.MustTileFromCode("1m"),
				[2]tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m")},
				*player.MustID(0),
			),
			want: block.MustSequence(*tile.MustTileFromCode("1m")),
		},
		{
			name: "5r-46p",
			m: meld.MustChii(
				*tile.MustTileFromCode("5pr"),
				[2]tile.Tile{*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("6p")},
				*player.MustID(0),
			),
			want: block.MustSequence(*tile.MustTileFromCode("4p")),
		},
		{
			name: "6-5r4p",
			m: meld.MustChii(
				*tile.MustTileFromCode("6p"),
				[2]tile.Tile{*tile.MustTileFromCode("5pr"), *tile.MustTileFromCode("4p")},
				*player.MustID(0),
			),
			want: block.MustSequence(*tile.MustTileFromCode("4p")),
		},
		{
			name: "7-5r6s",
			m: meld.MustChii(
				*tile.MustTileFromCode("7s"),
				[2]tile.Tile{*tile.MustTileFromCode("5sr"), *tile.MustTileFromCode("6s")},
				*player.MustID(0),
			),
			want: block.MustSequence(*tile.MustTileFromCode("5s")),
		},
		{
			name: "1m-1m1m to 1m triplet",
			m: meld.MustPon(
				*tile.MustTileFromCode("1m"),
				[2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
				*player.MustID(0),
			),
			want: block.MustTriplet(*tile.MustTileFromCode("1m")),
		},
		{
			name: "5sr-5s5s to 5s triplet",
			m: meld.MustPon(
				*tile.MustTileFromCode("5sr"),
				[2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
				*player.MustID(0),
			),
			want: block.MustTriplet(*tile.MustTileFromCode("5s")),
		},
		{
			name: "5s-5s5sr to 5s triplet",
			m: meld.MustPon(
				*tile.MustTileFromCode("5s"),
				[2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")},
				*player.MustID(0),
			),
			want: block.MustTriplet(*tile.MustTileFromCode("5s")),
		},
		{
			name: "1m-1m1m1m to 1m quad",
			m: meld.MustCalledKan(
				*tile.MustTileFromCode("1m"),
				[3]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
				*player.MustID(0),
			),
			want: block.MustQuad(*tile.MustTileFromCode("1m")),
		},
		{
			name: "5sr-5s5s5s to 5s quad",
			m: meld.MustCalledKan(
				*tile.MustTileFromCode("5sr"),
				[3]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
				*player.MustID(0),
			),
			want: block.MustQuad(*tile.MustTileFromCode("5s")),
		},
		{
			name: "5s-5s5s5sr to 5s quad",
			m: meld.MustCalledKan(
				*tile.MustTileFromCode("5s"),
				[3]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")},
				*player.MustID(0),
			),
			want: block.MustQuad(*tile.MustTileFromCode("5s")),
		},
		{
			name: "1m1m1m1m to 1m quad",
			m: meld.MustConcealedKan(
				[4]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			),
			want: block.MustQuad(*tile.MustTileFromCode("1m")),
		},
		{
			name: "5s5s5s5sr to 5s quad",
			m: meld.MustConcealedKan(
				[4]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")},
			),
			want: block.MustQuad(*tile.MustTileFromCode("5s")),
		},
		{
			name: "1p-1p1p-1p to 1p quad",
			m: meld.MustPromotedKan(
				*tile.MustTileFromCode("1p"),
				[2]tile.Tile{*tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p")},
				*tile.MustTileFromCode("1p"),
				*player.MustID(0),
			),
			want: block.MustQuad(*tile.MustTileFromCode("1p")),
		},
		{
			name: "5p-5p5p-5pr to 5p quad",
			m: meld.MustPromotedKan(
				*tile.MustTileFromCode("5p"),
				[2]tile.Tile{*tile.MustTileFromCode("5p"), *tile.MustTileFromCode("5p")},
				*tile.MustTileFromCode("5pr"),
				*player.MustID(0),
			),
			want: block.MustQuad(*tile.MustTileFromCode("5p")),
		},
		{
			name: "5pr-5p5p-5p to 5p quad",
			m: meld.MustPromotedKan(
				*tile.MustTileFromCode("5pr"),
				[2]tile.Tile{*tile.MustTileFromCode("5p"), *tile.MustTileFromCode("5p")},
				*tile.MustTileFromCode("5p"),
				*player.MustID(0),
			),
			want: block.MustQuad(*tile.MustTileFromCode("5p")),
		},
		{
			name: "5p-5p5pr-5p to 5p quad",
			m: meld.MustPromotedKan(
				*tile.MustTileFromCode("5p"),
				[2]tile.Tile{*tile.MustTileFromCode("5p"), *tile.MustTileFromCode("5pr")},
				*tile.MustTileFromCode("5p"),
				*player.MustID(0),
			),
			want: block.MustQuad(*tile.MustTileFromCode("5p")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := block.NewBlockFromMeld(tt.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlockFromMeld() = %v, want %v", got, tt.want)
			}
		})
	}
}
