package block_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func TestNewSequence(t *testing.T) {
	tests := []struct {
		name    string
		tile    tile.Tile
		want    []tile.Tile
		wantErr bool
	}{
		{
			name:    "cannot create sequence starting with 8",
			tile:    *tile.MustTileFromCode("8m"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create sequence from honors",
			tile:    *tile.MustTileFromCode("C"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create sequence from unknown tile",
			tile:    *tile.MustTileFromCode("?"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create sequence from red five",
			tile:    *tile.MustTileFromCode("5mr"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := block.NewSequence(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewSequence() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewSequence() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewSequence().ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTriplet(t *testing.T) {
	tests := []struct {
		name    string
		tile    tile.Tile
		want    []tile.Tile
		wantErr bool
	}{
		{
			name:    "C triplet",
			tile:    *tile.MustTileFromCode("C"),
			want:    []tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			wantErr: false,
		},
		{
			name:    "cannot create triplet from unknown tile",
			tile:    *tile.MustTileFromCode("?"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create triplet from red five",
			tile:    *tile.MustTileFromCode("5mr"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := block.NewTriplet(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewTriplet() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewTriplet() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewTriplet().ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQuad(t *testing.T) {
	tests := []struct {
		name    string
		tile    tile.Tile
		want    []tile.Tile
		wantErr bool
	}{
		{
			name:    "C quad",
			tile:    *tile.MustTileFromCode("C"),
			want:    []tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			wantErr: false,
		},
		{
			name:    "cannot create quad from unknown tile",
			tile:    *tile.MustTileFromCode("?"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create quad from red five",
			tile:    *tile.MustTileFromCode("5mr"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := block.NewQuad(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewQuad() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewQuad() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewQuad().ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPair(t *testing.T) {
	tests := []struct {
		name    string
		tile    tile.Tile
		want    []tile.Tile
		wantErr bool
	}{
		{
			name:    "C pair",
			tile:    *tile.MustTileFromCode("C"),
			want:    []tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			wantErr: false,
		},
		{
			name:    "cannot create pair from unknown tile",
			tile:    *tile.MustTileFromCode("?"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create pair from red five",
			tile:    *tile.MustTileFromCode("5mr"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := block.NewPair(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewPair() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewPair() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewPair().ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
