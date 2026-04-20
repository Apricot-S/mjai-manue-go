package mjai_test

import (
	"io"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai"
)

func TestParseTsumo(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		want    *event.Draw
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := mjai.ParseTsumo(tt.r)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseTsumo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseTsumo() succeeded unexpectedly")
			}
			if *got != *tt.want {
				t.Errorf("ParseTsumo() = %v, want %v", got, tt.want)
			}
		})
	}
}
