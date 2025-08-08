package inbound

import (
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name    string
		want    *Error
		wantErr bool
	}{
		{
			name:    "test NewError()",
			want:    &Error{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewError()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}
