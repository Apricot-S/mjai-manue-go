package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/testutil"
)

func TestNewTsumo(t *testing.T) {
	type args struct {
		actor int
		pai   base.Pai
	}
	tests := []struct {
		name    string
		args    args
		want    *Tsumo
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				actor: 0,
				pai:   *testutil.MustPai("?"),
			},
			want: &Tsumo{
				Actor: 0,
				Pai:   *testutil.MustPai("?"),
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor: -1,
				pai:   *testutil.MustPai("?"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor: 4,
				pai:   *testutil.MustPai("?"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTsumo(tt.args.actor, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTsumo() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTsumo() = %v, want %v", got, tt.want)
			}
		})
	}
}
