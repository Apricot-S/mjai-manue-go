package inbound

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func parseOpenCallFields(actorValue, targetValue int, pai string, consumedValue []string) (
	actorSeat, targetSeat *seat.Seat,
	taken *tile.Tile,
	consumed [2]tile.Tile,
	err error,
) {
	actor, err := parseSeatField("actor", actorValue)
	if err != nil {
		return nil, nil, nil, [2]tile.Tile{}, err
	}
	target, err := parseSeatField("target", targetValue)
	if err != nil {
		return nil, nil, nil, [2]tile.Tile{}, err
	}
	takenTile, err := parseTileField("pai", pai)
	if err != nil {
		return nil, nil, nil, [2]tile.Tile{}, err
	}
	consumed, err = parseConsumed2(consumedValue)
	if err != nil {
		return nil, nil, nil, [2]tile.Tile{}, err
	}
	return actor, target, takenTile, consumed, nil
}
