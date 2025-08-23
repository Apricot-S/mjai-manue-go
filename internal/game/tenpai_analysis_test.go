package game

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func TestIsTenpaiGeneral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{
			name:    "empty : An empty hand is one step away from being a pair wait -> noten",
			input:   "",
			want:    false,
			wantErr: false,
		},
		{
			name:    "chitoitsu",
			input:   "1m 1m 8m 8m 2p 8p 8p 5s 5s E E C C",
			want:    false,
			wantErr: false,
		},
		{
			name:    "thirteen orphans",
			input:   "1m 9m 1p 9p 1s 9s E S W N P F C",
			want:    false,
			wantErr: false,
		},
		{
			name:    "tenpai",
			input:   "1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S",
			want:    true,
			wantErr: false,
		},
		{
			name:    "win",
			input:   "1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S S",
			want:    true,
			wantErr: false,
		},
		{
			name:    "with meld",
			input:   "1m 2m 3m 4p 5pr 6p 7s 8s 9s E",
			want:    true,
			wantErr: false,
		},
		{
			name:    "5 identical tiles",
			input:   "1m 1m 1m 1m 1m",
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pais, err := base.StrToPais(tt.input)
			if err != nil {
				t.Errorf("StrToPais() error = %v", err)
				return
			}
			paiSet, err := base.NewPaiSet(pais)
			if err != nil {
				t.Errorf("NewPaiSet() error = %v", err)
				return
			}

			got, err := IsTenpaiGeneral(paiSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsTenpaiGeneral() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsTenpaiGeneral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTenpaiAll(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{
			name:    "empty : An empty hand is one step away from being a pair wait -> noten",
			input:   "",
			want:    false,
			wantErr: false,
		},
		{
			name:    "chitoitsu",
			input:   "1m 1m 8m 8m 2p 8p 8p 5s 5s E E C C",
			want:    true,
			wantErr: false,
		},
		{
			name:    "thirteen orphans",
			input:   "1m 9m 1p 9p 1s 9s E S W N P F C",
			want:    true,
			wantErr: false,
		},
		{
			name:    "tenpai",
			input:   "1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S",
			want:    true,
			wantErr: false,
		},
		{
			name:    "win",
			input:   "1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S S",
			want:    true,
			wantErr: false,
		},
		{
			name:    "with meld",
			input:   "1m 2m 3m 4p 5pr 6p 7s 8s 9s E",
			want:    true,
			wantErr: false,
		},
		{
			name:    "5 identical tiles",
			input:   "1m 1m 1m 1m 1m",
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pais, err := base.StrToPais(tt.input)
			if err != nil {
				t.Errorf("StrToPais() error = %v", err)
				return
			}
			paiSet, err := base.NewPaiSet(pais)
			if err != nil {
				t.Errorf("NewPaiSet() error = %v", err)
				return
			}

			got, err := IsTenpaiAll(paiSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsTenpaiAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsTenpaiAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWaitedPaisAll(t *testing.T) {
	tests := []struct {
		name    string
		tehai   string
		want    string
		wantErr bool
	}{
		{
			name:    "tenpai",
			tehai:   "1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S",
			want:    "E S",
			wantErr: false,
		},
		{
			name:    "chitoitsu",
			tehai:   "1m 1m 8m 8m 2p 8p 8p 5s 5s E E C C",
			want:    "2p",
			wantErr: false,
		},
		{
			name:    "thirteen orphans",
			tehai:   "9m 9m 1p 9p 1s 9s E S W N P F C",
			want:    "1m",
			wantErr: false,
		},
		{
			name:    "thirteen orphans 13 waits",
			tehai:   "1m 9m 1p 9p 1s 9s E S W N P F C",
			want:    "1m 9m 1p 9p 1s 9s E S W N P F C",
			wantErr: false,
		},
		{
			name:    "An empty hand does not have tusmo tile",
			tehai:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "win",
			tehai:   "1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S S",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tehaiPais, err := base.StrToPais(tt.tehai)
			if err != nil {
				t.Errorf("StrToPais() error = %v", err)
				return
			}
			tehaiPaiSet, err := base.NewPaiSet(tehaiPais)
			if err != nil {
				t.Errorf("NewPaiSet() error = %v", err)
				return
			}

			var waited *base.PaiSet
			if tt.want == "" {
				waited = nil
			} else {
				waitedPais, err := base.StrToPais(tt.want)
				if err != nil {
					t.Errorf("StrToPais() error = %v", err)
					return
				}
				waited, err = base.NewPaiSet(waitedPais)
				if err != nil {
					t.Errorf("NewPaiSet() error = %v", err)
					return
				}
			}

			got, err := GetWaitedPaisAll(tehaiPaiSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWaitedPaisAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, waited) {
				t.Errorf("GetWaitedPaisAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
