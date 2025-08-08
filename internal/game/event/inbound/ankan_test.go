package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestNewAnkan(t *testing.T) {
	type args struct {
		actor    int
		consumed [4]base.Pai
	}
	tests := []struct {
		name    string
		args    args
		want    *Ankan
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				actor:    1,
				consumed: [4]base.Pai(mustPais("6s", "6s", "6s", "6s")),
			},
			want: &Ankan{
				Actor:    1,
				Consumed: [4]base.Pai(mustPais("6s", "6s", "6s", "6s")),
			},
			wantErr: false,
		},
		{
			name: "invalid consumed tiles",
			args: args{
				actor:    2,
				consumed: [4]base.Pai(mustPais("6s", "6s", "6s", "7s")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown consumed tiles",
			args: args{
				actor:    2,
				consumed: [4]base.Pai(mustPais("?", "?", "?", "?")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "red tile in consumed",
			args: args{
				actor:    2,
				consumed: [4]base.Pai(mustPais("5s", "5s", "5s", "5sr")),
			},
			want: &Ankan{
				Actor:    2,
				Consumed: [4]base.Pai(mustPais("5s", "5s", "5s", "5sr")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAnkan(tt.args.actor, tt.args.consumed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnkan() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnkan() = %v, want %v", got, tt.want)
			}
		})
	}
}
