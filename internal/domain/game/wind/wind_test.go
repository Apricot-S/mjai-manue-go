package wind_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestNewWind(t *testing.T) {
	tests := []struct {
		name    string
		w       string
		want    wind.Wind
		wantErr bool
	}{
		{
			name:    "East",
			w:       "E",
			want:    wind.East,
			wantErr: false,
		},
		{
			name:    "South",
			w:       "S",
			want:    wind.South,
			wantErr: false,
		},
		{
			name:    "West",
			w:       "W",
			want:    wind.West,
			wantErr: false,
		},
		{
			name:    "North",
			w:       "N",
			want:    wind.North,
			wantErr: false,
		},
		{
			name:    "invalid",
			w:       "East",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := wind.NewWind(tt.w)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewWind() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewWind() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("NewWind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWind_String(t *testing.T) {
	tests := []struct {
		name string
		w    wind.Wind
		want string
	}{
		{
			name: "East",
			w:    wind.East,
			want: "E",
		},
		{
			name: "South",
			w:    wind.South,
			want: "S",
		},
		{
			name: "West",
			w:    wind.West,
			want: "W",
		},
		{
			name: "North",
			w:    wind.North,
			want: "N",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.w.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWind_Next(t *testing.T) {
	tests := []struct {
		name string
		w    wind.Wind
		want wind.Wind
	}{
		{
			name: "East -> South",
			w:    wind.East,
			want: wind.South,
		},
		{
			name: "South -> West",
			w:    wind.South,
			want: wind.West,
		},
		{
			name: "West -> North",
			w:    wind.West,
			want: wind.North,
		},
		{
			name: "North -> East",
			w:    wind.North,
			want: wind.East,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.w.Next()
			if got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
