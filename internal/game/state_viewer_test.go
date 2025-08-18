package game

import (
	"reflect"
	"slices"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func getDefaultStateForTest() *StateImpl {
	east, _ := base.NewPaiWithName("E")
	players := [NumPlayers]base.Player{}

	for i := range NumPlayers {
		tehais := make([]base.Pai, 13)
		for j := range 13 {
			tehais[j] = *base.Unknown
		}

		players[i] = *base.NewPlayerForTest(
			i,
			"",
			tehais,
			make([]base.Furo, 0, base.MaxNumFuro),
			make([]base.Pai, 0),
			make([]base.Pai, 0),
			make([]base.Pai, 0),
			base.NotReach,
			-1,
			-1,
			InitScore,
			false,
			true,
		)
	}

	return &StateImpl{
		players:     players,
		bakaze:      *east,
		kyokuNum:    1,
		honba:       0,
		oya:         &players[0],
		chicha:      &players[0],
		doraMarkers: make([]base.Pai, 0, MaxNumDoraMarkers),
		numPipais:   NumInitPipais,

		prevDahaiActor: -1,
		prevDahaiPai:   nil,
		currentEvent:   nil,
	}
}

func TestState_Anpais(t *testing.T) {
	type args struct {
		player *base.Player
	}
	type testCase struct {
		name  string
		state State
		args  args
		want  []base.Pai
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
		pais, _ := base.StrToPais("1m 7z")
		state.players[3] = *base.NewPlayerForTest(
			state.players[3].ID(),
			state.players[3].Name(),
			state.players[3].Tehais(),
			state.players[3].Furos(),
			slices.Concat(state.players[3].Ho(), pais),
			state.players[3].Sutehais(),
			state.players[3].ExtraAnpais(),
			state.players[3].ReachState(),
			state.players[3].ReachHoIndex(),
			state.players[3].ReachSutehaiIndex(),
			state.players[3].Score(),
			state.players[3].CanDahai(),
			state.players[3].IsMenzen(),
		)
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
