package service_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/wind"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/service"
)

func TestCalculateFuHan(t *testing.T) {
	tests := []struct {
		name           string
		handCodes      []string
		handBlocks     []block.Block
		melds          []meld.Meld
		prevalentWind  wind.Wind
		seatWind       wind.Wind
		doraIndicators []tile.Tile
		tsumo          bool
		riichi         bool
		wantFu         int
		wantHan        int
		wantYakus      map[string]int
	}{
		{
			name:      "no Yaku",
			handCodes: []string{"1m", "1m", "1m", "2p", "3p", "4p", "3s", "4s", "5s", "6s", "6s", "6s", "9s", "9s"},
			handBlocks: []block.Block{
				block.MustTriplet(*tile.MustTileFromCode("1m")),
				block.MustSequence(*tile.MustTileFromCode("2p")),
				block.MustSequence(*tile.MustTileFromCode("3s")),
				block.MustTriplet(*tile.MustTileFromCode("6s")),
				block.MustPair(*tile.MustTileFromCode("9s")),
			},
			melds:          nil,
			prevalentWind:  wind.East,
			seatWind:       wind.South,
			doraIndicators: nil,
			tsumo:          false,
			riichi:         false,
			wantFu:         40,
			wantHan:        0,
			wantYakus:      map[string]int{},
		},
		{
			name:      "only Riichi",
			handCodes: []string{"1m", "1m", "1m", "2p", "3p", "4p", "3s", "4s", "5s", "6s", "6s", "6s", "9s", "9s"},
			handBlocks: []block.Block{
				block.MustTriplet(*tile.MustTileFromCode("1m")),
				block.MustSequence(*tile.MustTileFromCode("2p")),
				block.MustSequence(*tile.MustTileFromCode("3s")),
				block.MustTriplet(*tile.MustTileFromCode("6s")),
				block.MustPair(*tile.MustTileFromCode("9s")),
			},
			melds:          nil,
			prevalentWind:  wind.East,
			seatWind:       wind.South,
			doraIndicators: nil,
			tsumo:          false,
			riichi:         true,
			wantFu:         40,
			wantHan:        1,
			wantYakus:      map[string]int{"reach": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.handCodes)
			fu, han, yakus := service.CalculateFuHan(hand, tt.handBlocks, tt.melds, tt.prevalentWind, tt.seatWind, tt.doraIndicators, tt.tsumo, tt.riichi)
			if tt.wantFu != fu {
				t.Errorf("CalculateFuHan() = %v, want %v", fu, tt.wantFu)
			}
			if tt.wantHan != han {
				t.Errorf("CalculateFuHan() = %v, want %v", han, tt.wantHan)
			}
			if !reflect.DeepEqual(tt.wantYakus, yakus) {
				t.Errorf("CalculateFuHan() = %v, want %v", yakus, tt.wantYakus)
			}
		})
	}
}
