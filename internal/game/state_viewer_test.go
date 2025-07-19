package game

import (
	"reflect"
	"slices"
	"testing"
)

func getDefaultStateForTest() *StateImpl {
	east, _ := NewPaiWithName("E")
	players := [NumPlayers]Player{}

	for i := range NumPlayers {
		tehais := make([]Pai, 13)
		for j := range 13 {
			tehais[j] = *Unknown
		}

		players[i] = Player{
			id:                i,
			name:              "",
			tehais:            tehais,
			furos:             make([]Furo, 0, maxNumFuro),
			ho:                make([]Pai, 0),
			sutehais:          make([]Pai, 0),
			extraAnpais:       make([]Pai, 0),
			reachState:        NotReach,
			reachHoIndex:      -1,
			reachSutehaiIndex: -1,
			score:             InitScore,
			canDahai:          false,
			isMenzen:          true,
		}
	}

	return &StateImpl{
		players:     players,
		bakaze:      *east,
		kyokuNum:    1,
		honba:       0,
		oya:         &players[0],
		chicha:      &players[0],
		doraMarkers: make([]Pai, 0, MaxNumDoraMarkers),
		numPipais:   NumInitPipais,

		prevEventType:    "",
		prevDahaiActor:   -1,
		prevDahaiPai:     nil,
		currentEventType: "",
	}
}

func TestState_Anpais(t *testing.T) {
	type args struct {
		player *Player
	}
	type testCase struct {
		name  string
		state State
		args  args
		want  []Pai
	}
	tests := []testCase{}

	{
		state := getDefaultStateForTest()
		player := state.players[0]

		tests = append(tests, testCase{
			name:  "",
			state: state,
			args:  args{player: &player},
			want:  nil,
		})
	}
	{
		state := getDefaultStateForTest()
		pais, _ := StrToPais("1m 7z")
		state.players[3].ho = slices.Concat(state.players[3].ho, pais)
		player := state.players[3]

		tests = append(tests, testCase{
			name:  "",
			state: state,
			args:  args{player: &player},
			want:  pais,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.Anpais(tt.args.player); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.Anpais() = %v, want %v", got, tt.want)
			}
		})
	}
}
