package game

import (
	"testing"
)

func TestIsHoraForm(t *testing.T) {
	type args struct {
		ps *PaiSet
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsHoraForm(tt.args.ps)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsHoraForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsHoraForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHoraFormGeneral(t *testing.T) {
	type args struct {
		ps         *PaiSet
		numMentsus int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHoraFormGeneral(tt.args.ps, tt.args.numMentsus); got != tt.want {
				t.Errorf("isHoraFormGeneral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHoraFormChitoitsu(t *testing.T) {
	type args struct {
		ps *PaiSet
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHoraFormChitoitsu(tt.args.ps); got != tt.want {
				t.Errorf("isHoraFormChitoitsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHoraFormKokushimuso(t *testing.T) {
	type args struct {
		ps *PaiSet
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHoraFormKokushimuso(tt.args.ps); got != tt.want {
				t.Errorf("isHoraFormKokushimuso() = %v, want %v", got, tt.want)
			}
		})
	}
}
