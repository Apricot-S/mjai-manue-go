package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/testutil"
)

func TestNewKakan(t *testing.T) {
	type args struct {
		actor    int
		target   int
		taken    base.Pai
		consumed [2]base.Pai
		added    base.Pai
	}
	tests := []struct {
		name    string
		args    args
		want    *Kakan
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				actor:    1,
				target:   0,
				taken:    *testutil.MustPai("6s"),
				consumed: [2]base.Pai(testutil.MustPais("6s", "6s")),
				added:    *testutil.MustPai("6s"),
			},
			want: &Kakan{
				Actor:    1,
				Target:   0,
				Taken:    *testutil.MustPai("6s"),
				Consumed: [2]base.Pai(testutil.MustPais("6s", "6s")),
				Added:    *testutil.MustPai("6s"),
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
				target:   0,
				taken:    *testutil.MustPai("6s"),
				consumed: [2]base.Pai(testutil.MustPais("6s", "6s")),
				added:    *testutil.MustPai("6s"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:    4,
				target:   3,
				taken:    *testutil.MustPai("6s"),
				consumed: [2]base.Pai(testutil.MustPais("6s", "6s")),
				added:    *testutil.MustPai("6s"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: args{
				actor:    0,
				target:   -1,
				taken:    *testutil.MustPai("6s"),
				consumed: [2]base.Pai(testutil.MustPais("6s", "6s")),
				added:    *testutil.MustPai("6s"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: args{
				actor:    0,
				target:   4,
				taken:    *testutil.MustPai("6s"),
				consumed: [2]base.Pai(testutil.MustPais("6s", "6s")),
				added:    *testutil.MustPai("6s"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "actor and target are the same",
			args: args{
				actor:    0,
				target:   0,
				taken:    *testutil.MustPai("6s"),
				consumed: [2]base.Pai(testutil.MustPais("6s", "6s")),
				added:    *testutil.MustPai("6s"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown tile",
			args: args{
				actor:    0,
				target:   1,
				taken:    *testutil.MustPai("?"),
				consumed: [2]base.Pai(testutil.MustPais("?", "?")),
				added:    *testutil.MustPai("?"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid taken tile",
			args: args{
				actor:    0,
				target:   1,
				taken:    *testutil.MustPai("6s"),
				consumed: [2]base.Pai(testutil.MustPais("7s", "7s")),
				added:    *testutil.MustPai("7s"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid consumed tiles",
			args: args{
				actor:    2,
				target:   3,
				taken:    *testutil.MustPai("6s"),
				consumed: [2]base.Pai(testutil.MustPais("6s", "7s")),
				added:    *testutil.MustPai("6s"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid added tile",
			args: args{
				actor:    2,
				target:   3,
				taken:    *testutil.MustPai("5s"),
				consumed: [2]base.Pai(testutil.MustPais("5s", "5s")),
				added:    *testutil.MustPai("6s"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "red tile in taken",
			args: args{
				actor:    2,
				target:   3,
				taken:    *testutil.MustPai("5sr"),
				consumed: [2]base.Pai(testutil.MustPais("5s", "5s")),
				added:    *testutil.MustPai("5s"),
			},
			want: &Kakan{
				Actor:    2,
				Target:   3,
				Taken:    *testutil.MustPai("5sr"),
				Consumed: [2]base.Pai(testutil.MustPais("5s", "5s")),
				Added:    *testutil.MustPai("5s"),
			},
			wantErr: false,
		},
		{
			name: "red tile in consumed",
			args: args{
				actor:    2,
				target:   3,
				taken:    *testutil.MustPai("5s"),
				consumed: [2]base.Pai(testutil.MustPais("5s", "5sr")),
				added:    *testutil.MustPai("5s"),
			},
			want: &Kakan{
				Actor:    2,
				Target:   3,
				Taken:    *testutil.MustPai("5s"),
				Consumed: [2]base.Pai(testutil.MustPais("5s", "5sr")),
				Added:    *testutil.MustPai("5s"),
			},
			wantErr: false,
		},
		{
			name: "red tile in added",
			args: args{
				actor:    2,
				target:   3,
				taken:    *testutil.MustPai("5s"),
				consumed: [2]base.Pai(testutil.MustPais("5s", "5s")),
				added:    *testutil.MustPai("5sr"),
			},
			want: &Kakan{
				Actor:    2,
				Target:   3,
				Taken:    *testutil.MustPai("5s"),
				Consumed: [2]base.Pai(testutil.MustPais("5s", "5s")),
				Added:    *testutil.MustPai("5sr"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKakan(tt.args.actor, tt.args.target, tt.args.taken, tt.args.consumed, tt.args.added)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKakan() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKakan() = %v, want %v", got, tt.want)
			}
		})
	}
}
