package outbound

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
		log       string
	}
	tests := []struct {
		name    string
		args    args
		want    *Dahai
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:     0,
				pai:       *testutil.MustPai("1m"),
				tsumogiri: false,
				log:       "",
			},
			want: &Dahai{
				action: action{
					Actor: 0,
					Log:   "",
				},
				Pai:       *testutil.MustPai("1m"),
				Tsumogiri: false,
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:     3,
				pai:       *testutil.MustPai("5sr"),
				tsumogiri: true,
				log:       "test",
			},
			want: &Dahai{
				action: action{
					Actor: 3,
					Log:   "test",
				},
				Pai:       *testutil.MustPai("5sr"),
				Tsumogiri: true,
			},
			wantErr: false,
		},
		{
			name: "invalid pai unknown",
			args: args{
				actor:     -1,
				pai:       *testutil.MustPai("?"),
				tsumogiri: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:     -1,
				pai:       *testutil.MustPai("1m"),
				tsumogiri: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:     4,
				pai:       *testutil.MustPai("1m"),
				tsumogiri: true,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDahai(tt.args.actor, tt.args.pai, tt.args.tsumogiri, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDahai() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDahai() = %v, want %v", got, tt.want)
			}
		})
	}
}
