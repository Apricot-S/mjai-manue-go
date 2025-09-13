package inbound

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/testutil"
)

func mustTehaiPai(tehai ...string) [13]base.Pai {
	tehaiPai := [13]base.Pai{}
	for i, paiStr := range tehai {
		tehaiPai[i] = *testutil.MustPai(paiStr)
	}
	return tehaiPai
}

func TestNewStartKyoku(t *testing.T) {
	type args struct {
		bakaze     base.Pai
		kyoku      int
		honba      int
		kyotaku    int
		oya        int
		doraMarker base.Pai
		scores     *[4]int
		tehais     [4][13]base.Pai
	}
	tests := []struct {
		name    string
		args    args
		want    *StartKyoku
		wantErr bool
	}{
		{
			name: "without scores",
			args: args{
				bakaze:     *testutil.MustPai("E"),
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: *testutil.MustPai("7s"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want: &StartKyoku{
				Bakaze:     *testutil.MustPai("E"),
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: *testutil.MustPai("7s"),
				Scores:     nil,
				Tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: args{
				bakaze:     *testutil.MustPai("S"),
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: *testutil.MustPai("7s"),
				scores:     &[4]int{25000, 25000, 25000, 25000},
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want: &StartKyoku{
				Bakaze:     *testutil.MustPai("S"),
				Kyoku:      1,
				Honba:      0,
				Kyotaku:    0,
				Oya:        0,
				DoraMarker: *testutil.MustPai("7s"),
				Scores:     &[4]int{25000, 25000, 25000, 25000},
				Tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid bakaze",
			args: args{
				bakaze:     *testutil.MustPai("P"),
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: *testutil.MustPai("7s"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid kyoku min",
			args: args{
				bakaze:     *testutil.MustPai("E"),
				kyoku:      0,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: *testutil.MustPai("7s"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid kyoku max",
			args: args{
				bakaze:     *testutil.MustPai("E"),
				kyoku:      5,
				honba:      0,
				kyotaku:    0,
				oya:        0,
				doraMarker: *testutil.MustPai("7s"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid honba",
			args: args{
				bakaze:     *testutil.MustPai("E"),
				kyoku:      1,
				honba:      -1,
				kyotaku:    0,
				oya:        0,
				doraMarker: *testutil.MustPai("7s"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid kyotaku",
			args: args{
				bakaze:     *testutil.MustPai("E"),
				kyoku:      1,
				honba:      0,
				kyotaku:    -1,
				oya:        0,
				doraMarker: *testutil.MustPai("7s"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid oya min",
			args: args{
				bakaze:     *testutil.MustPai("E"),
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        -1,
				doraMarker: *testutil.MustPai("7s"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid oya max",
			args: args{
				bakaze:     *testutil.MustPai("E"),
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        4,
				doraMarker: *testutil.MustPai("7s"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown dora marker",
			args: args{
				bakaze:     *testutil.MustPai("E"),
				kyoku:      1,
				honba:      0,
				kyotaku:    0,
				oya:        3,
				doraMarker: *testutil.MustPai("?"),
				scores:     nil,
				tehais: [4][13]base.Pai{
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("3m", "4m", "3p", "5pr", "7p", "9p", "4s", "4s", "5sr", "7s", "7s", "W", "N"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
					mustTehaiPai("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStartKyoku(tt.args.bakaze, tt.args.kyoku, tt.args.honba, tt.args.kyotaku, tt.args.oya, tt.args.doraMarker, tt.args.scores, tt.args.tehais)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStartKyoku() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStartKyoku() = %v, want %v", got, tt.want)
			}
		})
	}
}
