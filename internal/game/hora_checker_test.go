package game

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func convertStrToPaiSetForTest(paiStr string) *base.PaiSet {
	pais, _ := base.StrToPais(paiStr)
	ps, _ := base.NewPaiSet(pais)
	return ps
}

func TestIsHoraForm(t *testing.T) {
	type args struct {
		ps *base.PaiSet
	}
	type testCase struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}
	tests := []testCase{}

	{
		paiStr := "3m 3m 1p 2p 3p 1s 2s 3s 4s 5s 6s 7s 8s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true, wantErr: false})
	}
	{
		paiStr := "1m 1m E E E"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true, wantErr: false})
	}
	{
		paiStr := "1p 2p 3p 4p 5p 1s 2s 3s 4s 5s 6s 7s 8s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false, wantErr: false})
	}
	{
		paiStr := "2p 2p 2p 2s 3s 4s 6s 7s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false, wantErr: false})
	}
	{
		paiStr := "1p 1p 9p 9p 1s 1s 3s 3s 5s 5s 7s 7s 9s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true, wantErr: false})
	}
	{
		paiStr := "1p 1p 1p 1p 3p 3p 4p 4p 5p 5p 6p 6p 7p 7p"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false, wantErr: false})
	}
	{
		paiStr := "1m 9m 1p 9p 1s 9s E S W N P F C C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true, wantErr: false})
	}
	{
		paiStr := "1m 2m 9m 1p 9p 1s 9s E S W N P F C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false, wantErr: false})
	}
	{
		paiStr := "3m 3m 3m 1p 2p 3p 1s 2s 3s 4s 5s 6s 7s 8s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false, wantErr: true})
	}
	{
		paiStr := ""
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false, wantErr: true})
	}
	{
		paiStr := "1p 1p 1p 1p 1p 1s 2s 3s 4s 5s 6s 7s 8s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false, wantErr: true})
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
		ps *base.PaiSet
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{}

	{
		paiStr := "3m 3m 1p 2p 3p 1s 2s 3s 4s 5s 6s 7s 8s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1p 1p 1p 2p 3p 1s 2s 3s 4s 5s 6s 7s 8s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1s 2s 3s 4s 5s 6s 7s 8s 9s E E C C C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1s 2s 3s 4s 5s 5s 5s 6s 7s 7s 8s 8s 9s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1s 1s 1s 2s 3s 4s 5s 6s 7s 8s 8s 9s 9s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "3m 4m 5m 7p 8p 9p 2s 3s 3s 3s 3s 4s P P"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1m 1m 1m 2m 3m 3m 3m 4m 4m 4m 5m 5m 9m 9m"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1m 1m"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1m 1m E E E"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}

	{
		paiStr := "1p 2p 3p 4p 5p 1s 2s 3s 4s 5s 6s 7s 8s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false})
	}
	{
		paiStr := "1p 1p 1p 4p 5p 1s 1s 1s 2s 2s 2s 4s 4s 4s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false})
	}
	{
		paiStr := "1s 1s 1s 2s 2s 2s 3s 3s 3s 5s 6s 8s 8s 8s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false})
	}
	{
		paiStr := "2p 2p 2p 2s 3s 4s 6s 7s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHoraFormGeneral(tt.args.ps); got != tt.want {
				t.Errorf("isHoraFormGeneral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSingleColorHoraFormWithoutPair(t *testing.T) {
	type args struct {
		ps []int
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{}

	{
		paiStr := "1m 1m 1m 1m 2m 2m 3m 3m 4m"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps[:]}, want: true})
	}
	{
		paiStr := "1m 1m 1m 1m 2m 2m 3m 3m"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps[:]}, want: false})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSingleColorHoraFormWithoutPair(tt.args.ps); got != tt.want {
				t.Errorf("isSingleColorHoraFormWithoutPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSingleColorHoraFormWithPair(t *testing.T) {
	type args struct {
		ps []int
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{}

	{
		paiStr := "1m 1m 1m 1m 2m 2m 3m 3m 4m"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps[:]}, want: false})
	}
	{
		paiStr := "1m 1m 1m 1m 2m 2m 3m 3m"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps[:]}, want: true})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSingleColorHoraFormWithPair(tt.args.ps); got != tt.want {
				t.Errorf("isSingleColorHoraFormWithPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHoraFormChitoitsu(t *testing.T) {
	type args struct {
		ps *base.PaiSet
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{}

	{
		paiStr := "1p 1p 9p 9p 1s 1s 3s 3s 5s 5s 7s 7s 9s 9s"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1m 1m 1p 1p 9p 9p 2s 2s 4s 4s S S C C"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1m 1m 2m 2m 3m 3m 4m 4m 5m 5m 6m 6m 7m 7m"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: true})
	}
	{
		paiStr := "1p 1p 1p 1p 3p 3p 4p 4p 5p 5p 6p 6p 7p 7p"
		ps := convertStrToPaiSetForTest(paiStr)
		tests = append(tests, testCase{name: paiStr, args: args{ps: ps}, want: false})
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
		ps *base.PaiSet
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
