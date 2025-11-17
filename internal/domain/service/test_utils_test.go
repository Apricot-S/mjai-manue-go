package service_test

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func codesToHand(codes []string) *hand.Hand {
	tiles := make([]tile.Tile, len(codes))
	for i, code := range codes {
		tiles[i] = *tile.MustTileFromCode(code)
	}

	h, err := hand.NewHand(tiles)
	if err != nil {
		panic(err)
	}

	return h
}
