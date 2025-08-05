package mjai

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewEndGame(t *testing.T) {
	tests := []struct {
		name    string
		want    *EndGame
		wantErr bool
	}{
		{
			name: "test NewEndGame()",
			want: &EndGame{
				Message: Message{Type: TypeEndGame},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEndGame()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndGame_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *EndGame
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			args: &EndGame{
				Message: Message{Type: TypeEndGame},
			},
			want:    `{"type":"end_game"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &EndGame{
				Message: Message{Type: ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &EndGame{
				Message: Message{Type: TypeHello},
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

func TestEndGame_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    EndGame
		wantErr bool
	}{
		{
			name: "valid",
			args: `{"type":"end_game"}`,
			want: EndGame{
				Message: Message{Type: TypeEndGame},
			},
			wantErr: false,
		},
		{
			name: "with metadata",
			args: `{
				"type":"end_game",
				"metadata": {
					"foo": "bar"
				}
			}`,
			want: EndGame{
				Message: Message{Type: TypeEndGame},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: EndGame{
				Message: Message{Type: ""},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello"}`,
			want: EndGame{
				Message: Message{Type: TypeHello},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got EndGame
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
