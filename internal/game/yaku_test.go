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
