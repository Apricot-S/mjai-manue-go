package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/testutil"
)

func TestNewDora(t *testing.T) {
	type args struct {
		doraMarker base.Pai
	}
	tests := []struct {
		name    string
		args    args
		want    *Dora
		wantErr bool
	}{
		{
			name: "valid dora marker",
			args: args{
				doraMarker: *testutil.MustPai("6s"),
			},
			want: &Dora{
				DoraMarker: *testutil.MustPai("6s"),
			},
			wantErr: false,
		},
		{
			name: "unknown dora marker",
			args: args{
				doraMarker: *testutil.MustPai("?"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDora(tt.args.doraMarker)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDora() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDora() = %v, want %v", got, tt.want)
			}
		})
	}
}
