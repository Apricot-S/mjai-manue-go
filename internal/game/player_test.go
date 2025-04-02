package game

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewPlayer(t *testing.T) {
	type args struct {
		id        int
		name      string
		initScore int
	}
	type testCase struct {
		name    string
		args    args
		want    *Player
		wantErr bool
	}
	var tests []testCase

	// valid id
	for i := range 4 {
		tests = append(tests, testCase{
			name: fmt.Sprintf("validID_%d", i),
			args: args{id: i, name: "", initScore: 25_000},
			want: &Player{
				id: i, name: "",
				tehais:     make([]Pai, 0, 14),
				furos:      make([]Furo, 0, 4),
				ho:         make([]Pai, 0, 24),
				sutehais:   make([]Pai, 0, 27),
				reachState: None,
				score:      25_000,
				isMenzen:   true,
			},
			wantErr: false,
		})
	}

	// invalid id
	for _, i := range []int{-1, 4} {
		tests = append(tests, testCase{
			name:    fmt.Sprintf("invalidID_%d", i),
			args:    args{id: i, name: "", initScore: 25_000},
			want:    nil,
			wantErr: true,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPlayer(tt.args.id, tt.args.name, tt.args.initScore)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_onStartKyoku(t *testing.T) {
	type fields struct {
		id        int
		name      string
		initScore int
	}
	type args struct {
		tehais []Pai
		score  *int
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *Player
		wantErr bool
	}
	var tests []testCase

	// valid cases without score
	{
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tests = append(tests, testCase{
			name:   "validNoScore",
			fields: fields{id: 0, name: "", initScore: 25_000},
			args:   args{tehais: tehais, score: nil},
			want: &Player{
				id:       0,
				name:     "",
				furos:    make([]Furo, 0, 4),
				tehais:   tehais,
				score:    25_000,
				isMenzen: true,
			},
			wantErr: false,
		})
	}

	// valid cases with score
	{
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		initScore := 30_000
		tests = append(tests, testCase{
			name:   "validWithScore",
			fields: fields{id: 0, name: "", initScore: 25_000},
			args:   args{tehais: tehais, score: &initScore},
			want: &Player{
				id:       0,
				name:     "",
				furos:    make([]Furo, 0, 4),
				tehais:   tehais,
				score:    initScore,
				isMenzen: true,
			},
			wantErr: false,
		})
	}

	// short hand
	{
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N")
		tests = append(tests, testCase{
			name:    "shortHand",
			fields:  fields{id: 0, name: "", initScore: 25_000},
			args:    args{tehais: tehais, score: nil},
			want:    nil,
			wantErr: true,
		})
	}

	// long hand
	{
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N N")
		tests = append(tests, testCase{
			name:    "longHand",
			fields:  fields{id: 0, name: "", initScore: 25_000},
			args:    args{tehais: tehais, score: nil},
			want:    nil,
			wantErr: true,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := NewPlayer(tt.fields.id, tt.fields.name, tt.fields.initScore)
			if err := p.onStartKyoku(tt.args.tehais, tt.args.score); (err != nil) != tt.wantErr {
				t.Errorf("Player.onStartKyoku() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlayer_onTsumo(t *testing.T) {
	// valid cases
	t.Run("valid", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("4m")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onTsumo(*tsumoPai)
		if err != nil {
			t.Errorf("Player.onTsumo() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            append(tehais, *tsumoPai),
			furos:             make([]Furo, 0, 4),
			ho:                make([]Pai, 0, 24),
			sutehais:          make([]Pai, 0, 27),
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			canDahai:          true,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// invalid: after tsumo
	t.Run("after tsumo", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai1, _ := NewPaiWithName("4m")
		tsumoPai2, _ := NewPaiWithName("5m")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai1)

		err := p.onTsumo(*tsumoPai2)
		if err == nil {
			t.Errorf("Player.onTsumo() error = %v", err)
		}
	})
}

func TestPlayer_onDahai(t *testing.T) {
	// before reach
	t.Run("before reach", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("4m")
		dahai, _ := NewPaiWithName("2m")
		tehaisAfterDahai, _ := StrToPais("1m 3m 4m 6m 7m 8m 1p 2p 3p 6p 8p N N")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		err := p.onDahai(*dahai)

		if err != nil {
			t.Errorf("Player.onDahai() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterDahai,
			furos:             make([]Furo, 0, 4),
			ho:                []Pai{*dahai},
			sutehais:          []Pai{*dahai},
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// after reach
	t.Run("after reach", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai1, _ := NewPaiWithName("E")
		tsumoPai2, _ := NewPaiWithName("4m")
		dahai1, _ := NewPaiWithName("E")
		dahai2, _ := NewPaiWithName("4m")
		tehaisAfterDahai, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		reachHoIndex := 0
		reachSutehaiIndex := 0

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai1)
		p.onReach()
		p.onDahai(*dahai1)
		p.onReachAccepted(nil)
		p.AddExtraAnpais(*dahai1)
		p.onTsumo(*tsumoPai2)

		err := p.onDahai(*dahai2)
		if err != nil {
			t.Errorf("Player.onDahai() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterDahai,
			furos:             make([]Furo, 0, 4),
			ho:                []Pai{*dahai1, *dahai2},
			sutehais:          []Pai{*dahai1, *dahai2},
			extraAnpais:       []Pai{*dahai1},
			reachState:        Accepted,
			reachHoIndex:      &reachHoIndex,
			reachSutehaiIndex: &reachSutehaiIndex,
			score:             24_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// unknown instead not have pai
	t.Run("unknown instead not have pai", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N ?")
		tsumoPai, _ := NewPaiWithName("4m")
		dahai, _ := NewPaiWithName("5m")
		tehaisAfterDahai, _ := StrToPais("1m 2m 3m 4m 6m 7m 8m 1p 2p 3p 6p 8p N")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)

		err := p.onDahai(*dahai)
		if err != nil {
			t.Errorf("Player.onDahai() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterDahai,
			furos:             make([]Furo, 0, 4),
			ho:                []Pai{*dahai},
			sutehais:          []Pai{*dahai},
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// invalid: before tsumo
	t.Run("before tsumo", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		dahai, _ := NewPaiWithName("2m")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onDahai(*dahai)
		if err == nil {
			t.Errorf("Player.onDahai() error = %v", err)
		}
	})

	// invalid: not have pai
	t.Run("not have pai", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("4m")
		dahai, _ := NewPaiWithName("5m")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)

		err := p.onDahai(*dahai)
		if err == nil {
			t.Errorf("Player.onDahai() error = %v", err)
		}
	})
}

func TestPlayer_onChiPonKan(t *testing.T) {
	// on chi
	t.Run("on chi", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		taken, _ := NewPaiWithName("7p")
		consumed, _ := StrToPais("6p 8p")
		target := 3
		furo, _ := NewChi(*taken, [2]Pai(consumed), target)
		tehaisAfterFuro, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p N N")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onChiPonKan(furo)
		if err != nil {
			t.Errorf("Player.onChiPonKan() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterFuro,
			furos:             []Furo{furo},
			ho:                make([]Pai, 0, 24),
			sutehais:          make([]Pai, 0, 27),
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			canDahai:          true,
			isMenzen:          false,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// on pon
	t.Run("on pon", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		taken, _ := NewPaiWithName("N")
		consumed, _ := StrToPais("N N")
		target := 2
		furo, _ := NewPon(*taken, [2]Pai(consumed), target)
		tehaisAfterFuro, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onChiPonKan(furo)
		if err != nil {
			t.Errorf("Player.onChiPonKan() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterFuro,
			furos:             []Furo{furo},
			ho:                make([]Pai, 0, 24),
			sutehais:          make([]Pai, 0, 27),
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			canDahai:          true,
			isMenzen:          false,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// on daiminkan
	t.Run("on daiminkan", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p N N N")
		taken, _ := NewPaiWithName("N")
		consumed, _ := StrToPais("N N N")
		target := 2
		furo, _ := NewDaiminkan(*taken, [3]Pai(consumed), target)
		tehaisAfterFuro, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onChiPonKan(furo)
		if err != nil {
			t.Errorf("Player.onChiPonKan() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterFuro,
			furos:             []Furo{furo},
			ho:                make([]Pai, 0, 24),
			sutehais:          make([]Pai, 0, 27),
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			canDahai:          false,
			isMenzen:          false,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// cannot 5th furo
	t.Run("cannot 5th furo", func(t *testing.T) {
		tehais, _ := StrToPais("2m 3m 1p 1p 2p 3p 1s 1s 1s 2s N N N")
		target := 2

		taken1, _ := NewPaiWithName("N")
		consumed1, _ := StrToPais("N N N")
		furo1, _ := NewDaiminkan(*taken1, [3]Pai(consumed1), target)
		tsumoPai1, _ := NewPaiWithName("4m")
		dahai1, _ := NewPaiWithName("4m")

		taken2, _ := NewPaiWithName("1s")
		consumed2, _ := StrToPais("1s 1s 1s")
		furo2, _ := NewDaiminkan(*taken2, [3]Pai(consumed2), target)
		tsumoPai2, _ := NewPaiWithName("4m")
		dahai2, _ := NewPaiWithName("4m")

		taken3, _ := NewPaiWithName("1p")
		consumed3, _ := StrToPais("1p 1p")
		furo3, _ := NewPon(*taken3, [2]Pai(consumed3), target)
		dahai3, _ := NewPaiWithName("2s")

		taken4, _ := NewPaiWithName("1m")
		consumed4, _ := StrToPais("2m 3m")
		furo4, _ := NewChi(*taken4, [2]Pai(consumed4), target)
		dahai4, _ := NewPaiWithName("3p")

		taken5, _ := NewPaiWithName("4p")
		consumed5, _ := StrToPais("2p 3p")
		furo5, _ := NewChi(*taken5, [2]Pai(consumed5), target)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onChiPonKan(furo1)
		p.onTsumo(*tsumoPai1)
		p.onDahai(*dahai1)
		p.onChiPonKan(furo2)
		p.onTsumo(*tsumoPai2)
		p.onDahai(*dahai2)
		p.onChiPonKan(furo3)
		p.onDahai(*dahai3)
		p.onChiPonKan(furo4)
		p.onDahai(*dahai4)

		err := p.onChiPonKan(furo5)
		if err == nil {
			t.Errorf("Player.onChiPonKan() error = %v", err)
		}
	})

	// after tsumo
	t.Run("after tsumo", func(t *testing.T) {
		tehais, _ := StrToPais("2m 3m 1p 1p 2p 3p 1s 1s 1s 2s N N N")
		tsumoPai, _ := NewPaiWithName("4m")
		target := 2
		taken, _ := NewPaiWithName("N")
		consumed, _ := StrToPais("N N N")
		furo, _ := NewDaiminkan(*taken, [3]Pai(consumed), target)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)

		err := p.onChiPonKan(furo)
		if err == nil {
			t.Errorf("Player.onChiPonKan() error = %v", err)
		}
	})

	// after reach
	t.Run("after reach", func(t *testing.T) {
		tehais, _ := StrToPais("2m 3m 1p 1p 2p 3p 1s 1s 1s 2s N N N")
		tsumoPai, _ := NewPaiWithName("4m")
		target := 2
		taken, _ := NewPaiWithName("N")
		consumed, _ := StrToPais("N N N")
		furo, _ := NewDaiminkan(*taken, [3]Pai(consumed), target)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onReach()
		p.onDahai(*tsumoPai)
		p.onReachAccepted(nil)

		err := p.onChiPonKan(furo)
		if err == nil {
			t.Errorf("Player.onChiPonKan() error = %v", err)
		}
	})

	// cannot ankan
	t.Run("cannot ankan", func(t *testing.T) {
		tehais, _ := StrToPais("2m 3m 1p 1p 2p 3p 1s 1s 1s N N N N")
		consumed, _ := StrToPais("N N N N")
		furo, _ := NewAnkan([4]Pai(consumed))

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onChiPonKan(furo)
		if err == nil {
			t.Errorf("Player.onChiPonKan() error = %v", err)
		}
	})

	// cannot kakan
	t.Run("cannot kakan", func(t *testing.T) {
		tehais, _ := StrToPais("2m 3m 1p 1p 2p 3p 1s 1s 1s 2s N N N")
		taken, _ := NewPaiWithName("N")
		consumed, _ := StrToPais("N N N")
		target := 2
		furo, _ := NewKakan(*taken, [3]Pai(consumed), &target)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onChiPonKan(furo)
		if err == nil {
			t.Errorf("Player.onChiPonKan() error = %v", err)
		}
	})
}

func TestPlayer_onAnkan(t *testing.T) {
	// on ankan
	t.Run("on ankan", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p N N N")
		tsumoPai, _ := NewPaiWithName("N")
		consumed, _ := StrToPais("N N N N")
		furo, _ := NewAnkan([4]Pai(consumed))
		tehaisAfterFuro, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)

		err := p.onAnkan(furo)
		if err != nil {
			t.Errorf("Player.onAnkan() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterFuro,
			furos:             []Furo{furo},
			ho:                make([]Pai, 0, 24),
			sutehais:          make([]Pai, 0, 27),
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// cannot 5th furo
	t.Run("cannot 5th furo", func(t *testing.T) {
		tehais, _ := StrToPais("2m 3m 1p 1p 2p 3p 1s 1s 1s 2s N N N")
		target := 2

		taken1, _ := NewPaiWithName("N")
		consumed1, _ := StrToPais("N N N")
		furo1, _ := NewDaiminkan(*taken1, [3]Pai(consumed1), target)
		tsumoPai1, _ := NewPaiWithName("4m")
		dahai1, _ := NewPaiWithName("4m")

		taken2, _ := NewPaiWithName("1s")
		consumed2, _ := StrToPais("1s 1s 1s")
		furo2, _ := NewDaiminkan(*taken2, [3]Pai(consumed2), target)
		tsumoPai2, _ := NewPaiWithName("4m")
		dahai2, _ := NewPaiWithName("4m")

		taken3, _ := NewPaiWithName("1p")
		consumed3, _ := StrToPais("1p 1p")
		furo3, _ := NewPon(*taken3, [2]Pai(consumed3), target)
		dahai3, _ := NewPaiWithName("2s")

		taken4, _ := NewPaiWithName("1m")
		consumed4, _ := StrToPais("2m 3m")
		furo4, _ := NewChi(*taken4, [2]Pai(consumed4), target)
		dahai4, _ := NewPaiWithName("3p")

		consumed5, _ := StrToPais("E E E E")
		furo5, _ := NewAnkan([4]Pai(consumed5))

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onChiPonKan(furo1)
		p.onTsumo(*tsumoPai1)
		p.onDahai(*dahai1)
		p.onChiPonKan(furo2)
		p.onTsumo(*tsumoPai2)
		p.onDahai(*dahai2)
		p.onChiPonKan(furo3)
		p.onDahai(*dahai3)
		p.onChiPonKan(furo4)
		p.onDahai(*dahai4)

		err := p.onAnkan(furo5)
		if err == nil {
			t.Errorf("Player.onAnkan() error = %v", err)
		}
	})

	// before tsumo
	t.Run("before tsumo", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p N N N N")
		consumed, _ := StrToPais("N N N N")
		furo, _ := NewAnkan([4]Pai(consumed))

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onAnkan(furo)
		if err == nil {
			t.Errorf("Player.onAnkan() error = %v", err)
		}
	})
}

func TestPlayer_onKakan(t *testing.T) {
	// on kakan
	t.Run("on kakan", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p N N N")
		tsumoPai, _ := NewPaiWithName("6p")

		taken1, _ := NewPaiWithName("N")
		consumed1, _ := StrToPais("N N")
		target := 3
		furo1, _ := NewPon(*taken1, [2]Pai(consumed1), target)
		dahai1, _ := NewPaiWithName("6p")

		taken2, _ := NewPaiWithName("N")
		consumed2, _ := StrToPais("N N N")
		furo2, _ := NewKakan(*taken2, [3]Pai(consumed2), nil)
		tehaisAfterFuro, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p")

		kantsu := Kakan{
			taken:    *taken2,
			consumed: [3]Pai(consumed2),
			target:   &target,
			pais:     append(consumed2, *taken2),
		}

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onChiPonKan(furo1)
		p.onDahai(*dahai1)
		p.onTsumo(*tsumoPai)

		err := p.onKakan(furo2)
		if err != nil {
			t.Errorf("Player.onKakan() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterFuro,
			furos:             []Furo{&kantsu},
			ho:                []Pai{*dahai1},
			sutehais:          []Pai{*dahai1},
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			isMenzen:          false,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// 5th furo
	t.Run("5th furo", func(t *testing.T) {
		tehais, _ := StrToPais("2m 3m 1p 1p 2p 3p 1s 1s 1s 2s N N N")
		target := 2

		taken1, _ := NewPaiWithName("N")
		consumed1, _ := StrToPais("N N N")
		furo1, _ := NewDaiminkan(*taken1, [3]Pai(consumed1), target)
		tsumoPai1, _ := NewPaiWithName("4m")
		dahai1, _ := NewPaiWithName("4m")

		taken2, _ := NewPaiWithName("1s")
		consumed2, _ := StrToPais("1s 1s 1s")
		furo2, _ := NewDaiminkan(*taken2, [3]Pai(consumed2), target)
		tsumoPai2, _ := NewPaiWithName("4m")
		dahai2, _ := NewPaiWithName("4m")

		taken3, _ := NewPaiWithName("1p")
		consumed3, _ := StrToPais("1p 1p")
		furo3, _ := NewPon(*taken3, [2]Pai(consumed3), target)
		dahai3, _ := NewPaiWithName("2s")

		taken4, _ := NewPaiWithName("1m")
		consumed4, _ := StrToPais("2m 3m")
		furo4, _ := NewChi(*taken4, [2]Pai(consumed4), target)
		dahai4, _ := NewPaiWithName("3p")

		tsumoPai, _ := NewPaiWithName("1p")
		taken5, _ := NewPaiWithName("1p")
		consumed5, _ := StrToPais("1p 1p 1p")
		furo5, _ := NewKakan(*taken5, [3]Pai(consumed5), nil)
		tehaisAfterFuro, _ := StrToPais("2p")

		kantsu := Kakan{
			taken:    *taken5,
			consumed: [3]Pai(consumed5),
			target:   &target,
			pais:     append(consumed5, *taken5),
		}

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onChiPonKan(furo1)
		p.onTsumo(*tsumoPai1)
		p.onDahai(*dahai1)
		p.onChiPonKan(furo2)
		p.onTsumo(*tsumoPai2)
		p.onDahai(*dahai2)
		p.onChiPonKan(furo3)
		p.onDahai(*dahai3)
		p.onChiPonKan(furo4)
		p.onDahai(*dahai4)
		p.onTsumo(*tsumoPai)

		err := p.onKakan(furo5)
		if err != nil {
			t.Errorf("Player.onKakan() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterFuro,
			furos:             []Furo{furo1, furo2, &kantsu, furo4},
			ho:                []Pai{*dahai1, *dahai2, *dahai3, *dahai4},
			sutehais:          []Pai{*dahai1, *dahai2, *dahai3, *dahai4},
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			canDahai:          false,
			isMenzen:          false,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// before tsumo
	t.Run("before tsumo", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p N N N N")
		target := 2

		taken1, _ := NewPaiWithName("N")
		consumed1, _ := StrToPais("N N")
		furo1, _ := NewPon(*taken1, [2]Pai(consumed1), target)
		dahai1, _ := NewPaiWithName("6m")

		taken2, _ := NewPaiWithName("N")
		consumed2, _ := StrToPais("N N N")
		furo2, _ := NewKakan(*taken2, [3]Pai(consumed2), nil)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onChiPonKan(furo1)
		p.onDahai(*dahai1)

		err := p.onKakan(furo2)
		if err == nil {
			t.Errorf("Player.onKakan() error = %v", err)
		}
	})

	// without pon
	t.Run("without pon", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p N N N")
		tsumoPai, _ := NewPaiWithName("7p")
		taken, _ := NewPaiWithName("N")
		consumed, _ := StrToPais("N N N")
		furo, _ := NewKakan(*taken, [3]Pai(consumed), nil)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)

		err := p.onKakan(furo)
		if err == nil {
			t.Errorf("Player.onKakan() error = %v", err)
		}
	})
}

func TestPlayer_onReach(t *testing.T) {
	// valid case
	t.Run("on reach", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")
		tehaisAfterTsumo, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N E")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)

		err := p.onReach()
		if err != nil {
			t.Errorf("Player.onReach() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehaisAfterTsumo,
			furos:             make([]Furo, 0, 4),
			ho:                make([]Pai, 0, 24),
			sutehais:          make([]Pai, 0, 27),
			extraAnpais:       nil,
			reachState:        Declared,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			canDahai:          true,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// cannot twice
	t.Run("cannot twice", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onReach()

		err := p.onReach()
		if err == nil {
			t.Errorf("Player.onReach() error = %v", err)
		}
	})

	// cannot after reach
	t.Run("cannot after reach", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai1, _ := NewPaiWithName("E")
		tsumoPai2, _ := NewPaiWithName("E")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai1)
		p.onReach()
		p.onDahai(*tsumoPai1)
		p.onReachAccepted(nil)
		p.onTsumo(*tsumoPai2)

		err := p.onReach()
		if err == nil {
			t.Errorf("Player.onReach() error = %v", err)
		}
	})

	// cannot before tsumo
	t.Run("cannot before tsumo", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onReach()
		if err == nil {
			t.Errorf("Player.onReach() error = %v", err)
		}
	})

	// cannot after furo
	t.Run("cannot after furo", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai1, _ := NewPaiWithName("E")

		target := 3
		taken, _ := NewPaiWithName("N")
		consumed, _ := StrToPais("N N")
		furo, _ := NewPon(*taken, [2]Pai(consumed), target)
		dahai, _ := NewPaiWithName("8p")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onChiPonKan(furo)
		p.onDahai(*dahai)
		p.onTsumo(*tsumoPai1)

		err := p.onReach()
		if err == nil {
			t.Errorf("Player.onReach() error = %v", err)
		}
	})
}

func TestPlayer_onReachAccepted(t *testing.T) {
	// on reach accepted
	t.Run("on reach accepted", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")
		dahaiIndex := 0

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onReach()
		p.onDahai(*tsumoPai)

		err := p.onReachAccepted(nil)
		if err != nil {
			t.Errorf("Player.onReachAccepted() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehais,
			furos:             make([]Furo, 0, 4),
			ho:                []Pai{*tsumoPai},
			sutehais:          []Pai{*tsumoPai},
			extraAnpais:       nil,
			reachState:        Accepted,
			reachHoIndex:      &dahaiIndex,
			reachSutehaiIndex: &dahaiIndex,
			score:             24_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// with score
	t.Run("with score", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")
		dahaiIndex := 0
		scoreAfterReach := 23_000

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onReach()
		p.onDahai(*tsumoPai)

		err := p.onReachAccepted(&scoreAfterReach)
		if err != nil {
			t.Errorf("Player.onReachAccepted() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehais,
			furos:             make([]Furo, 0, 4),
			ho:                []Pai{*tsumoPai},
			sutehais:          []Pai{*tsumoPai},
			extraAnpais:       nil,
			reachState:        Accepted,
			reachHoIndex:      &dahaiIndex,
			reachSutehaiIndex: &dahaiIndex,
			score:             23_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// cannot twice
	t.Run("cannot twice", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onReach()
		p.onDahai(*tsumoPai)
		p.onReachAccepted(nil)

		err := p.onReachAccepted(nil)
		if err == nil {
			t.Errorf("Player.onReachAccepted() error = %v", err)
		}
	})

	// cannot before reach
	t.Run("cannot before reach", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)

		err := p.onReachAccepted(nil)
		if err == nil {
			t.Errorf("Player.onReachAccepted() error = %v", err)
		}
	})
}

func TestPlayer_onTargeted(t *testing.T) {
	// on chi
	t.Run("on chi", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("7p")

		taken, _ := NewPaiWithName("7p")
		consumed, _ := StrToPais("6p 8p")
		target := 0
		furo, _ := NewChi(*taken, [2]Pai(consumed), target)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onDahai(*tsumoPai)

		err := p.onTargeted(furo)
		if err != nil {
			t.Errorf("Player.onTargeted() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehais,
			furos:             make([]Furo, 0, 4),
			ho:                []Pai{},
			sutehais:          []Pai{*tsumoPai},
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// on pon
	t.Run("on pon", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")

		taken, _ := NewPaiWithName("E")
		consumed, _ := StrToPais("E E")
		target := 0
		furo, _ := NewPon(*taken, [2]Pai(consumed), target)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onDahai(*tsumoPai)

		err := p.onTargeted(furo)
		if err != nil {
			t.Errorf("Player.onTargeted() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehais,
			furos:             make([]Furo, 0, 4),
			ho:                []Pai{},
			sutehais:          []Pai{*tsumoPai},
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// on daiminkan
	t.Run("on daiminkan", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")

		taken, _ := NewPaiWithName("E")
		consumed, _ := StrToPais("E E E")
		target := 0
		furo, _ := NewDaiminkan(*taken, [3]Pai(consumed), target)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onDahai(*tsumoPai)

		err := p.onTargeted(furo)
		if err != nil {
			t.Errorf("Player.onTargeted() error = %v", err)
		}

		want := &Player{
			id:                0,
			name:              "",
			tehais:            tehais,
			furos:             make([]Furo, 0, 4),
			ho:                []Pai{},
			sutehais:          []Pai{*tsumoPai},
			extraAnpais:       nil,
			reachState:        None,
			reachHoIndex:      nil,
			reachSutehaiIndex: nil,
			score:             25_000,
			isMenzen:          true,
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("Player = %v, want %v", p, want)
		}
	})

	// cannot on ankan
	t.Run("cannot on ankan", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")

		consumed, _ := StrToPais("E E E E")
		furo, _ := NewAnkan([4]Pai(consumed))

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onDahai(*tsumoPai)

		err := p.onTargeted(furo)
		if err == nil {
			t.Errorf("Player.onTargeted() error = %v", err)
		}
	})

	// cannot on kakan
	t.Run("cannot on kakan", func(t *testing.T) {
		tehais, _ := StrToPais("1m 2m 3m 6m 7m 8m 1p 2p 3p 6p 8p N N")
		tsumoPai, _ := NewPaiWithName("E")

		taken, _ := NewPaiWithName("E")
		consumed, _ := StrToPais("E E E")
		furo, _ := NewKakan(*taken, [3]Pai(consumed), nil)

		p, _ := NewPlayer(0, "", 25_000)
		p.onStartKyoku(tehais, nil)
		p.onTsumo(*tsumoPai)
		p.onDahai(*tsumoPai)

		err := p.onTargeted(furo)
		if err == nil {
			t.Errorf("Player.onTargeted() error = %v", err)
		}
	})
}
