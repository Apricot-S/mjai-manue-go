package tile_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func TestNewTileFromID(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantID  int
		wantErr bool
	}{
		{
			name:    "0 is minimum valid ID",
			id:      0,
			wantID:  0,
			wantErr: false,
		},
		{
			name:    "37 is maximum valid ID",
			id:      37,
			wantID:  37,
			wantErr: false,
		},
		{
			name:    "-1 is an invalid ID",
			id:      -1,
			wantID:  -1,
			wantErr: true,
		},
		{
			name:    "38 is an invalid ID",
			id:      38,
			wantID:  38,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tile.NewTileFromID(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewTileFromID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewTileFromID() succeeded unexpectedly")
			}
			if got.ID() != tt.wantID {
				t.Errorf("NewTileFromID().ID() = %v, want %v", got, tt.wantID)
			}
		})
	}
}

func TestNewTileFromCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantID  int
		wantErr bool
	}{
		{
			name:    "1m is ID 0",
			code:    "1m",
			wantID:  0,
			wantErr: false,
		},
		{
			name:    "C is ID 33",
			code:    "C",
			wantID:  33,
			wantErr: false,
		},
		{
			name:    "5sr is ID 36",
			code:    "5sr",
			wantID:  36,
			wantErr: false,
		},
		{
			name:    "? is ID 37",
			code:    "?",
			wantID:  37,
			wantErr: false,
		},
		{
			name:    "1z is an invalid code",
			code:    "1z",
			wantID:  27,
			wantErr: true,
		},
		{
			name:    "0m is an invalid code",
			code:    "0m",
			wantID:  34,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tile.NewTileFromCode(tt.code)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewTileFromCode() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewTileFromCode() succeeded unexpectedly")
			}
			if got.ID() != tt.wantID {
				t.Errorf("NewTileFromCode().ID() = %v, want %v", got, tt.wantID)
			}
		})
	}
}

