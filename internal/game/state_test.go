package game

import (
	"reflect"
	"slices"
	"testing"
)

func Test_GetPlayerDistance(t *testing.T) {
	type args struct {
		p1 *Player
		p2 *Player
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "p1 and p2 are the same, p1 is 0",
			args: args{p1: &Player{id: 0}, p2: &Player{id: 0}},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 1",
			args: args{p1: &Player{id: 1}, p2: &Player{id: 1}},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 2",
			args: args{p1: &Player{id: 2}, p2: &Player{id: 2}},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 3",
			args: args{p1: &Player{id: 3}, p2: &Player{id: 3}},
			want: 0,
		},
		{
			name: "p2 is shimocha of p1, p1 is 0",
			args: args{p1: &Player{id: 0}, p2: &Player{id: 1}},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 1",
			args: args{p1: &Player{id: 1}, p2: &Player{id: 2}},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 2",
			args: args{p1: &Player{id: 2}, p2: &Player{id: 3}},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 3",
			args: args{p1: &Player{id: 3}, p2: &Player{id: 0}},
			want: 3,
		},
		{
			name: "p2 is toimen of p1, p1 is 0",
			args: args{p1: &Player{id: 0}, p2: &Player{id: 2}},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 1",
			args: args{p1: &Player{id: 1}, p2: &Player{id: 3}},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 2",
			args: args{p1: &Player{id: 2}, p2: &Player{id: 0}},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 3",
			args: args{p1: &Player{id: 3}, p2: &Player{id: 1}},
			want: 2,
		},
		{
			name: "p2 is kamicha of p1, p1 is 0",
			args: args{p1: &Player{id: 0}, p2: &Player{id: 3}},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 1",
			args: args{p1: &Player{id: 1}, p2: &Player{id: 0}},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 2",
			args: args{p1: &Player{id: 2}, p2: &Player{id: 1}},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 3",
			args: args{p1: &Player{id: 3}, p2: &Player{id: 2}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPlayerDistance(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("GetPlayerDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNextKyoku(t *testing.T) {
	east, _ := NewPaiWithName("E")
	south, _ := NewPaiWithName("S")
	west, _ := NewPaiWithName("W")
	north, _ := NewPaiWithName("N")

	type args struct {
		bakaze   *Pai
		kyokuNum int
	}
	tests := []struct {
		name  string
		args  args
		want  *Pai
		want1 int
	}{
		{
			name:  "E1 -> E2",
			args:  args{bakaze: east, kyokuNum: 1},
			want:  east,
			want1: 2,
		},
		{
			name:  "E4 -> S1",
			args:  args{bakaze: east, kyokuNum: 4},
			want:  south,
			want1: 1,
		},
		{
			name:  "S2 -> S3",
			args:  args{bakaze: south, kyokuNum: 2},
			want:  south,
			want1: 3,
		},
		{
			name:  "S4 -> W1",
			args:  args{bakaze: south, kyokuNum: 4},
			want:  west,
			want1: 1,
		},
		{
			name:  "W3 -> W4",
			args:  args{bakaze: west, kyokuNum: 3},
			want:  west,
			want1: 4,
		},
		{
			name:  "W4 -> N1",
			args:  args{bakaze: west, kyokuNum: 4},
			want:  north,
			want1: 1,
		},
		{
			name:  "N4 -> E1",
			args:  args{bakaze: north, kyokuNum: 4},
			want:  east,
			want1: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getNextKyoku(tt.args.bakaze, tt.args.kyokuNum)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNextKyoku() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getNextKyoku() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func getDefaultStateForTest() *StateImpl {
	east, _ := NewPaiWithName("E")
	players := [numPlayers]Player{}

	for i := range numPlayers {
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
			reachState:        None,
			reachHoIndex:      -1,
			reachSutehaiIndex: -1,
			score:             initScore,
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
		doraMarkers: make([]Pai, 0, maxNumDoraMarkers),
		numPipais:   numInitPipais,

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
