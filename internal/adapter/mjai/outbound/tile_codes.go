package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"

func tileCodes2(tiles [2]tile.Tile) []string {
	return []string{tiles[0].String(), tiles[1].String()}
}

func tileCodes3(tiles [3]tile.Tile) []string {
	return []string{tiles[0].String(), tiles[1].String(), tiles[2].String()}
}

func tileCodes4(tiles [4]tile.Tile) []string {
	return []string{tiles[0].String(), tiles[1].String(), tiles[2].String(), tiles[3].String()}
}
