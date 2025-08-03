package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestNewDora(t *testing.T) {
	type args struct {
		doraMarker base.Pai
	}
	tests := []struct {
		name string
		args args
		want *Dora
	}{
		{
			name: "dora",
			args: args{
				doraMarker: *mustPai("6s"),
			},
			want: &Dora{
				DoraMarker: *mustPai("6s"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDora(tt.args.doraMarker)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDora() = %v, want %v", got, tt.want)
			}
		})
	}
}
