package outbound

import (
	"reflect"
	"testing"
)

func TestNewNone(t *testing.T) {
	tests := []struct {
		name    string
		want    *None
		wantErr bool
	}{
		{
			name:    "test NewNone()",
			want:    &None{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNone()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNone() = %v, want %v", got, tt.want)
			}
		})
	}
}
