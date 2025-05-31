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

func Test_isIpeko(t *testing.T) {
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
		name := "is ipeko"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("1m 2m 3m")
		pais2, _ := StrToPais("1m 2m 3m")
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
		name := "is not ipeko"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("2m 3m 4m")
		pais2, _ := StrToPais("1m 2m 3m")
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
			if got := isIpeko(tt.args.allMentsus); got != tt.want {
				t.Errorf("isIpeko() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSanshokuDojun(t *testing.T) {
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
		name := "is sanshoku dojun"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("1m 2m 3m")
		pais2, _ := StrToPais("1p 2p 3p")
		pais3, _ := StrToPais("1s 2s 3s")
		pais4, _ := StrToPais("E E E E")
		pais5, _ := StrToPais("N N")
		mentsus[0] = NewShuntsu(pais1[0], pais1[1], pais1[2])
		mentsus[1] = NewShuntsu(pais2[0], pais2[1], pais2[2])
		mentsus[2] = NewShuntsu(pais3[0], pais3[1], pais3[2])
		mentsus[3] = NewKantsu(pais4[0], pais4[1], pais4[2], pais4[3])
		mentsus[4] = NewToitsu(pais5[0], pais5[1])

		tests = append(tests, testCase{
			name: name,
			args: args{allMentsus: mentsus},
			want: true,
		})
	}
	{
		name := "is not sanshoku dojun"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("1m 2m 3m")
		pais2, _ := StrToPais("1s 2s 3s")
		pais3, _ := StrToPais("1s 2s 3s")
		pais4, _ := StrToPais("E E E E")
		pais5, _ := StrToPais("N N")
		mentsus[0] = NewShuntsu(pais1[0], pais1[1], pais1[2])
		mentsus[1] = NewShuntsu(pais2[0], pais2[1], pais2[2])
		mentsus[2] = NewShuntsu(pais3[0], pais3[1], pais3[2])
		mentsus[3] = NewKantsu(pais4[0], pais4[1], pais4[2], pais4[3])
		mentsus[4] = NewToitsu(pais5[0], pais5[1])

		tests = append(tests, testCase{
			name: name,
			args: args{allMentsus: mentsus},
			want: false,
		})
	}
	{
		name := "is not sanshoku dojun for same kotsu"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("1m 1m 1m")
		pais2, _ := StrToPais("1p 1p 1p")
		pais3, _ := StrToPais("1s 1s 1s")
		pais4, _ := StrToPais("W W W")
		pais5, _ := StrToPais("N N")
		mentsus[0] = NewKotsu(pais1[0], pais1[1], pais1[2])
		mentsus[1] = NewKotsu(pais2[0], pais2[1], pais2[2])
		mentsus[2] = NewKotsu(pais3[0], pais3[1], pais3[2])
		mentsus[3] = NewKotsu(pais4[0], pais4[1], pais4[2])
		mentsus[4] = NewToitsu(pais5[0], pais5[1])

		tests = append(tests, testCase{
			name: name,
			args: args{allMentsus: mentsus},
			want: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSanshokuDojun(tt.args.allMentsus); got != tt.want {
				t.Errorf("isSanshokuDojun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isIkkiTsukan(t *testing.T) {
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
		name := "is ikki tsukan"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("1m 2m 3m")
		pais2, _ := StrToPais("4m 5m 6m")
		pais3, _ := StrToPais("7m 8m 9m")
		pais4, _ := StrToPais("E E E E")
		pais5, _ := StrToPais("N N")
		mentsus[0] = NewShuntsu(pais1[0], pais1[1], pais1[2])
		mentsus[1] = NewShuntsu(pais2[0], pais2[1], pais2[2])
		mentsus[2] = NewShuntsu(pais3[0], pais3[1], pais3[2])
		mentsus[3] = NewKantsu(pais4[0], pais4[1], pais4[2], pais4[3])
		mentsus[4] = NewToitsu(pais5[0], pais5[1])

		tests = append(tests, testCase{
			name: name,
			args: args{allMentsus: mentsus},
			want: true,
		})
	}
	{
		name := "is not ikki tsukan"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("1m 2m 3m")
		pais2, _ := StrToPais("4p 5p 6p")
		pais3, _ := StrToPais("7m 8m 9m")
		pais4, _ := StrToPais("E E E E")
		pais5, _ := StrToPais("N N")
		mentsus[0] = NewShuntsu(pais1[0], pais1[1], pais1[2])
		mentsus[1] = NewShuntsu(pais2[0], pais2[1], pais2[2])
		mentsus[2] = NewShuntsu(pais3[0], pais3[1], pais3[2])
		mentsus[3] = NewKantsu(pais4[0], pais4[1], pais4[2], pais4[3])
		mentsus[4] = NewToitsu(pais5[0], pais5[1])

		tests = append(tests, testCase{
			name: name,
			args: args{allMentsus: mentsus},
			want: false,
		})
	}
	{
		name := "is not ikki tsukan for kotsu"
		mentsus := make([]Mentsu, 5)
		pais1, _ := StrToPais("1m 1m 1m")
		pais2, _ := StrToPais("4m 4m 4m")
		pais3, _ := StrToPais("7m 7m 7m")
		pais4, _ := StrToPais("W W W")
		pais5, _ := StrToPais("N N")
		mentsus[0] = NewKotsu(pais1[0], pais1[1], pais1[2])
		mentsus[1] = NewKotsu(pais2[0], pais2[1], pais2[2])
		mentsus[2] = NewKotsu(pais3[0], pais3[1], pais3[2])
		mentsus[3] = NewKotsu(pais4[0], pais4[1], pais4[2])
		mentsus[4] = NewToitsu(pais5[0], pais5[1])

		tests = append(tests, testCase{
			name: name,
			args: args{allMentsus: mentsus},
			want: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIkkiTsukan(tt.args.allMentsus); got != tt.want {
				t.Errorf("isIkkiTsukan() = %v, want %v", got, tt.want)
			}
		})
	}
}
