package inbound

import (
	"reflect"
	"testing"
)

func TestNewHello(t *testing.T) {
	tests := []struct {
		name    string
		want    *Hello
		wantErr bool
	}{
		{
			name:    "test NewHello()",
			want:    &Hello{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHello()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHello() = %v, want %v", got, tt.want)
			}
		})
	}
}
