package tile_test

import (
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

func TestTile_IsSuit(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "9s is suit",
			code: "9s",
			want: true,
		},
		{
			name: "5mr is suit",
			code: "5mr",
			want: true,
		},
		{
			name: "E is not suit",
			code: "E",
			want: false,
		},
		{
			name: "? is not suit",
			code: "?",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.IsSuit()
			if got != tt.want {
				t.Errorf("IsSuit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTile_IsHonor(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "9s is not honor",
			code: "9s",
			want: false,
		},
		{
			name: "5mr is not honor",
			code: "5mr",
			want: false,
		},
		{
			name: "E is honor",
			code: "E",
			want: true,
		},
		{
			name: "? is not honor",
			code: "?",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tile.MustTileFromCode(tt.code)
			got := ti.IsHonor()
			if got != tt.want {
				t.Errorf("IsHonor() = %v, want %v", got, tt.want)
			}
		})
	}
}
