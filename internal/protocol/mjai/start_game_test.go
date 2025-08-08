package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/go-json-experiment/json"
)

func TestNewStartGame(t *testing.T) {
	type args struct {
		id    int
		names []string
	}
	tests := []struct {
		name    string
		args    args
		want    *StartGame
		wantErr bool
	}{
		{
			name: "without names",
			args: args{
				id:    0,
				names: nil,
			},
			want: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      0,
				Names:   nil,
			},
			wantErr: false,
		},
		{
			name: "empty string names",
			args: args{
				id:    0,
				names: []string{"", "", "", ""},
			},
			want: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      0,
				Names:   []string{"", "", "", ""},
			},
			wantErr: false,
		},
		{
			name: "with names",
			args: args{
				id:    3,
				names: []string{"shanten0", "shanten1", "shanten2", "shanten3"},
			},
			want: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      3,
				Names:   []string{"shanten0", "shanten1", "shanten2", "shanten3"},
			},
			wantErr: false,
		},
		{
			name: "invalid id min",
			args: args{
				id:    -1,
				names: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid id max",
			args: args{
				id:    4,
				names: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid names empty",
			args: args{
				id:    0,
				names: []string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid names 3",
			args: args{
				id:    0,
				names: []string{"shanten0", "shanten1", "shanten2"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid names 5",
			args: args{
				id:    0,
				names: []string{"shanten0", "shanten1", "shanten2", "shanten3", "shanten4"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStartGame(tt.args.id, tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStartGame() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStartGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartGame_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *StartGame
		want    string
		wantErr bool
	}{
		{
			name: "without names",
			args: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      0,
				Names:   nil,
			},
			want:    `{"type":"start_game","id":0}`,
			wantErr: false,
		},
		{
			name: "with names",
			args: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      3,
				Names:   []string{"shanten0", "shanten1", "shanten2", "shanten3"},
			},
			want:    `{"type":"start_game","id":3,"names":["shanten0","shanten1","shanten2","shanten3"]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &StartGame{
				Message: Message{Type: ""},
				ID:      0,
				Names:   nil,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &StartGame{
				Message: Message{Type: TypeNone},
				ID:      3,
				Names:   nil,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid id min",
			args: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      -1,
				Names:   nil,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid id max",
			args: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      4,
				Names:   nil,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid names empty",
			args: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      1,
				Names:   []string{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid names 3",
			args: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      1,
				Names:   []string{"shanten0", "shanten1", "shanten2"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid names 5",
			args: &StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      2,
				Names:   []string{"shanten0", "shanten1", "shanten2", "shanten3", "shanten4"},
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal error = %v, want %v", err, tt.wantErr)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestStartGame_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    StartGame
		wantErr bool
	}{
		{
			name: "without id",
			args: `{"type":"start_game"}`,
			want: StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      0,
				Names:   nil,
			},
			wantErr: false,
		},
		{
			name: "without names",
			args: `{"type":"start_game","id":2}`,
			want: StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      2,
				Names:   nil,
			},
			wantErr: false,
		},
		{
			name: "with names",
			args: `{"type":"start_game","id":0,"names":["shanten0","shanten1","shanten2","shanten3"]}`,
			want: StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      0,
				Names:   []string{"shanten0", "shanten1", "shanten2", "shanten3"},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: StartGame{
				Message: Message{Type: ""},
				ID:      0,
				Names:   nil,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none"}`,
			want: StartGame{
				Message: Message{Type: TypeNone},
				ID:      0,
				Names:   nil,
			},
			wantErr: true,
		},
		{
			name: "invalid id min",
			args: `{"type":"start_game","id":-1}`,
			want: StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      -1,
				Names:   nil,
			},
			wantErr: true,
		},
		{
			name: "invalid id max",
			args: `{"type":"start_game","id":4}`,
			want: StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      4,
				Names:   nil,
			},
			wantErr: true,
		},
		{
			name: "invalid names empty",
			args: `{"type":"start_game","id":0,"names":[]}`,
			want: StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      0,
				Names:   []string{},
			},
			wantErr: true,
		},
		{
			name: "invalid names 3",
			args: `{"type":"start_game","id":0,"names":["shanten0","shanten1","shanten2"]}`,
			want: StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      0,
				Names:   []string{"shanten0", "shanten1", "shanten2"},
			},
			wantErr: true,
		},
		{
			name: "invalid names 5",
			args: `{"type":"start_game","id":0,"names":["shanten0","shanten1","shanten2","shanten3","shanten4"]}`,
			want: StartGame{
				Message: Message{Type: TypeStartGame},
				ID:      0,
				Names:   []string{"shanten0", "shanten1", "shanten2", "shanten3", "shanten4"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got StartGame
			err := json.Unmarshal([]byte(tt.args), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal error = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartGame_ToEvent(t *testing.T) {
	tests := []struct {
		name    string
		args    *StartGame
		want    *inbound.StartGame
		wantErr bool
	}{
		{
			name: "without names",
			args: &StartGame{
				Message: Message{TypeStartGame},
				ID:      3,
				Names:   nil,
			},
			want: &inbound.StartGame{
				ID:    3,
				Names: [4]string{"", "", "", ""},
			},
			wantErr: false,
		},
		{
			name: "with names",
			args: &StartGame{
				Message: Message{TypeStartGame},
				ID:      3,
				Names:   []string{"a", "b", "c", "d"},
			},
			want: &inbound.StartGame{
				ID:    3,
				Names: [4]string{"a", "b", "c", "d"},
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: &StartGame{
				Message: Message{TypeStartGame},
				ID:      4,
				Names:   nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.ToEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("StartGame.ToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartGame.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
