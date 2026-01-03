package player_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/player"
)

func TestNewID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		wantIndex int
		wantErr   bool
	}{
		{
			name:      "valid ID: 0",
			id:        0,
			wantIndex: 0,
			wantErr:   false,
		},
		{
			name:      "valid ID: 3",
			id:        3,
			wantIndex: 3,
			wantErr:   false,
		},
		{
			name:      "invalid ID: -1",
			id:        -1,
			wantIndex: -1,
			wantErr:   true,
		},
		{
			name:      "invalid ID: 4",
			id:        4,
			wantIndex: 4,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := player.NewID(tt.id)
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
