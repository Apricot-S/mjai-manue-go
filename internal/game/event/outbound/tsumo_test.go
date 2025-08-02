package outbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestNewTsumo(t *testing.T) {
	type args struct {
		actor int
		pai   base.Pai
		log   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Tsumo
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor: 0,
				pai:   *mustPai("?"),
				log:   "",
			},
			want: &Tsumo{
				action: action{
					Actor: 0,
					Log:   "",
				},
				Pai: *mustPai("?"),
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor: 3,
				pai:   *mustPai("5sr"),
				log:   "test",
			},
			want: &Tsumo{
				action: action{
					Actor: 3,
					Log:   "test",
				},
				Pai: *mustPai("5sr"),
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor: -1,
				pai:   *mustPai("?"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor: 4,
				pai:   *mustPai("?"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTsumo(tt.args.actor, tt.args.pai, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTsumo() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTsumo() = %v, want %v", got, tt.want)
			}
		})
	}
}
