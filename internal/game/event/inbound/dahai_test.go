package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/testutil"
)

func TestNewDahai(t *testing.T) {
	type args struct {
		actor     int
		pai       base.Pai
		tsumogiri bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Dahai
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				actor:     0,
				pai:       *testutil.MustPai("7p"),
				tsumogiri: false,
			},
			want: &Dahai{
				Actor:     0,
				Pai:       *testutil.MustPai("7p"),
				Tsumogiri: false,
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:     -1,
				pai:       *testutil.MustPai("7p"),
				tsumogiri: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:     4,
				pai:       *testutil.MustPai("7p"),
				tsumogiri: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown pai",
			args: args{
				actor:     2,
				pai:       *testutil.MustPai("?"),
				tsumogiri: true,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDahai(tt.args.actor, tt.args.pai, tt.args.tsumogiri)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDahai() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDahai() = %v, want %v", got, tt.want)
			}
		})
	}
}
