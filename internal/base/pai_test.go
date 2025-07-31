package base

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

var allNames = [...]string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
	"E", "S", "W", "N", "P", "F", "C",
	"5mr", "5pr", "5sr",
	"?",
}

var allTypes = [...]rune{
	'm', 'm', 'm', 'm', 'm', 'm', 'm', 'm', 'm',
	'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
	's', 's', 's', 's', 's', 's', 's', 's', 's',
	't', 't', 't', 't', 't', 't', 't',
	'm', 'p', 's',
	'?',
}

var allNumbers = [...]uint8{
	1, 2, 3, 4, 5, 6, 7, 8, 9,
	1, 2, 3, 4, 5, 6, 7, 8, 9,
	1, 2, 3, 4, 5, 6, 7, 8, 9,
	1, 2, 3, 4, 5, 6, 7,
	5, 5, 5,
	10,
}

var allIsReds = [...]bool{
	false, false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false,
	true, true, true,
	false,
}

func getTestCaseValid(i int) struct {
	name    string
	id      uint8
	typ     rune
	number  uint8
	isRed   bool
	want    *Pai
	wantErr bool
} {
	return struct {
		name    string
		id      uint8
		typ     rune
		number  uint8
		isRed   bool
		want    *Pai
		wantErr bool
	}{
		allNames[i],
		uint8(i),
		allTypes[i],
		allNumbers[i],
		allIsReds[i],
		&Pai{uint8(i), allTypes[i], allNumbers[i], allIsReds[i]},
		false,
	}
}

func TestNewPaiWithID_Valid(t *testing.T) {
	var tests = []struct {
		name    string
		id      uint8
		typ     rune
		number  uint8
		isRed   bool
		want    *Pai
		wantErr bool
	}{}

	for i := range len(allNames) {
		tests = append(tests, getTestCaseValid(i))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaiWithID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaiWithID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaiWithID() = %v, want %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.ID(), tt.id) {
				t.Errorf("NewPaiWithID().ID() = %v, id %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.Type(), tt.typ) {
				t.Errorf("NewPaiWithID().Type() = %v, typ %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.Number(), tt.number) {
				t.Errorf("NewPaiWithID().Number() = %v, number %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.IsRed(), tt.isRed) {
				t.Errorf("NewPaiWithID().IsRed() = %v, isRed %v", got, tt.want)
				return
			}
		})
	}
}

