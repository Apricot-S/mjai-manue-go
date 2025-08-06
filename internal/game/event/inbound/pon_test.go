package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestNewPon(t *testing.T) {
	type args struct {
		actor    int
		target   int
		taken    base.Pai
		consumed [2]base.Pai
	}
	tests := []struct {
		name    string
		args    args
		want    *Pon
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				actor:    1,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [2]base.Pai(mustPais("6s", "6s")),
			},
			want: &Pon{
				Actor:    1,
				Target:   0,
				Taken:    *mustPai("6s"),
				Consumed: [2]base.Pai(mustPais("6s", "6s")),
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
				target:   0,
				taken:    *mustPai("6s"),
				consumed: [2]base.Pai(mustPais("6s", "6s")),
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
				consumed: [2]base.Pai(mustPais("6s", "6s")),
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
				consumed: [2]base.Pai(mustPais("6s", "6s")),
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
				consumed: [2]base.Pai(mustPais("6s", "6s")),
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
				consumed: [2]base.Pai(mustPais("6s", "6s")),
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
				consumed: [2]base.Pai(mustPais("?", "?")),
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
				consumed: [2]base.Pai(mustPais("7s", "7s")),
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
				consumed: [2]base.Pai(mustPais("6s", "7s")),
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
				consumed: [2]base.Pai(mustPais("5s", "5s")),
			},
			want: &Pon{
				Actor:    2,
				Target:   3,
				Taken:    *mustPai("5sr"),
				Consumed: [2]base.Pai(mustPais("5s", "5s")),
			},
			wantErr: false,
		},
		{
			name: "red tile in consumed",
			args: args{
				actor:    2,
				target:   3,
				taken:    *mustPai("5s"),
				consumed: [2]base.Pai(mustPais("5s", "5sr")),
			},
			want: &Pon{
				Actor:    2,
				Target:   3,
				Taken:    *mustPai("5s"),
				Consumed: [2]base.Pai(mustPais("5s", "5sr")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPon(tt.args.actor, tt.args.target, tt.args.taken, tt.args.consumed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPon() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPon() = %v, want %v", got, tt.want)
			}
		})
	}
}
