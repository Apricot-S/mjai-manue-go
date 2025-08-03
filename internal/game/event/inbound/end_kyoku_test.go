package inbound

import (
	"reflect"
	"testing"
)

func TestNewEndKyoku(t *testing.T) {
	tests := []struct {
		name    string
		want    *EndKyoku
		wantErr bool
	}{
		{
			name:    "test NewEndKyoku()",
			want:    &EndKyoku{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEndKyoku()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndKyoku() = %v, want %v", got, tt.want)
			}
		})
	}
}
