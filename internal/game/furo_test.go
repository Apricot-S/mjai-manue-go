package game

import (
	"reflect"
	"testing"
)

func TestNewFuro(t *testing.T) {
	type args struct {
		t        FuroType
		taken    *Pai
		consumed []Pai
		target   *int
	}
	type testCase struct {
		name    string
		args    args
		want    *Furo
		wantErr bool
	}
	tests := []testCase{}

	// Chi valid
	for _, strs := range [][2]string{{"1m", "2m 3m"}, {"1m", "3m 2m"}, {"5mr", "4m 6m"}} {
		taken, _ := NewPaiWithName(strs[0])
		consumed, _ := StrToPais(strs[1])
		target := 0
		want, _ := NewFuro(Chi, taken, consumed, &target)
		wantErr := false
		tests = append(tests, testCase{"Chi", args{Chi, taken, consumed, &target}, want, wantErr})
	}

	// Chi invalid
	{
		taken, _ := NewPaiWithName("1m")
		consumed2, _ := StrToPais("2m 3m")
		consumed1, _ := StrToPais("2m")
		consumed3, _ := StrToPais("2m 3m 4m")
		target := 0

		tests = append(tests, testCase{"Chi invalid: taken: nil", args{Chi, nil, consumed2, &target}, nil, true})
		tests = append(tests, testCase{"Chi invalid: taken: nil, target: nil", args{Chi, nil, consumed2, nil}, nil, true})
		tests = append(tests, testCase{"Chi invalid: taken: nil, consumed: 1", args{Chi, nil, consumed1, &target}, nil, true})
		tests = append(tests, testCase{"Chi invalid: taken: nil, consumed: 1, target: nil", args{Chi, nil, consumed1, nil}, nil, true})
		tests = append(tests, testCase{"Chi invalid: taken: nil, consumed: 3", args{Chi, nil, consumed3, &target}, nil, true})
		tests = append(tests, testCase{"Chi invalid: taken: nil, consumed: 3, target: nil", args{Chi, nil, consumed3, nil}, nil, true})
		tests = append(tests, testCase{"Chi invalid: target: nil", args{Chi, taken, consumed2, nil}, nil, true})
		tests = append(tests, testCase{"Chi invalid: consumed: 1", args{Chi, taken, consumed1, &target}, nil, true})
		tests = append(tests, testCase{"Chi invalid: consumed: 1, target: nil", args{Chi, taken, consumed1, nil}, nil, true})
		tests = append(tests, testCase{"Chi invalid: consumed: 3", args{Chi, taken, consumed3, &target}, nil, true})
		tests = append(tests, testCase{"Chi invalid: consumed: 3, target: nil", args{Chi, taken, consumed3, nil}, nil, true})
	}

	// Pon valid
	for _, strs := range [][2]string{{"E", "E E"}, {"5mr", "5m 5m"}, {"5p", "5pr 5p"}, {"5s", "5s 5sr"}} {
		taken, _ := NewPaiWithName(strs[0])
		consumed, _ := StrToPais(strs[1])
		target := 1
		want, _ := NewFuro(Pon, taken, consumed, &target)
		wantErr := false
		tests = append(tests, testCase{"Pon", args{Pon, taken, consumed, &target}, want, wantErr})
	}

	// Pon invalid
	{
		taken, _ := NewPaiWithName("1m")
		consumed2, _ := StrToPais("1m 1m")
		consumed1, _ := StrToPais("1m")
		consumed3, _ := StrToPais("1m 1m 1m")
		target := 0

		tests = append(tests, testCase{"Pon invalid: taken: nil", args{Pon, nil, consumed2, &target}, nil, true})
		tests = append(tests, testCase{"Pon invalid: taken: nil, target: nil", args{Pon, nil, consumed2, nil}, nil, true})
		tests = append(tests, testCase{"Pon invalid: taken: nil, consumed: 1", args{Pon, nil, consumed1, &target}, nil, true})
		tests = append(tests, testCase{"Pon invalid: taken: nil, consumed: 1, target: nil", args{Pon, nil, consumed1, nil}, nil, true})
		tests = append(tests, testCase{"Pon invalid: taken: nil, consumed: 3", args{Pon, nil, consumed3, &target}, nil, true})
		tests = append(tests, testCase{"Pon invalid: taken: nil, consumed: 3, target: nil", args{Pon, nil, consumed3, nil}, nil, true})
		tests = append(tests, testCase{"Pon invalid: target: nil", args{Pon, taken, consumed2, nil}, nil, true})
		tests = append(tests, testCase{"Pon invalid: consumed: 1", args{Pon, taken, consumed1, &target}, nil, true})
		tests = append(tests, testCase{"Pon invalid: consumed: 1, target: nil", args{Pon, taken, consumed1, nil}, nil, true})
		tests = append(tests, testCase{"Pon invalid: consumed: 3", args{Pon, taken, consumed3, &target}, nil, true})
		tests = append(tests, testCase{"Pon invalid: consumed: 3, target: nil", args{Pon, taken, consumed3, nil}, nil, true})
	}

	// Daiminkan valid
	for _, strs := range [][2]string{{"E", "E E E"}, {"5mr", "5m 5m 5m"}, {"5p", "5pr 5p 5p"}, {"5s", "5s 5sr 5s"}, {"5s", "5s 5s 5sr"}} {
		taken, _ := NewPaiWithName(strs[0])
		consumed, _ := StrToPais(strs[1])
		target := 2
		want, _ := NewFuro(Daiminkan, taken, consumed, &target)
		wantErr := false

		tests = append(tests, testCase{"Daiminkan", args{Daiminkan, taken, consumed, &target}, want, wantErr})
	}

	// Daiminkan invalid
	{
		taken, _ := NewPaiWithName("1m")
		consumed3, _ := StrToPais("1m 1m 1m")
		consumed2, _ := StrToPais("1m 1m")
		consumed4, _ := StrToPais("1m 1m 1m 1m")
		target := 0

		tests = append(tests, testCase{"Daiminkan invalid: taken: nil", args{Daiminkan, nil, consumed3, &target}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: taken: nil, target: nil", args{Daiminkan, nil, consumed3, nil}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: taken: nil, consumed: 2", args{Daiminkan, nil, consumed2, &target}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: taken: nil, consumed: 2, target: nil", args{Daiminkan, nil, consumed2, nil}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: taken: nil, consumed: 4", args{Daiminkan, nil, consumed4, &target}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: taken: nil, consumed: 4, target: nil", args{Daiminkan, nil, consumed4, nil}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: target: nil", args{Daiminkan, taken, consumed3, nil}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: consumed: 2", args{Daiminkan, taken, consumed2, &target}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: consumed: 2, target: nil", args{Daiminkan, taken, consumed2, nil}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: consumed: 4", args{Daiminkan, taken, consumed4, &target}, nil, true})
		tests = append(tests, testCase{"Daiminkan invalid: consumed: 4, target: nil", args{Daiminkan, taken, consumed4, nil}, nil, true})
	}

	// Ankan valid
	for _, consumed := range []string{"E E E E", "5mr 5m 5m 5m", "5p 5pr 5p 5p", "5s 5s 5sr 5s", "5s 5s 5s 5sr"} {
		consumed, _ := StrToPais(consumed)
		want, _ := NewFuro(Ankan, nil, consumed, nil)
		wantErr := false

		tests = append(tests, testCase{"Ankan", args{Ankan, nil, consumed, nil}, want, wantErr})
	}

	// Ankan invalid
	{
		taken, _ := NewPaiWithName("1m")
		consumed4, _ := StrToPais("1m 1m 1m 1m")
		consumed3, _ := StrToPais("1m 1m 1m")
		consumed5, _ := StrToPais("1m 1m 1m 1m 1m")
		target := 0

		tests = append(tests, testCase{"Ankan invalid: taken: not nil, target: not nil", args{Ankan, taken, consumed4, &target}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: taken: not nil", args{Ankan, taken, consumed4, nil}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: taken: not nil, consumed: 3, target: not nil", args{Ankan, taken, consumed3, &target}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: taken: not nil, consumed: 3", args{Ankan, taken, consumed3, nil}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: taken: not nil, consumed: 5, target: not nil", args{Ankan, taken, consumed5, &target}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: taken: not nil, consumed: 5", args{Ankan, taken, consumed5, nil}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: target: not nil", args{Ankan, nil, consumed4, &target}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: consumed: 3, target: not nil", args{Ankan, nil, consumed3, &target}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: consumed: 3", args{Ankan, nil, consumed3, nil}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: consumed: 5, target: not nil", args{Ankan, nil, consumed5, &target}, nil, true})
		tests = append(tests, testCase{"Ankan invalid: consumed: 5", args{Ankan, nil, consumed5, nil}, nil, true})
	}

	// Kakan valid target: present
	for _, strs := range [][2]string{{"E", "E E E"}, {"5mr", "5m 5m 5m"}, {"5p", "5pr 5p 5p"}, {"5s", "5s 5sr 5s"}, {"5s", "5s 5s 5sr"}} {
		taken, _ := NewPaiWithName(strs[0])
		consumed, _ := StrToPais(strs[1])
		target := 3
		want, _ := NewFuro(Kakan, taken, consumed, &target)
		wantErr := false

		tests = append(tests, testCase{"Kakan", args{Kakan, taken, consumed, &target}, want, wantErr})
	}

	// Kakan valid target: nil
	for _, strs := range [][2]string{{"E", "E E E"}, {"5mr", "5m 5m 5m"}, {"5p", "5pr 5p 5p"}, {"5s", "5s 5sr 5s"}, {"5s", "5s 5s 5sr"}} {
		taken, _ := NewPaiWithName(strs[0])
		consumed, _ := StrToPais(strs[1])
		want, _ := NewFuro(Kakan, taken, consumed, nil)
		wantErr := false

		tests = append(tests, testCase{"Kakan", args{Kakan, taken, consumed, nil}, want, wantErr})
	}

	// Kakan invalid
	{
		taken, _ := NewPaiWithName("1m")
		consumed3, _ := StrToPais("1m 1m 1m")
		consumed2, _ := StrToPais("1m 1m")
		consumed4, _ := StrToPais("1m 1m 1m 1m")
		target := 0

		tests = append(tests, testCase{"Kakan invalid: taken: nil", args{Kakan, nil, consumed3, &target}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: taken: nil, target: nil", args{Kakan, nil, consumed3, nil}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: taken: nil, consumed: 2", args{Kakan, nil, consumed2, &target}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: taken: nil, consumed: 2, target: nil", args{Kakan, nil, consumed2, nil}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: taken: nil, consumed: 4", args{Kakan, nil, consumed4, &target}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: taken: nil, consumed: 4, target: nil", args{Kakan, nil, consumed4, nil}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: consumed: 2", args{Kakan, taken, consumed2, &target}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: consumed: 2, target: nil", args{Kakan, taken, consumed2, nil}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: consumed: 4", args{Kakan, taken, consumed4, &target}, nil, true})
		tests = append(tests, testCase{"Kakan invalid: consumed: 4, target: nil", args{Kakan, taken, consumed4, nil}, nil, true})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFuro(tt.args.t, tt.args.taken, tt.args.consumed, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFuro() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFuro() = %v, want %v", got, tt.want)
			}
		})
	}
}
