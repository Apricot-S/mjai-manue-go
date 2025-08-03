package outbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestNewAnkan(t *testing.T) {
	type args struct {
		actor    int
		target   int
		consumed [4]base.Pai
		log      string
	}
	tests := []struct {
		name    string
		args    args
		want    *Ankan
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:    1,
				target:   0,
				consumed: [4]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "",
			},
			want: &Ankan{
				action: action{
					Actor: 1,
					Log:   "",
				},
				Consumed: [4]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:    1,
				consumed: [4]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "test",
			},
			want: &Ankan{
				action: action{
					Actor: 1,
					Log:   "test",
				},
				Consumed: [4]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
			},
			wantErr: false,
		},
		{
			name: "invalid consumed tiles",
			args: args{
				actor:    2,
				target:   3,
				consumed: [4]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s"), *mustPai("7s")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown consumed tiles",
			args: args{
				actor:    2,
				target:   3,
				consumed: [4]base.Pai{*mustPai("?"), *mustPai("?"), *mustPai("?"), *mustPai("?")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "red tile in consumed",
			args: args{
				actor:    2,
				target:   3,
				consumed: [4]base.Pai{*mustPai("5s"), *mustPai("5s"), *mustPai("5s"), *mustPai("5sr")},
				log:      "",
			},
			want: &Ankan{
				action: action{
					Actor: 2,
					Log:   "",
				},
				Consumed: [4]base.Pai{*mustPai("5s"), *mustPai("5s"), *mustPai("5s"), *mustPai("5sr")},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAnkan(tt.args.actor, tt.args.consumed, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnkan() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnkan() = %v, want %v", got, tt.want)
			}
		})
	}
}