func TestNewPaiWithID_Invalid(t *testing.T) {
	type args struct {
		id uint8
	}
	tests := []struct {
		name    string
		args    args
		want    *Pai
		wantErr bool
	}{
		{"? + 1", args{38}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaiWithID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaiWithID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaiWithID() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestNewPaiWithName_Valid(t *testing.T) {
	var tests = []struct {
		name    string
		id      uint8
		typ     rune
		number  uint8
		isRed   bool
		want    *Pai
		wantErr bool
	}{}

	for i := range len(allNames) {
		tests = append(tests, getTestCaseValid(i))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaiWithName(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaiWithName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaiWithName() = %v, want %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.ID(), tt.id) {
				t.Errorf("NewPaiWithName().ID() = %v, id %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.Type(), tt.typ) {
				t.Errorf("NewPaiWithName().Type() = %v, typ %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.Number(), tt.number) {
				t.Errorf("NewPaiWithName().Number() = %v, number %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.IsRed(), tt.isRed) {
				t.Errorf("NewPaiWithName().IsRed() = %v, isRed %v", got, tt.want)
				return
			}
		})
	}
}

func TestNewPaiWithName_Invalid(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *Pai
		wantErr bool
	}{
		{"", args{""}, nil, true},
		{"1", args{"1"}, nil, true},
		{"m", args{"m"}, nil, true},
		{"0m", args{"0m"}, nil, true},
		{"10p", args{"10p"}, nil, true},
		{"4sr", args{"4sr"}, nil, true},
		{"mp", args{"mp"}, nil, true},
		{"e", args{"e"}, nil, true},
		{"!", args{"!"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaiWithName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaiWithName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaiWithName() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestNewPaiWithDetail_Valid(t *testing.T) {
	var tests = []struct {
		name    string
		id      uint8
		typ     rune
		number  uint8
		isRed   bool
		want    *Pai
		wantErr bool
	}{}

	// Exclude Unknown
	for i := range len(allNames) - 1 {
		tests = append(tests, getTestCaseValid(i))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaiWithDetail(tt.typ, tt.number, tt.isRed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaiWithDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaiWithDetail() = %v, want %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.ID(), tt.id) {
				t.Errorf("NewPaiWithDetail().ID() = %v, id %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.Type(), tt.typ) {
				t.Errorf("NewPaiWithDetail().Type() = %v, typ %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.Number(), tt.number) {
				t.Errorf("NewPaiWithDetail().Number() = %v, number %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got.IsRed(), tt.isRed) {
				t.Errorf("NewPaiWithDetail().IsRed() = %v, isRed %v", got, tt.want)
				return
			}
		})
	}
}

func TestNewPaiWithDetail_Invalid(t *testing.T) {
	type args struct {
		typ    rune
		number uint8
		isRed  bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Pai
		wantErr bool
	}{
		{"?", args{'?', 10, false}, nil, true},
		{"1", args{' ', 1, false}, nil, true},
		{"0m", args{'m', 0, false}, nil, true},
		{"10p", args{'p', 10, false}, nil, true},
		{"8t", args{'t', 8, false}, nil, true},
		{"4sr", args{'s', 4, true}, nil, true},
		{"5tr", args{'t', 5, true}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaiWithDetail(tt.args.typ, tt.args.number, tt.args.isRed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaiWithDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaiWithDetail() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestPai_IsUnknown(t *testing.T) {
	type testCase struct {
		name string
		id   uint8
		want bool
	}
	tests := []testCase{{"?", 37, true}}

	// Exclude Unknown
	for i := range len(allNames) - 1 {
		tests = append(tests, testCase{allNames[i], uint8(i), false})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := NewPaiWithID(tt.id)
			if got := p.IsUnknown(); got != tt.want {
				t.Errorf("Pai.IsUnknown() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestPai_HasSameSymbol(t *testing.T) {
	type testCase struct {
		name   string
		first  *Pai
		second *Pai
		want   bool
	}
	var tests []testCase

	// same Pai has same symbol
	for i, n := range allNames {
		first, _ := NewPaiWithName(n)
		second, _ := NewPaiWithName(n)
		tests = append(tests, testCase{allNames[i], first, second, true})
	}

	// normal 5 suits and red 5 suits has same symbol
	for _, typ := range []rune{'m', 'p', 's'} {
		first, _ := NewPaiWithDetail(typ, 5, false)
		second, _ := NewPaiWithDetail(typ, 5, true)
		tests = append(tests, testCase{fmt.Sprintf("5%c 5%cr", typ, typ), first, second, true})
	}

	// different Pai has not same symbol
	for _, n := range [][]string{{"1m", "2m"}, {"9m", "9p"}, {"E", "S"}, {"1m", "?"}} {
		first, _ := NewPaiWithName(n[0])
		second, _ := NewPaiWithName(n[1])
		tests = append(tests, testCase{fmt.Sprintf("%s %s", n[0], n[1]), first, second, false})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.first
			if got := p.HasSameSymbol(tt.second); got != tt.want {
				t.Errorf("Pai.HasSameSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPai_NextForDora(t *testing.T) {
	type testCase struct {
		name       string
		doraMarker *Pai
		dora       *Pai
		want       bool
	}
	var tests []testCase

	// suits
	for i := range 9 * 3 {
		var nextN int
		if i%9 == 8 {
			nextN = i - 8
		} else {
			nextN = i + 1
		}

		doraMarker, _ := NewPaiWithID(uint8(i))
		dora, _ := NewPaiWithID(uint8(nextN))
		tests = append(tests, testCase{allNames[i], doraMarker, dora, true})
	}

	// winds
	winds := []string{"E", "S", "W", "N"}
	for i := range winds {
		var nextN int
		if i == 3 {
			nextN = 0
		} else {
			nextN = i + 1
		}

		doraMarker, _ := NewPaiWithName(winds[i])
		dora, _ := NewPaiWithName(winds[nextN])
		tests = append(tests, testCase{winds[i], doraMarker, dora, true})
	}

	// dragons
	dragons := []string{"P", "F", "C"}
	for i := range dragons {
		var nextN int
		if i == 2 {
			nextN = 0
		} else {
			nextN = i + 1
		}

		doraMarker, _ := NewPaiWithName(dragons[i])
		dora, _ := NewPaiWithName(dragons[nextN])
		tests = append(tests, testCase{dragons[i], doraMarker, dora, true})
	}

	// red 5 suits
	for _, typ := range []rune{'m', 'p', 's'} {
		doraMarker, _ := NewPaiWithDetail(typ, 5, true)
		dora, _ := NewPaiWithDetail(typ, 6, false)
		tests = append(tests, testCase{fmt.Sprintf("5%cr", typ), doraMarker, dora, true})
	}

	// unknown
	{
		doraMarker, _ := NewPaiWithName("?")
		dora, _ := NewPaiWithName("?")
		tests = append(tests, testCase{"?", doraMarker, dora, true})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.doraMarker
			if got := p.NextForDora(); !reflect.DeepEqual(*tt.dora == *got, tt.want) {
				t.Errorf("Pai.NextForDora() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPai_IsYaochu(t *testing.T) {
	type testCase struct {
		name string
		want bool
	}
	tests := []testCase{
		{"1m", true},
		{"9m", true},
		{"1p", true},
		{"9p", true},
		{"1s", true},
		{"9s", true},
		{"E", true},
		{"S", true},
		{"W", true},
		{"N", true},
		{"P", true},
		{"F", true},
		{"C", true},
	}

	for _, typ := range []rune{'m', 'p', 's'} {
		for _, n := range []uint8{2, 3, 4, 5, 6, 7, 8} {
			p := fmt.Sprintf("%d%c", n, typ)
			tests = append(tests, testCase{p, false})
		}
		p := fmt.Sprintf("5%cr", typ)
		tests = append(tests, testCase{p, false})
	}

	tests = append(tests, testCase{"?", false})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := NewPaiWithName(tt.name)
			if got := p.IsYaochu(); got != tt.want {
				t.Errorf("Pai.IsYaochu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPai_AddRed(t *testing.T) {
	type testCase struct {
		name string
		want *Pai
	}
	var tests []testCase

	for _, n := range allNames {
		p, _ := NewPaiWithName(n)
		if p.Type() != TsupaiType && p.Number() == 5 && !p.IsRed() {
			red, _ := NewPaiWithName(n + "r")
			tests = append(tests, testCase{n, red})
		} else {
			tests = append(tests, testCase{n, p})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := NewPaiWithName(tt.name)
			if got := p.AddRed(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pai.AddRed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPai_RemoveRed(t *testing.T) {
	type testCase struct {
		name string
		want *Pai
	}
	var tests []testCase

	for _, n := range allNames {
		p, _ := NewPaiWithName(n)
		if p.IsRed() {
			normal, _ := NewPaiWithName(n[0:2])
			tests = append(tests, testCase{n, normal})
		} else {
			tests = append(tests, testCase{n, p})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := NewPaiWithName(tt.name)
			if got := p.RemoveRed(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pai.RemoveRed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPai_Next(t *testing.T) {
	type args struct {
		n int8
	}
	type testCase struct {
		name string
		pai  *Pai
		args args
		want *Pai
	}
	var tests []testCase

	// suits
	for _, typ := range []rune{'m', 'p', 's'} {
		for i := 1; i <= 9; i++ {
			p, _ := NewPaiWithDetail(typ, uint8(i), false)

			// plus
			for j := 0; j <= 9-i; j++ {
				n := uint8(i + j)
				next, _ := NewPaiWithDetail(typ, n, false)
				name := fmt.Sprintf("%d%c + %d: %d%c", i, typ, j, n, typ)
				tests = append(tests, testCase{name, p, args{int8(j)}, next})
			}

			// plus out of bound
			for j := 10 - i; j <= 9; j++ {
				name := fmt.Sprintf("%d%c + %d: nil", i, typ, j)
				tests = append(tests, testCase{name, p, args{int8(j)}, nil})
			}

			// minus
			for j := range i {
				n := uint8(i - j)
				prev, _ := NewPaiWithDetail(typ, n, false)
				name := fmt.Sprintf("%d%c - %d: %d%c", i, typ, j, n, typ)
				tests = append(tests, testCase{name, p, args{int8(-j)}, prev})
			}

			// minus out of bound
			for j := i; j <= 9; j++ {
				name := fmt.Sprintf("%d%c - %d: nil", i, typ, j)
				tests = append(tests, testCase{name, p, args{int8(-j)}, nil})
			}
		}

		// red 5
		red, _ := NewPaiWithDetail(typ, uint8(5), true)

		// plus
		for j := 0; j <= 5; j++ {
			n := uint8(5 + j)
			next, _ := NewPaiWithDetail(typ, n, false)
			name := fmt.Sprintf("5%cr + %d: %d%c", typ, j, n, typ)
			tests = append(tests, testCase{name, red, args{int8(j)}, next})
		}

		// plus out of bound
		for j := 5; j <= 9; j++ {
			name := fmt.Sprintf("5%cr + %d: nil", typ, j)
			tests = append(tests, testCase{name, red, args{int8(j)}, nil})
		}

		// minus
		for j := range 5 {
			n := uint8(5 - j)
			prev, _ := NewPaiWithDetail(typ, n, false)
			name := fmt.Sprintf("5%cr - %d: %d%c", typ, j, n, typ)
			tests = append(tests, testCase{name, red, args{int8(-j)}, prev})
		}

		// minus out of bound
		for j := 5; j <= 9; j++ {
			name := fmt.Sprintf("5%cr - %d: nil", typ, j)
			tests = append(tests, testCase{name, red, args{int8(-j)}, nil})
		}
	}

	// honors, unknown
	for _, symbol := range []string{"E", "S", "W", "N", "P", "F", "C", "?"} {
		p, _ := NewPaiWithName(symbol)
		tests = append(tests, testCase{fmt.Sprintf("%s + 0: nil", symbol), p, args{0}, nil})
		tests = append(tests, testCase{fmt.Sprintf("%s + 1: nil", symbol), p, args{1}, nil})
		tests = append(tests, testCase{fmt.Sprintf("%s - 1: nil", symbol), p, args{-1}, nil})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.pai
			if got := p.Next(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pai.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPai_ToString(t *testing.T) {
	type args struct {
		name string
	}
	type testCase struct {
		name string
		args args
		want string
	}
	tests := []testCase{}

	for _, n := range allNames {
		tests = append(tests, testCase{n, args{n}, n})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := NewPaiWithName(tt.args.name)
			if got := p.ToString(); got != tt.want {
				t.Errorf("Pai.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaisToStr(t *testing.T) {
	type args struct {
		pais []Pai
	}
	type testCase struct {
		name string
		args args
		want string
	}
	tests := []testCase{
		{"empty", args{[]Pai{}}, ""},
	}

	// 1 pai
	for _, name := range allNames {
		pai, _ := NewPaiWithName(name)
		pais := []Pai{*pai}
		tests = append(tests, testCase{name, args{pais}, name})
	}

	// 2 pais
	for _, name := range allNames {
		names := fmt.Sprintf("%s %s", name, name)
		pai, _ := NewPaiWithName(name)
		pais := []Pai{*pai, *pai}
		tests = append(tests, testCase{names, args{pais}, names})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PaisToStr(tt.args.pais); got != tt.want {
				t.Errorf("PaisToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToPais(t *testing.T) {
	type args struct {
		names string
	}
	type testCase struct {
		name string
		args args
		want []Pai
	}
	tests := []testCase{
		{"empty", args{""}, []Pai{}},
	}

	// 1 pai
	for _, name := range allNames {
		pai, _ := NewPaiWithName(name)
		pais := []Pai{*pai}
		tests = append(tests, testCase{name, args{name}, pais})
	}

	// 2 pais
	for _, name := range allNames {
		names := fmt.Sprintf("%s %s", name, name)
		pai, _ := NewPaiWithName(name)
		pais := []Pai{*pai, *pai}
		tests = append(tests, testCase{names, args{names}, pais})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := StrToPais(tt.args.names); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrToPais() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortPais(t *testing.T) {
	names := [...]string{
		"?",
		"5sr", "5pr", "5mr",
		"C", "F", "P", "N", "W", "S", "E",
		"9s", "8s", "7s", "6s", "5s", "4s", "3s", "2s", "1s",
		"9p", "8p", "7p", "6p", "5p", "4p", "3p", "2p", "1p",
		"9m", "8m", "7m", "6m", "5m", "4m", "3m", "2m", "1m",
	}

	pais := make(Pais, 0, len(names))
	for _, name := range names {
		pai, _ := NewPaiWithName(name)
		pais = append(pais, *pai)
	}
	sort.Sort(pais)

	sortedNames := [...]string{
		"1m", "2m", "3m", "4m", "5m", "5mr", "6m", "7m", "8m", "9m",
		"1p", "2p", "3p", "4p", "5p", "5pr", "6p", "7p", "8p", "9p",
		"1s", "2s", "3s", "4s", "5s", "5sr", "6s", "7s", "8s", "9s",
		"E", "S", "W", "N", "P", "F", "C",
		"?",
	}

	for i, sortedName := range sortedNames {
		if pais[i].ToString() != sortedName {
			t.Errorf("Expected %s but got %s", sortedName, pais[i].ToString())
		}
	}
}

func TestGetUniquePais(t *testing.T) {
	type args struct {
		ps  Pais
		del func(Pai) bool
	}
	type testCase struct {
		name string
		args args
		want Pais
	}
	tests := []testCase{
		{
			name: "nil",
			args: args{
				ps:  nil,
				del: func(p Pai) bool { return false },
			},
			want: nil,
		},
		{
			name: "empty",
			args: args{
				ps:  Pais{},
				del: func(p Pai) bool { return false },
			},
			want: Pais{},
		},
		{
			name: "one element",
			args: args{
				ps: func() Pais {
					p, _ := NewPaiWithName("1m")
					return Pais{*p}
				}(),
				del: func(p Pai) bool { return false },
			},
			want: func() Pais {
				p, _ := NewPaiWithName("1m")
				return Pais{*p}
			}(),
		},
		{
			name: "two elements, no duplicate",
			args: args{
				ps: func() Pais {
					p1, _ := NewPaiWithName("1m")
					p2, _ := NewPaiWithName("2m")
					return Pais{*p1, *p2}
				}(),
				del: func(p Pai) bool { return false },
			},
			want: func() Pais {
				p1, _ := NewPaiWithName("1m")
				p2, _ := NewPaiWithName("2m")
				return Pais{*p1, *p2}
			}(),
		},
		{
			name: "two elements, duplicate",
			args: args{
				ps: func() Pais {
					p1, _ := NewPaiWithName("1m")
					p2, _ := NewPaiWithName("1m")
					return Pais{*p1, *p2}
				}(),
				del: func(p Pai) bool { return false },
			},
			want: func() Pais {
				p1, _ := NewPaiWithName("1m")
				return Pais{*p1}
			}(),
		},
		{
			name: "two elements, red and black 5m",
			args: args{
				ps: func() Pais {
					p1, _ := NewPaiWithName("5m")
					p2, _ := NewPaiWithName("5mr")
					return Pais{*p1, *p2}
				}(),
				del: func(p Pai) bool { return false },
			},
			want: func() Pais {
				p1, _ := NewPaiWithName("5m")
				p2, _ := NewPaiWithName("5mr")
				return Pais{*p1, *p2}
			}(),
		},
		{
			name: "three elements, 1m 2m 1m",
			args: args{
				ps: func() Pais {
					p1, _ := NewPaiWithName("1m")
					p2, _ := NewPaiWithName("2m")
					p3, _ := NewPaiWithName("1m")
					return Pais{*p1, *p2, *p3}
				}(),
				del: func(p Pai) bool { return false },
			},
			want: func() Pais {
				p1, _ := NewPaiWithName("1m")
				p2, _ := NewPaiWithName("2m")
				return Pais{*p1, *p2}
			}(),
		},
		{
			name: "three elements, all unique",
			args: args{
				ps: func() Pais {
					p1, _ := NewPaiWithName("1m")
					p2, _ := NewPaiWithName("2m")
					p3, _ := NewPaiWithName("3m")
					return Pais{*p1, *p2, *p3}
				}(),
				del: func(p Pai) bool { return false },
			},
			want: func() Pais {
				p1, _ := NewPaiWithName("1m")
				p2, _ := NewPaiWithName("2m")
				p3, _ := NewPaiWithName("3m")
				return Pais{*p1, *p2, *p3}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUniquePais(tt.args.ps, tt.args.del); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUniquePais() = %v, want %v", got, tt.want)
			}
		})
	}
}
