package seat_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func TestNewSeat(t *testing.T) {
	tests := []struct {
		name      string
		index     int
		wantIndex int
		wantErr   bool
	}{
		{
			name:      "valid Seat: 0",
			index:     0,
			wantIndex: 0,
			wantErr:   false,
		},
		{
			name:      "valid Seat: 3",
			index:     3,
			wantIndex: 3,
			wantErr:   false,
		},
		{
			name:      "invalid Seat: -1",
			index:     -1,
			wantIndex: -1,
			wantErr:   true,
		},
		{
			name:      "invalid Seat: 4",
			index:     4,
			wantIndex: 4,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := seat.NewSeat(tt.index)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewSeat() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewSeat() succeeded unexpectedly")
			}
			if got.Index() != tt.wantIndex {
				t.Errorf("NewSeat().Index() = %v, want %v", got, tt.wantIndex)
			}
		})
	}
}

func TestSeat_IsShimochaOf(t *testing.T) {
	tests := []struct {
		name   string
		seat   seat.Seat
		target seat.Seat
		want   bool
	}{
		{
			name:   "1 is shimocha of 0",
			seat:   *seat.MustSeat(1),
			target: *seat.MustSeat(0),
			want:   true,
		},
		{
			name:   "0 is shimocha of 3",
			seat:   *seat.MustSeat(0),
			target: *seat.MustSeat(3),
			want:   true,
		},
		{
			name:   "2 is not shimocha of 0",
			seat:   *seat.MustSeat(2),
			target: *seat.MustSeat(0),
			want:   false,
		},
		{
			name:   "same seat is not shimocha",
			seat:   *seat.MustSeat(0),
			target: *seat.MustSeat(0),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.seat.IsShimochaOf(tt.target); got != tt.want {
				t.Errorf("IsShimochaOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