func TestTile_Code(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want string
	}{
		{
			name: "1m is ID 0",
			id:   0,
			want: "1m",
		},
		{
			name: "C is ID 33",
			id:   33,
			want: "C",
		},
		{
			name: "5sr is ID 36",
			id:   36,
			want: "5sr",
		},
		{
			name: "? is ID 37",
			id:   37,
			want: "?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti, err := tile.NewTileFromID(tt.id)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := ti.Code()
			if got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_Color(t *testing.T) {
	tests := []struct {
		name string
		code string
		want rune
	}{
		{
			name: "1m's color is m",
			code: "1m",
			want: 'm',
		},
		{
			name: "1p's color is p",
			code: "1p",
			want: 'p',
		},
		{
			name: "1s's color is s",
			code: "1s",
			want: 's',
		},
		{
			name: "E's color is t",
			code: "E",
			want: 't',
		},
		{
			name: "5mr's color is m",
			code: "5mr",
			want: 'm',
		},
		{
			name: "?'s color is ?",
			code: "?",
			want: '?',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.Color()
			if got != tt.want {
				t.Errorf("Color() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_Number(t *testing.T) {
	tests := []struct {
		name string
		code string
		want int
	}{
		{
			name: "5p's number is 5",
			code: "5p",
			want: 5,
		},
		{
			name: "5pr's number is 5",
			code: "5pr",
			want: 5,
		},
		{
			name: "E's number is 1",
			code: "E",
			want: 1,
		},
		{
			name: "?'s number is 0",
			code: "?",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.Number()
			if got != tt.want {
				t.Errorf("Number() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_IsRed(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "5m is not red",
			code: "5m",
			want: false,
		},
		{
			name: "5mr is red",
			code: "5mr",
			want: true,
		},
		{
			name: "P is not red",
			code: "P",
			want: false,
		},
		{
			name: "? is not red",
			code: "?",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.IsRed()
			if got != tt.want {
				t.Errorf("IsRed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_IsSuits(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "9s is suits",
			code: "9s",
			want: true,
		},
		{
			name: "5mr is suits",
			code: "5mr",
			want: true,
		},
		{
			name: "E is not suits",
			code: "E",
			want: false,
		},
		{
			name: "? is not suits",
			code: "?",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.IsSuits()
			if got != tt.want {
				t.Errorf("IsSuits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_IsHonors(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "9s is not honors",
			code: "9s",
			want: false,
		},
		{
			name: "5mr is not honors",
			code: "5mr",
			want: false,
		},
		{
			name: "E is honors",
			code: "E",
			want: true,
		},
		{
			name: "? is not honors",
			code: "?",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.IsHonors()
			if got != tt.want {
				t.Errorf("IsHonors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_IsUnknown(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "9s is not unknown",
			code: "9s",
			want: false,
		},
		{
			name: "5mr is not unknown",
			code: "5mr",
			want: false,
		},
		{
			name: "E is not unknown",
			code: "E",
			want: false,
		},
		{
			name: "? is unknown",
			code: "?",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.IsUnknown()
			if got != tt.want {
				t.Errorf("IsUnknown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_Next(t *testing.T) {
	tests := []struct {
		name string
		code string
		n    int
		want *tile.Tile
	}{
		{
			name: "no next tile for unknown",
			code: "?",
			n:    0,
			want: nil,
		},
		{
			name: "no next tile for honors",
			code: "E",
			n:    1,
			want: nil,
		},
		{
			name: "no next tile for 9",
			code: "9m",
			n:    1,
			want: nil,
		},
		{
			name: "no prev tile for 1",
			code: "1s",
			n:    -1,
			want: nil,
		},
		{
			name: "0 after of 9 is 9",
			code: "9p",
			n:    0,
			want: tile.MustTileFromCode("9p"),
		},
		{
			name: "1 after of 8 is 9",
			code: "8p",
			n:    1,
			want: tile.MustTileFromCode("9p"),
		},
		{
			name: "2 before of 3 is 1",
			code: "3m",
			n:    -2,
			want: tile.MustTileFromCode("1m"),
		},
		{
			name: "1 next of 5r is 6",
			code: "5pr",
			n:    1,
			want: tile.MustTileFromCode("6p"),
		},
		{
			name: "0 next of 5r is 5",
			code: "5pr",
			n:    0,
			want: tile.MustTileFromCode("5p"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.Next(tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_AddRed(t *testing.T) {
	tests := []struct {
		name string
		code string
		want *tile.Tile
	}{
		{
			name: "unknown will not be red",
			code: "?",
			want: tile.MustTileFromCode("?"),
		},
		{
			name: "suits other than 5 will not be red",
			code: "4s",
			want: tile.MustTileFromCode("4s"),
		},
		{
			name: "red 5 stays red",
			code: "5sr",
			want: tile.MustTileFromCode("5sr"),
		},
		{
			name: "normal 5 becomes red 5",
			code: "5s",
			want: tile.MustTileFromCode("5sr"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.AddRed()
			if *got != *tt.want {
				t.Errorf("AddRed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_RemoveRed(t *testing.T) {
	tests := []struct {
		name string
		code string
		want *tile.Tile
	}{
		{
			name: "no change unknown",
			code: "?",
			want: tile.MustTileFromCode("?"),
		},
		{
			name: "no change except for 5 suits",
			code: "4s",
			want: tile.MustTileFromCode("4s"),
		},
		{
			name: "normal 5 stays normal",
			code: "5m",
			want: tile.MustTileFromCode("5m"),
		},
		{
			name: "red 5 becomes normal 5",
			code: "5mr",
			want: tile.MustTileFromCode("5m"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.RemoveRed()
			if *got != *tt.want {
				t.Errorf("RemoveRed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_HasSameSymbol(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		otherCode string
		want      bool
	}{
		{
			name:      "1m-1m same",
			code:      "1m",
			otherCode: "1m",
			want:      true,
		},
		{
			name:      "1m-1p not same",
			code:      "1m",
			otherCode: "1p",
			want:      false,
		},
		{
			name:      "5p-5pr same",
			code:      "5p",
			otherCode: "5pr",
			want:      true,
		},
		{
			name:      "?-? same",
			code:      "?",
			otherCode: "?",
			want:      true,
		},
		{
			name:      "1m-? not same",
			code:      "1m",
			otherCode: "?",
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			other := tile.MustTileFromCode(tt.otherCode)
			got := ti.HasSameSymbol(other)
			if got != tt.want {
				t.Errorf("HasSameSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortTiles(t *testing.T) {
	names := [...]string{
		"?",
		"5sr", "5pr", "5mr",
		"C", "F", "P", "N", "W", "S", "E",
		"9s", "8s", "7s", "6s", "5s", "4s", "3s", "2s", "1s",
		"9p", "8p", "7p", "6p", "5p", "4p", "3p", "2p", "1p",
		"9m", "8m", "7m", "6m", "5m", "4m", "3m", "2m", "1m",
	}

	tiles := make(tile.Tiles, 0, len(names))
	for _, name := range names {
		t := tile.MustTileFromCode(name)
		tiles = append(tiles, *t)
	}
	sort.Sort(tiles)

	sortedNames := [...]string{
		"1m", "2m", "3m", "4m", "5m", "5mr", "6m", "7m", "8m", "9m",
		"1p", "2p", "3p", "4p", "5p", "5pr", "6p", "7p", "8p", "9p",
		"1s", "2s", "3s", "4s", "5s", "5sr", "6s", "7s", "8s", "9s",
		"E", "S", "W", "N", "P", "F", "C",
		"?",
	}

	for i, sortedName := range sortedNames {
		if tiles[i].Code() != sortedName {
			t.Errorf("Expected %s but got %s", sortedName, tiles[i].Code())
		}
	}
}
