package inbound

import (
	"reflect"
	"testing"
)

func TestNewEndGame(t *testing.T) {
	tests := []struct {
		name    string
		want    *EndGame
		wantErr bool
	}{
		{
			name:    "test NewEndGame()",
			want:    &EndGame{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEndGame()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndGame() = %v, want %v", got, tt.want)
			}
		})
	}
}
