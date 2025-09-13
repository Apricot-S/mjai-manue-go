package outbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/testutil"
)

func TestNewHora(t *testing.T) {
	type args struct {
		actor  int
		target int
		pai    base.Pai
		log    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Hora
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:  1,
				target: 0,
				pai:    *testutil.MustPai("6s"),
				log:    "",
			},
			want: &Hora{
				action: action{
					Actor: 1,
					Log:   "",
				},
				Target: 0,
				Pai:    *testutil.MustPai("6s"),
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:  3,
				target: 2,
				pai:    *testutil.MustPai("5sr"),
				log:    "test",
			},
			want: &Hora{
				action: action{
					Actor: 3,
					Log:   "test",
				},
				Target: 2,
				Pai:    *testutil.MustPai("5sr"),
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:  -1,
				target: 0,
				pai:    *testutil.MustPai("6s"),
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:  4,
				target: 3,
				pai:    *testutil.MustPai("6s"),
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: args{
				actor:  0,
				target: -1,
				pai:    *testutil.MustPai("6s"),
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: args{
				actor:  0,
				target: 4,
				pai:    *testutil.MustPai("6s"),
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown pai",
			args: args{
				actor:  0,
				target: 3,
				pai:    *testutil.MustPai("?"),
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHora(tt.args.actor, tt.args.target, tt.args.pai, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHora() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHora() = %v, want %v", got, tt.want)
			}
		})
	}
}
