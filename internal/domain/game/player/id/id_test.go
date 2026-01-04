package id_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/player/id"
)

func TestNewID(t *testing.T) {
	tests := []struct {
		name      string
		index     int
		wantIndex int
		wantErr   bool
	}{
		{
			name:      "valid ID: 0",
			index:     0,
			wantIndex: 0,
			wantErr:   false,
		},
		{
			name:      "valid ID: 3",
			index:     3,
			wantIndex: 3,
			wantErr:   false,
		},
		{
			name:      "invalid ID: -1",
			index:     -1,
			wantIndex: -1,
			wantErr:   true,
		},
		{
			name:      "invalid ID: 4",
			index:     4,
			wantIndex: 4,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := id.NewID(tt.index)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewID() succeeded unexpectedly")
			}
			if got.Index() != tt.wantIndex {
				t.Errorf("NewID().Index() = %v, want %v", got, tt.wantIndex)
			}
		})
	}
}
