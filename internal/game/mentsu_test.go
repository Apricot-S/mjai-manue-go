package game

import (
	"reflect"
	"testing"
)

func s2P(t *testing.T, str string) Pai {
	pai, err := NewPaiWithName(str)
	if err != nil {
		t.Fatalf("failed to parse pais: %v", err)
	}
	return *pai
}

func TestNewShuntsu(t *testing.T) {
	type args struct {
		pai1 Pai
		pai2 Pai
		pai3 Pai
	}
	tests := []struct {
		name    string
		args    args
		want    Shuntsu
		wantErr bool
	}{
		{
			name:    "valid shuntsu 1",
			args:    args{s2P(t, "1m"), s2P(t, "2m"), s2P(t, "3m")},
			want:    Shuntsu{s2P(t, "1m"), s2P(t, "2m"), s2P(t, "3m")},
			wantErr: false,
		},
		{
			name:    "valid shuntsu 2",
			args:    args{s2P(t, "1m"), s2P(t, "3m"), s2P(t, "2m")},
			want:    Shuntsu{s2P(t, "1m"), s2P(t, "3m"), s2P(t, "2m")},
			wantErr: false,
		},
		{
			name:    "valid shuntsu 3",
			args:    args{s2P(t, "5mr"), s2P(t, "4m"), s2P(t, "6m")},
			want:    Shuntsu{s2P(t, "5mr"), s2P(t, "4m"), s2P(t, "6m")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShuntsu(tt.args.pai1, tt.args.pai2, tt.args.pai3)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewShuntsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKotsu(t *testing.T) {
	type args struct {
		pai1 Pai
		pai2 Pai
		pai3 Pai
	}
	tests := []struct {
		name    string
		args    args
		want    Kotsu
		wantErr bool
	}{
		{
			name:    "valid kotsu 1",
			args:    args{s2P(t, "E"), s2P(t, "E"), s2P(t, "E")},
			want:    Kotsu{s2P(t, "E"), s2P(t, "E"), s2P(t, "E")},
			wantErr: false,
		},
		{
			name:    "valid kotsu 2",
			args:    args{s2P(t, "5mr"), s2P(t, "5m"), s2P(t, "5m")},
			want:    Kotsu{s2P(t, "5mr"), s2P(t, "5m"), s2P(t, "5m")},
			wantErr: false,
		},
		{
			name:    "valid kotsu 3",
			args:    args{s2P(t, "5p"), s2P(t, "5pr"), s2P(t, "5p")},
			want:    Kotsu{s2P(t, "5p"), s2P(t, "5pr"), s2P(t, "5p")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKotsu(tt.args.pai1, tt.args.pai2, tt.args.pai3)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewKotsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKantsu(t *testing.T) {
	type args struct {
		pai1 Pai
		pai2 Pai
		pai3 Pai
		pai4 Pai
	}
	tests := []struct {
		name    string
		args    args
		want    Kantsu
		wantErr bool
	}{
		{
			name:    "valid kantsu 1",
			args:    args{s2P(t, "E"), s2P(t, "E"), s2P(t, "E"), s2P(t, "E")},
			want:    Kantsu{s2P(t, "E"), s2P(t, "E"), s2P(t, "E"), s2P(t, "E")},
			wantErr: false,
		},
		{
			name:    "valid kantsu 2",
			args:    args{s2P(t, "5mr"), s2P(t, "5m"), s2P(t, "5m"), s2P(t, "5m")},
			want:    Kantsu{s2P(t, "5mr"), s2P(t, "5m"), s2P(t, "5m"), s2P(t, "5m")},
			wantErr: false,
		},
		{
			name:    "valid kantsu 3",
			args:    args{s2P(t, "5s"), s2P(t, "5s"), s2P(t, "5sr"), s2P(t, "5s")},
			want:    Kantsu{s2P(t, "5s"), s2P(t, "5s"), s2P(t, "5sr"), s2P(t, "5s")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKantsu(tt.args.pai1, tt.args.pai2, tt.args.pai3, tt.args.pai4)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewKantsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewToitsu(t *testing.T) {
	type args struct {
		pai1 Pai
		pai2 Pai
	}
	tests := []struct {
		name    string
		args    args
		want    Toitsu
		wantErr bool
	}{
		{
			name:    "valid toitsu 1",
			args:    args{s2P(t, "E"), s2P(t, "E")},
			want:    Toitsu{s2P(t, "E"), s2P(t, "E")},
			wantErr: false,
		},
		{
			name:    "valid toitsu 2",
			args:    args{s2P(t, "5mr"), s2P(t, "5m")},
			want:    Toitsu{s2P(t, "5mr"), s2P(t, "5m")},
			wantErr: false,
		},
		{
			name:    "valid toitsu 3",
			args:    args{s2P(t, "5p"), s2P(t, "5pr")},
			want:    Toitsu{s2P(t, "5p"), s2P(t, "5pr")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewToitsu(tt.args.pai1, tt.args.pai2)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewToitsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuntsu_ToString(t *testing.T) {
	tests := []struct {
		name string
		s    Shuntsu
		want string
	}{
		{
			name: "Shuntsu.ToString()",
			s:    Shuntsu{s2P(t, "1m"), s2P(t, "2m"), s2P(t, "3m")},
			want: "shuntsu: [1m 2m 3m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ToString(); got != tt.want {
				t.Errorf("Shuntsu.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKotsu_ToString(t *testing.T) {
	tests := []struct {
		name string
		k    Kotsu
		want string
	}{
		{
			name: "Kotsu.ToString()",
			k:    Kotsu{s2P(t, "1m"), s2P(t, "1m"), s2P(t, "1m")},
			want: "kotsu: [1m 1m 1m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.ToString(); got != tt.want {
				t.Errorf("Kotsu.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKantsu_ToString(t *testing.T) {
	tests := []struct {
		name string
		k    Kantsu
		want string
	}{
		{
			name: "Kantsu.ToString()",
			k:    Kantsu{s2P(t, "1m"), s2P(t, "1m"), s2P(t, "1m"), s2P(t, "1m")},
			want: "kantsu: [1m 1m 1m 1m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.ToString(); got != tt.want {
				t.Errorf("Kantsu.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToitsu_ToString(t *testing.T) {
	tests := []struct {
		name string
		t    Toitsu
		want string
	}{
		{
			name: "Toitsu.ToString()",
			t:    Toitsu{s2P(t, "1m"), s2P(t, "1m")},
			want: "toitsu: [1m 1m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.ToString(); got != tt.want {
				t.Errorf("Toitsu.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
