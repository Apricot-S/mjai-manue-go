package game

import "testing"

func Test_isTanyaochu(t *testing.T) {
	type args struct {
		allPais Pais
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{}

	{
		name := "is tanyaochu"
		pais, _ := StrToPais("2m 2m 6m 7m 8m 2p 3p 4p 3s 3s 3s 5s 6s 7s")
		tests = append(tests, testCase{
			name: name,
			args: args{allPais: pais},
			want: true,
		})
	}
	{
		name := "contains 19"
		pais, _ := StrToPais("2m 2m 6m 7m 8m 1p 2p 3p 3s 3s 3s 5s 6s 7s")
		tests = append(tests, testCase{
			name: name,
			args: args{allPais: pais},
			want: false,
		})
	}
	{
		name := "contains honors"
		pais, _ := StrToPais("E E 6m 7m 8m 2p 3p 4p 3s 3s 3s 5s 6s 7s")
		tests = append(tests, testCase{
			name: name,
			args: args{allPais: pais},
			want: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTanyaochu(tt.args.allPais); got != tt.want {
				t.Errorf("isTanyaochu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isChantaiyao(t *testing.T) {
	type args struct {
		allMentsus []Mentsu
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{}

	{
		name := "is chantaiyao"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("1m 2m 3m")
		pais2, _ := StrToPais("7p 8p 9p")
		pais3, _ := StrToPais("1s 1s 1s")
		pais4, _ := StrToPais("E E E E")
		pais5, _ := StrToPais("N N")
		mentsus[0] = NewShuntsu(pais1[0], pais1[1], pais1[2])
		mentsus[1] = NewShuntsu(pais2[0], pais2[1], pais2[2])
		mentsus[2] = NewKotsu(pais3[0], pais3[1], pais3[2])
		mentsus[3] = NewKantsu(pais4[0], pais4[1], pais4[2], pais4[3])
		mentsus[4] = NewToitsu(pais5[0], pais5[1])

		tests = append(tests, testCase{
			name: name,
			args: args{allMentsus: mentsus},
			want: true,
		})
	}
	{
		name := "is not chantaiyao"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("2m 3m 4m")
		pais2, _ := StrToPais("7p 8p 9p")
		pais3, _ := StrToPais("1s 1s 1s")
		pais4, _ := StrToPais("E E E E")
		pais5, _ := StrToPais("N N")
		mentsus[0] = NewShuntsu(pais1[0], pais1[1], pais1[2])
		mentsus[1] = NewShuntsu(pais2[0], pais2[1], pais2[2])
		mentsus[2] = NewKotsu(pais3[0], pais3[1], pais3[2])
		mentsus[3] = NewKantsu(pais4[0], pais4[1], pais4[2], pais4[3])
		mentsus[4] = NewToitsu(pais5[0], pais5[1])

		tests = append(tests, testCase{
			name: name,
			args: args{allMentsus: mentsus},
			want: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isChantaiyao(tt.args.allMentsus); got != tt.want {
				t.Errorf("isChantaiyao() = %v, want %v", got, tt.want)
			}
		})
	}
}
