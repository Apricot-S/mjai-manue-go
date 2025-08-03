package outbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestNewDaiminkan(t *testing.T) {
	type args struct {
		actor    int
		target   int
		taken    base.Pai
		consumed [3]base.Pai
		log      string
	}
	tests := []struct {
		name    string
		args    args
		want    *Daiminkan
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:    1,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "",
			},
			want: &Daiminkan{
				action: action{
					Actor: 1,
					Log:   "",
				},
				Target:   0,
				Taken:    *mustPai("6s"),
				Consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:    1,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "test",
			},
			want: &Daiminkan{
				action: action{
					Actor: 1,
					Log:   "test",
				},
				Target:   0,
				Taken:    *mustPai("6s"),
				Consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:    4,
				target:   3,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: args{
				actor:    0,
				target:   -1,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: args{
				actor:    0,
				target:   4,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "actor and target are the same",
			args: args{
				actor:    0,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("6s")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown tile",
			args: args{
				actor:    0,
				target:   1,
				taken:    *mustPai("?"),
				consumed: [3]base.Pai{*mustPai("?"), *mustPai("?"), *mustPai("?")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid taken tile",
			args: args{
				actor:    0,
				target:   1,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("7s"), *mustPai("7s"), *mustPai("7s")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid consumed tiles",
			args: args{
				actor:    2,
				target:   3,
				taken:    *mustPai("6s"),
				consumed: [3]base.Pai{*mustPai("6s"), *mustPai("6s"), *mustPai("7s")},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "red tile in taken",
			args: args{
				actor:    2,
				target:   3,
				taken:    *mustPai("5sr"),
				consumed: [3]base.Pai{*mustPai("5s"), *mustPai("5s"), *mustPai("5s")},
				log:      "",
			},
			want: &Daiminkan{
				action: action{
					Actor: 2,
					Log:   "",
				},
				Target:   3,
				Taken:    *mustPai("5sr"),
				Consumed: [3]base.Pai{*mustPai("5s"), *mustPai("5s"), *mustPai("5s")},
			},
			wantErr: false,
		},
		{
			name: "red tile in consumed",
			args: args{
				actor:    2,
				target:   3,
				taken:    *mustPai("5s"),
				consumed: [3]base.Pai{*mustPai("5s"), *mustPai("5s"), *mustPai("5sr")},
				log:      "",
			},
			want: &Daiminkan{
				action: action{
					Actor: 2,
					Log:   "",
				},
				Target:   3,
				Taken:    *mustPai("5s"),
				Consumed: [3]base.Pai{*mustPai("5s"), *mustPai("5s"), *mustPai("5sr")},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDaiminkan(tt.args.actor, tt.args.target, tt.args.taken, tt.args.consumed, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDaiminkan() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDaiminkan() = %v, want %v", got, tt.want)
			}
		})
	}
}
