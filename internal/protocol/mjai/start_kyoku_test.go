package mjai

import (
	"encoding/json/v2"
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

func TestNewStartKyoku(t *testing.T) {
	type args struct {
		bakaze     string
		kyoku      int
		honba      int
		kyotaku    int
		oya        int
		doraMarker string
		scores     []int
		tehais     [4][13]string
	}
	tests := []struct {
		name    string
		args    args
		want    *StartKyoku
		wantErr bool
	}{
		{
			name: "without scores",
			args: args{
				bakaze:     "E",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: args{
				bakaze:     "S",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     []int{25000, 25000, 25000, 25000},
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "S",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     []int{25000, 25000, 25000, 25000},
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid bakaze",
			args: args{
				bakaze:     "P",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid kyoku min",
			args: args{
				bakaze:     "E",
				kyoku:      0,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid kyoku max",
			args: args{
				bakaze:     "E",
				kyoku:      5,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid honba",
			args: args{
				bakaze:     "E",
				kyoku:      1,
				honba:      -1,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid kyotaku",
			args: args{
				bakaze:     "E",
				kyoku:      1,
				honba:      0,
				kyotaku:    -1,
				oya:        0,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid oya min",
			args: args{
				bakaze:     "E",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        -1,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid oya max",
			args: args{
				bakaze:     "E",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        4,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid dora marker",
			args: args{
				bakaze:     "E",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7sr",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid tehais",
			args: args{
				bakaze:     "E",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     nil,
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "1z"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: args{
				bakaze:     "S",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     []int{},
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: args{
				bakaze:     "S",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     []int{25000, 25000, 25000},
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: args{
				bakaze:     "S",
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: "7s",
				scores:     []int{25000, 25000, 25000, 25000, 25000},
				tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStartKyoku(tt.args.bakaze, tt.args.kyoku, tt.args.honba, tt.args.kyotaku, tt.args.oya, tt.args.doraMarker, tt.args.scores, tt.args.tehais)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStartKyoku() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStartKyoku() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartKyoku_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *StartKyoku
		want    string
		wantErr bool
	}{
		{
			name: "without scores",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "W",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    `{"type":"start_kyoku","bakaze":"W","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			wantErr: false,
		},
		{
			name: "with scores",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "N",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "2s",
				Scores:     []int{25000, 25000, 25000, 25000},
				Tehais: [4][13]string{
					{"9p", "7m", "9s", "9s", "8m", "5s", "2p", "W", "C", "5s", "N", "5mr", "F"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    `{"type":"start_kyoku","bakaze":"N","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"2s","scores":[25000,25000,25000,25000],"tehais":[["9p","7m","9s","9s","8m","5s","2p","W","C","5s","N","5mr","F"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &StartKyoku{
				Message:    Message{Type: ""},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &StartKyoku{
				Message:    Message{Type: TypeHello},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid bakaze",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "F",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid kyokyu min",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      0,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid kyokyu max",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      5,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid honba",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      -1,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid kyotaku",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    -1,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid oya min",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        -1,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid oya max",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        4,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid dora marker",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7sr",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid tehais",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3mr", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "N",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "2s",
				Scores:     []int{},
				Tehais: [4][13]string{
					{"9p", "7m", "9s", "9s", "8m", "5s", "2p", "W", "C", "5s", "N", "5mr", "F"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "W",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "2s",
				Scores:     []int{25000, 25000, 25000},
				Tehais: [4][13]string{
					{"9p", "7m", "9s", "9s", "8m", "5s", "2p", "W", "C", "5s", "N", "5mr", "F"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "S",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "2s",
				Scores:     []int{25000, 25000, 25000, 25000, 25000},
				Tehais: [4][13]string{
					{"9p", "7m", "9s", "9s", "8m", "5s", "2p", "W", "C", "5s", "N", "5mr", "F"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
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

func TestStartKyoku_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    StartKyoku
		wantErr bool
	}{
		{
			name: "without scores",
			args: `{"type":"start_kyoku","bakaze":"W","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "W",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: `{"type":"start_kyoku","bakaze":"E","dora_marker":"2s","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"scores":[25000,25000,25000,25000],"tehais":[["9p","7m","9s","9s","8m","5s","2p","W","C","5s","N","5mr","F"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "2s",
				Scores:     []int{25000, 25000, 25000, 25000},
				Tehais: [4][13]string{
					{"9p", "7m", "9s", "9s", "8m", "5s", "2p", "W", "C", "5s", "N", "5mr", "F"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","bakaze":"W","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: ""},
				Bakaze:     "W",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello","bakaze":"W","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeHello},
				Bakaze:     "W",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid bakaze",
			args: `{"type":"start_kyoku","bakaze":"?","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "?",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid kyokyu min",
			args: `{"type":"start_kyoku","bakaze":"E","kyoku":0,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      0,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid kyokyu max",
			args: `{"type":"start_kyoku","bakaze":"E","kyoku":5,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      5,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid honba",
			args: `{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":-1,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      -1,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid kyotaku",
			args: `{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":-1,"oya":0,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    -1,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid oya min",
			args: `{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":-1,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        -1,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid oya max",
			args: `{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":4,"dora_marker":"7s","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        4,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid dora marker",
			args: `{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7sr","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7sr",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid tehais",
			args: `{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","tehais":[["1z","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"1z", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name:    "invalid tehais length",
			args:    `{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7sr","tehais":[["?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want:    StartKyoku{},
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: `{"type":"start_kyoku","bakaze":"W","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","scores":[],"tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "W",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     []int{},
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: `{"type":"start_kyoku","bakaze":"W","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","scores":[25000,25000,25000],"tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "W",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     []int{25000, 25000, 25000},
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: `{"type":"start_kyoku","bakaze":"W","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"7s","scores":[25000,25000,25000,25000,25000],"tehais":[["?","?","?","?","?","?","?","?","?","?","?","?","?"],["3m","4m","3p","5pr","7p","9p","4s","4s","5sr","7s","7s","W","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"]]}`,
			want: StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "W",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     []int{25000, 25000, 25000, 25000, 25000},
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got StartKyoku
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

func TestStartKyoku_ToEvent(t *testing.T) {
	tests := []struct {
		name    string
		args    *StartKyoku
		want    *inbound.StartKyoku
		wantErr bool
	}{
		{
			name: "without scores",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "E",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     nil,
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want: &inbound.StartKyoku{
				Bakaze:     *mustPai("E"),
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: *mustPai("7s"),
				Scores:     nil,
				Tehais: [4][13]base.Pai{
					[13]base.Pai(mustPais("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?")),
					[13]base.Pai(mustPais("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N")),
					[13]base.Pai(mustPais("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?")),
					[13]base.Pai(mustPais("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?")),
				},
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "S",
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     []int{25000, 25000, 25000, 25000},
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want: &inbound.StartKyoku{
				Bakaze:     *mustPai("S"),
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: *mustPai("7s"),
				Scores:     &[4]int{25000, 25000, 25000, 25000},
				Tehais: [4][13]base.Pai{
					[13]base.Pai(mustPais("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?")),
					[13]base.Pai(mustPais("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N")),
					[13]base.Pai(mustPais("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?")),
					[13]base.Pai(mustPais("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?")),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: &StartKyoku{
				Message:    Message{Type: TypeStartKyoku},
				Bakaze:     "W",
				Kyoku:      0,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: "7s",
				Scores:     []int{25000, 25000, 25000, 25000},
				Tehais: [4][13]string{
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
					{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.ToEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("StartKyoku.ToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartKyoku.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
