package inbound

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
	}
	tests := []struct {
		name    string
		args    args
		want    *Chi
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				actor:    1,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [2]base.Pai(mustPais("5sr", "7s")),
			},
			want: &Chi{
				Actor:    1,
				Target:   0,
				Taken:    *mustPai("6s"),
				Consumed: [2]base.Pai(mustPais("5sr", "7s")),
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
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChi(tt.args.actor, tt.args.target, tt.args.taken, tt.args.consumed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChi() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChi() = %v, want %v", got, tt.want)
			}
		})
	}
}
