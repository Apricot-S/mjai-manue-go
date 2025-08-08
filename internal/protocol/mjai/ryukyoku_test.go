package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/go-json-experiment/json"
)

func TestNewRyukyoku(t *testing.T) {
	type args struct {
		scores []int
	}
	tests := []struct {
		name    string
		args    args
		want    *Ryukyoku
		wantErr bool
	}{
		{
			name: "without scores",
			args: args{
				scores: nil,
			},
			want: &Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  nil,
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: args{
				scores: []int{28000, 23000, 24000, 24000},
			},
			want: &Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{28000, 23000, 24000, 24000},
			},
			wantErr: false,
		},
		{
			name: "invalid scores empty",
			args: args{
				scores: []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: args{
				scores: []int{28000, 23000, 24000},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: args{
				scores: []int{28000, 23000, 24000, 24000, 24000},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRyukyoku(tt.args.scores)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRyukyoku() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRyukyoku() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRyukyoku_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Ryukyoku
		want    string
		wantErr bool
	}{
		{
			name: "without scores",
			args: &Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  nil,
			},
			want:    `{"type":"ryukyoku"}`,
			wantErr: false,
		},
		{
			name: "with scores",
			args: &Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{28000, 23000, 24000, 24000},
			},
			want:    `{"type":"ryukyoku","scores":[28000,23000,24000,24000]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Ryukyoku{
				Message: Message{""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Ryukyoku{
				Message: Message{TypeHello},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: &Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: &Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{28000, 23000, 24000},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: &Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{28000, 23000, 24000, 24000, 24000},
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

func TestRyukyoku_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Ryukyoku
		wantErr bool
	}{
		{
			name: "without scores",
			args: `{"type":"ryukyoku"}`,
			want: Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  nil,
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: `{"type":"ryukyoku","reason":"fanpai","tehais":[["5m","5m","5mr","3s","3s","N","N"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?","?","?","?","?","?","?","?","?","?"],["?","?","?","?"]],"tenpais":[true,false,false,false],"deltas":[3000,-1000,-1000,-1000],"scores":[28000,24000,24000,24000]}`,
			want: Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{28000, 24000, 24000, 24000},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: Ryukyoku{
				Message: Message{""},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello"}`,
			want: Ryukyoku{
				Message: Message{TypeHello},
			},
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: `{"type":"ryukyoku","scores":[]}`,
			want: Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{},
			},
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: `{"type":"ryukyoku","scores":[28000,23000,24000]}`,
			want: Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{28000, 23000, 24000},
			},
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: `{"type":"ryukyoku","scores":[28000,23000,24000,24000,24000]}`,
			want: Ryukyoku{
				Message: Message{TypeRyukyoku},
				Scores:  []int{28000, 23000, 24000, 24000, 24000},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Ryukyoku
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

func TestRyukyoku_ToEvent(t *testing.T) {
	tests := []struct {
		name    string
		args    *Ryukyoku
		want    *inbound.Ryukyoku
		wantErr bool
	}{
		{
			name: "without scores",
			args: &Ryukyoku{
				Message: Message{Type: TypeRyukyoku},
				Scores:  nil,
			},
			want: &inbound.Ryukyoku{
				Scores: nil,
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: &Ryukyoku{
				Message: Message{Type: TypeRyukyoku},
				Scores:  []int{26000, 24000, 23000, 24000},
			},
			want: &inbound.Ryukyoku{
				Scores: &[4]int{26000, 24000, 23000, 24000},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.ToEvent()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ryukyoku.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
