package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestNewHora(t *testing.T) {
	type args struct {
		actor      int
		target     int
		pai        *base.Pai
		horaPoints *int
		scores     *[4]int
	}
	tests := []struct {
		name    string
		args    args
		want    *Hora
		wantErr bool
	}{
		{
			name: "without hora_points",
			args: args{
				actor:      1,
				target:     0,
				pai:        mustPai("6s"),
				horaPoints: nil,
				scores:     nil,
			},
			want: &Hora{
				Actor:      1,
				Target:     0,
				Pai:        mustPai("6s"),
				HoraPoints: nil,
				Scores:     nil,
			},
			wantErr: false,
		},
		{
			name: "without scores",
			args: args{
				actor:      1,
				target:     0,
				pai:        mustPai("6s"),
				horaPoints: func() *int { p := 2600; return &p }(),
				scores:     nil,
			},
			want: &Hora{
				Actor:      1,
				Target:     0,
				Pai:        mustPai("6s"),
				HoraPoints: func() *int { p := 2600; return &p }(),
				Scores:     nil,
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				actor:      1,
				target:     0,
				pai:        mustPai("6s"),
				horaPoints: func() *int { p := 2600; return &p }(),
				scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			want: &Hora{
				Actor:      1,
				Target:     0,
				Pai:        mustPai("6s"),
				HoraPoints: func() *int { p := 2600; return &p }(),
				Scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:      -1,
				target:     0,
				pai:        mustPai("6s"),
				horaPoints: func() *int { p := 2600; return &p }(),
				scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:      4,
				target:     3,
				pai:        mustPai("6s"),
				horaPoints: func() *int { p := 2600; return &p }(),
				scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: args{
				actor:      0,
				target:     -1,
				pai:        mustPai("6s"),
				horaPoints: func() *int { p := 2600; return &p }(),
				scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: args{
				actor:      0,
				target:     4,
				pai:        mustPai("6s"),
				horaPoints: func() *int { p := 2600; return &p }(),
				scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil pai",
			args: args{
				actor:      0,
				target:     0,
				pai:        nil,
				horaPoints: func() *int { p := 0; return &p }(),
				scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			want: &Hora{
				Actor:      0,
				Target:     0,
				Pai:        nil,
				HoraPoints: func() *int { p := 0; return &p }(),
				Scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "unknown pai",
			args: args{
				actor:      0,
				target:     1,
				pai:        mustPai("?"),
				horaPoints: func() *int { p := 0; return &p }(),
				scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid hora_points",
			args: args{
				actor:      0,
				target:     2,
				pai:        mustPai("6s"),
				horaPoints: func() *int { p := -1; return &p }(),
				scores:     &[4]int{27500, 22300, 24300, 25900},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHora(tt.args.actor, tt.args.target, tt.args.pai, tt.args.horaPoints, tt.args.scores)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHora() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHora() = %v, want %v", got, tt.want)
			}
		})
	}
}
