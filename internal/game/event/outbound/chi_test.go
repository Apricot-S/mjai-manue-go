package outbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestNewChi(t *testing.T) {
	type args struct {
		actor    int
		target   int
		taken    base.Pai
		consumed [2]base.Pai
		log      string
	}
	tests := []struct {
		name    string
		args    args
		want    *Chi
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:    1,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [2]base.Pai(mustPais("5sr", "7s")),
				log:      "",
			},
			want: &Chi{
				action: action{
					Actor: 1,
					Log:   "",
				},
				Target:   0,
				Taken:    *mustPai("6s"),
				Consumed: [2]base.Pai(mustPais("5sr", "7s")),
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:    3,
				target:   2,
				taken:    *mustPai("5sr"),
				consumed: [2]base.Pai(mustPais("4s", "6s")),
				log:      "test",
			},
			want: &Chi{
				action: action{
					Actor: 3,
					Log:   "test",
				},
				Target:   2,
				Taken:    *mustPai("5sr"),
				Consumed: [2]base.Pai(mustPais("4s", "6s")),
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [2]base.Pai(mustPais("5sr", "7s")),
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
				consumed: [2]base.Pai(mustPais("5sr", "7s")),
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
				consumed: [2]base.Pai(mustPais("5sr", "7s")),
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
				consumed: [2]base.Pai(mustPais("5sr", "7s")),
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "target is not kamicha",
			args: args{
				actor:    0,
				target:   1,
				taken:    *mustPai("6s"),
				consumed: [2]base.Pai(mustPais("5sr", "7s")),
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown is not allowed",
			args: args{
				actor:    0,
				target:   3,
				taken:    *mustPai("?"),
				consumed: [2]base.Pai(mustPais("?", "?")),
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "tsupai is not allowed",
			args: args{
				actor:    0,
				target:   3,
				taken:    *mustPai("E"),
				consumed: [2]base.Pai(mustPais("S", "W")),
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "different color is not allowed: taken is different from consumed",
			args: args{
				actor:    0,
				target:   3,
				taken:    *mustPai("6s"),
				consumed: [2]base.Pai(mustPais("5p", "7p")),
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "different color is not allowed: consumed contains different color",
			args: args{
				actor:    0,
				target:   3,
				taken:    *mustPai("6p"),
				consumed: [2]base.Pai(mustPais("5p", "7s")),
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid sequence: taken is not in sequence with consumed",
			args: args{
				actor:    0,
				target:   3,
				taken:    *mustPai("6p"),
				consumed: [2]base.Pai(mustPais("5p", "8p")),
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChi(tt.args.actor, tt.args.target, tt.args.taken, tt.args.consumed, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChi() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChi() = %v, want %v", got, tt.want)
			}
		})
	}
}
