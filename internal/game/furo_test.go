package game

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestNewChi(t *testing.T) {
	type args struct {
		taken    Pai
		consumed [2]Pai
		target   int
	}
	type testCase struct {
		name    string
		args    args
		want    *Chi
		wantErr bool
	}
	tests := []testCase{}

	// valid cases
	for _, strs := range [][2]string{{"1m", "2m 3m"}, {"1m", "3m 2m"}, {"5mr", "4m 6m"}} {
		for i := range 4 {
			taken, _ := NewPaiWithName(strs[0])
			consumed, _ := StrToPais(strs[1])
			target := i
			pais := Pais{*taken, consumed[0], consumed[1]}
			sort.Sort(pais)
			want := &Chi{
				taken:    *taken,
				consumed: [2]Pai(consumed),
				target:   target,
				pais:     pais,
			}

			testCase := testCase{
				fmt.Sprintf("chi valid: target %d", target),
				args{*taken, [2]Pai(consumed), target},
				want,
				false,
			}
			tests = append(tests, testCase)
		}
	}

	// invalid cases
	for _, strs := range [][2]string{{"1m", "2m 3m"}, {"1m", "3m 2m"}, {"5mr", "4m 6m"}} {
		for _, i := range [2]int{-1, 4} {
			taken, _ := NewPaiWithName(strs[0])
			consumed, _ := StrToPais(strs[1])
			target := i

			testCase := testCase{
				fmt.Sprintf("chi invalid: target %d", target),
				args{*taken, [2]Pai(consumed), target},
				nil,
				true,
			}
			tests = append(tests, testCase)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChi(tt.args.taken, tt.args.consumed, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPon(t *testing.T) {
	type args struct {
		taken    Pai
		consumed [2]Pai
		target   int
	}
	type testCase struct {
		name    string
		args    args
		want    *Pon
		wantErr bool
	}
	tests := []testCase{}

	// valid cases
	for _, strs := range [][2]string{{"E", "E E"}, {"5mr", "5m 5m"}, {"5p", "5pr 5p"}, {"5s", "5s 5sr"}} {
		for i := range 4 {
			taken, _ := NewPaiWithName(strs[0])
			consumed, _ := StrToPais(strs[1])
			target := i
			pais := Pais{*taken, consumed[0], consumed[1]}
			sort.Sort(pais)
			want := &Pon{
				taken:    *taken,
				consumed: [2]Pai(consumed),
				target:   target,
				pais:     pais,
			}

			testCase := testCase{
				fmt.Sprintf("pon valid: target %d", target),
				args{*taken, [2]Pai(consumed), target},
				want,
				false,
			}
			tests = append(tests, testCase)
		}
	}

	// invalid cases
	for _, strs := range [][2]string{{"E", "E E"}, {"5mr", "5m 5m"}, {"5p", "5pr 5p"}, {"5s", "5s 5sr"}} {
		for _, i := range [2]int{-1, 4} {
			taken, _ := NewPaiWithName(strs[0])
			consumed, _ := StrToPais(strs[1])
			target := i

			testCase := testCase{
				fmt.Sprintf("pon invalid: target %d", target),
				args{*taken, [2]Pai(consumed), target},
				nil,
				true,
			}
			tests = append(tests, testCase)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPon(tt.args.taken, tt.args.consumed, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDaiminkan(t *testing.T) {
	type args struct {
		taken    Pai
		consumed [3]Pai
		target   int
	}
	type testCase struct {
		name    string
		args    args
		want    *Daiminkan
		wantErr bool
	}
	tests := []testCase{}

	// valid cases
	for _, strs := range [][2]string{{"E", "E E E"}, {"5mr", "5m 5m 5m"}, {"5p", "5pr 5p 5p"}, {"5s", "5s 5sr 5s"}, {"5s", "5s 5s 5sr"}} {
		for i := range 4 {
			taken, _ := NewPaiWithName(strs[0])
			consumed, _ := StrToPais(strs[1])
			target := i
			pais := Pais{*taken, consumed[0], consumed[1], consumed[2]}
			sort.Sort(pais)
			want := &Daiminkan{
				taken:    *taken,
				consumed: [3]Pai(consumed),
				target:   target,
				pais:     pais,
			}

			testCase := testCase{
				fmt.Sprintf("daiminkan valid: target %d", target),
				args{*taken, [3]Pai(consumed), target},
				want,
				false,
			}
			tests = append(tests, testCase)
		}
	}

	// invalid cases
	for _, strs := range [][2]string{{"E", "E E E"}, {"5mr", "5m 5m 5m"}, {"5p", "5pr 5p 5p"}, {"5s", "5s 5sr 5s"}, {"5s", "5s 5s 5sr"}} {
		for _, i := range [2]int{-1, 4} {
			taken, _ := NewPaiWithName(strs[0])
			consumed, _ := StrToPais(strs[1])
			target := i

			testCase := testCase{
				fmt.Sprintf("daiminkan invalid: target %d", target),
				args{*taken, [3]Pai(consumed), target},
				nil,
				true,
			}
			tests = append(tests, testCase)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDaiminkan(tt.args.taken, tt.args.consumed, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDaiminkan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDaiminkan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAnkan(t *testing.T) {
	type args struct {
		consumed [4]Pai
	}
	type testCase struct {
		name    string
		args    args
		want    *Ankan
		wantErr bool
	}
	tests := []testCase{}

	// valid cases
	for _, strs := range []string{"E E E E", "5mr 5m 5m 5m", "5p 5pr 5p 5p", "5s 5s 5sr 5s", "5s 5s 5s 5sr"} {
		consumed, _ := StrToPais(strs)
		pais := Pais{consumed[0], consumed[1], consumed[2], consumed[3]}
		sort.Sort(pais)
		want := &Ankan{
			consumed: [4]Pai(consumed),
			pais:     pais,
		}

		testCase := testCase{
			"ankan valid: target",
			args{[4]Pai(consumed)},
			want,
			false,
		}
		tests = append(tests, testCase)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAnkan(tt.args.consumed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnkan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnkan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKakan(t *testing.T) {
	type args struct {
		taken    Pai
		consumed [3]Pai
		target   *int
	}
	type testCase struct {
		name    string
		args    args
		want    *Kakan
		wantErr bool
	}
	tests := []testCase{}

	// valid cases target: present
	for _, strs := range [][2]string{{"E", "E E E"}, {"5mr", "5m 5m 5m"}, {"5p", "5pr 5p 5p"}, {"5s", "5s 5sr 5s"}, {"5s", "5s 5s 5sr"}} {
		for i := range 4 {
			taken, _ := NewPaiWithName(strs[0])
			consumed, _ := StrToPais(strs[1])
			target := i
			pais := Pais{*taken, consumed[0], consumed[1], consumed[2]}
			sort.Sort(pais)
			want := &Kakan{
				taken:    *taken,
				consumed: [3]Pai(consumed),
				target:   &target,
				pais:     pais,
			}

			testCase := testCase{
				fmt.Sprintf("kakan valid: target %d", target),
				args{*taken, [3]Pai(consumed), &target},
				want,
				false,
			}
			tests = append(tests, testCase)
		}
	}

	// valid cases target: nil
	for _, strs := range [][2]string{{"E", "E E E"}, {"5mr", "5m 5m 5m"}, {"5p", "5pr 5p 5p"}, {"5s", "5s 5sr 5s"}, {"5s", "5s 5s 5sr"}} {
		taken, _ := NewPaiWithName(strs[0])
		consumed, _ := StrToPais(strs[1])
		pais := Pais{*taken, consumed[0], consumed[1], consumed[2]}
		sort.Sort(pais)
		want := &Kakan{
			taken:    *taken,
			consumed: [3]Pai(consumed),
			target:   nil,
			pais:     pais,
		}

		testCase := testCase{
			"kakan valid: target nil",
			args{*taken, [3]Pai(consumed), nil},
			want,
			false,
		}
		tests = append(tests, testCase)
	}

	// invalid cases
	for _, strs := range [][2]string{{"E", "E E E"}, {"5mr", "5m 5m 5m"}, {"5p", "5pr 5p 5p"}, {"5s", "5s 5sr 5s"}, {"5s", "5s 5s 5sr"}} {
		for _, i := range [2]int{-1, 4} {
			taken, _ := NewPaiWithName(strs[0])
			consumed, _ := StrToPais(strs[1])
			target := i

			testCase := testCase{
				fmt.Sprintf("kakan invalid: target %d", target),
				args{*taken, [3]Pai(consumed), &target},
				nil,
				true,
			}
			tests = append(tests, testCase)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKakan(tt.args.taken, tt.args.consumed, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKakan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKakan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChi_ToMentsu(t *testing.T) {
	type testCase struct {
		name string
		chi  *Chi
		want *Shuntsu
	}
	tests := []testCase{}

	for _, strs := range [][2]string{{"1m", "2m 3m"}, {"1m", "3m 2m"}, {"5mr", "4m 6m"}} {
		taken, _ := NewPaiWithName(strs[0])
		consumed, _ := StrToPais(strs[1])
		target := 0
		chi, _ := NewChi(*taken, [2]Pai(consumed), target)

		pais := Pais{*taken, consumed[0], consumed[1]}
		sort.Sort(pais)
		want := Shuntsu{pais[0], pais[1], pais[2]}

		testCase := testCase{
			name: fmt.Sprintf("chi %s %s", strs[0], strs[1]),
			chi:  chi,
			want: &want,
		}
		tests = append(tests, testCase)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.chi.ToMentsu()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chi.ToMentsu() = %v, want %v", got, tt.want)
			}
		})
	}
}
