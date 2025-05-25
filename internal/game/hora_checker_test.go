package game

import (
	"testing"
)

func convertStrToPaiSetForTest(paiStr string) *PaiSet {
	pais, _ := StrToPais(paiStr)
	ps, _ := NewPaiSetWithPais(pais)
	return ps
}

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
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{}

	{
		paiStr := "1m 9m 9m 1p 9p 1s 9s E S W N P F C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1m 9m 1p 9p 1s 9s E E S W N P F C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1m 9m 1p 9p 1s 9s E S W N P F C C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1m 2m 9m 1p 9p 1s 9s E S W N P F C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false})
	}
	{
		paiStr := "1m 9m 1p 9p 1s 9s E E E W N P F C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHoraFormKokushimuso(tt.args.ps); got != tt.want {
				t.Errorf("isHoraFormKokushimuso() = %v, want %v", got, tt.want)
			}
		})
	}
}
