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
		pais [3]Pai
	}
	tests := []struct {
		name    string
		args    args
		want    Shuntsu
		wantErr bool
	}{
		{
			name:    "valid shuntsu 1",
			args:    args{pais: [3]Pai{s2P(t, "1m"), s2P(t, "2m"), s2P(t, "3m")}},
			want:    Shuntsu([3]Pai{s2P(t, "1m"), s2P(t, "2m"), s2P(t, "3m")}),
			wantErr: false,
		},
		{
			name:    "valid shuntsu 2",
			args:    args{pais: [3]Pai{s2P(t, "1m"), s2P(t, "3m"), s2P(t, "2m")}},
			want:    Shuntsu([3]Pai{s2P(t, "1m"), s2P(t, "3m"), s2P(t, "2m")}),
			wantErr: false,
		},
		{
			name:    "valid shuntsu 3",
			args:    args{pais: [3]Pai{s2P(t, "5mr"), s2P(t, "4m"), s2P(t, "6m")}},
			want:    Shuntsu([3]Pai{s2P(t, "5mr"), s2P(t, "4m"), s2P(t, "6m")}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShuntsu(tt.args.pais)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewShuntsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKotsu(t *testing.T) {
	type args struct {
		pais [3]Pai
	}
	tests := []struct {
		name    string
		args    args
		want    Kotsu
		wantErr bool
	}{
		{
			name:    "valid kotsu 1",
			args:    args{pais: [3]Pai{s2P(t, "E"), s2P(t, "E"), s2P(t, "E")}},
			want:    Kotsu([3]Pai{s2P(t, "E"), s2P(t, "E"), s2P(t, "E")}),
			wantErr: false,
		},
		{
			name:    "valid kotsu 2",
			args:    args{pais: [3]Pai{s2P(t, "5mr"), s2P(t, "5m"), s2P(t, "5m")}},
			want:    Kotsu([3]Pai{s2P(t, "5mr"), s2P(t, "5m"), s2P(t, "5m")}),
			wantErr: false,
		},
		{
			name:    "valid kotsu 3",
			args:    args{pais: [3]Pai{s2P(t, "5p"), s2P(t, "5pr"), s2P(t, "5p")}},
			want:    Kotsu([3]Pai{s2P(t, "5p"), s2P(t, "5pr"), s2P(t, "5p")}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKotsu(tt.args.pais)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewKotsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKantsu(t *testing.T) {
	type args struct {
		pais [4]Pai
	}
	tests := []struct {
		name    string
		args    args
		want    Kantsu
		wantErr bool
	}{
		{
			name:    "valid kantsu 1",
			args:    args{pais: [4]Pai{s2P(t, "E"), s2P(t, "E"), s2P(t, "E"), s2P(t, "E")}},
			want:    Kantsu([4]Pai{s2P(t, "E"), s2P(t, "E"), s2P(t, "E"), s2P(t, "E")}),
			wantErr: false,
		},
		{
			name:    "valid kantsu 2",
			args:    args{pais: [4]Pai{s2P(t, "5mr"), s2P(t, "5m"), s2P(t, "5m"), s2P(t, "5m")}},
			want:    Kantsu([4]Pai{s2P(t, "5mr"), s2P(t, "5m"), s2P(t, "5m"), s2P(t, "5m")}),
			wantErr: false,
		},
		{
			name:    "valid kantsu 3",
			args:    args{pais: [4]Pai{s2P(t, "5s"), s2P(t, "5s"), s2P(t, "5sr"), s2P(t, "5s")}},
			want:    Kantsu([4]Pai{s2P(t, "5s"), s2P(t, "5s"), s2P(t, "5sr"), s2P(t, "5s")}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKantsu(tt.args.pais)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewKantsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewToitsu(t *testing.T) {
	type args struct {
		pais [2]Pai
	}
	tests := []struct {
		name    string
		args    args
		want    Toitsu
		wantErr bool
	}{
		{
			name:    "valid toitsu 1",
			args:    args{pais: [2]Pai{s2P(t, "E"), s2P(t, "E")}},
			want:    Toitsu([2]Pai{s2P(t, "E"), s2P(t, "E")}),
			wantErr: false,
		},
		{
			name:    "valid toitsu 2",
			args:    args{pais: [2]Pai{s2P(t, "5mr"), s2P(t, "5m")}},
			want:    Toitsu([2]Pai{s2P(t, "5mr"), s2P(t, "5m")}),
			wantErr: false,
		},
		{
			name:    "valid toitsu 3",
			args:    args{pais: [2]Pai{s2P(t, "5p"), s2P(t, "5pr")}},
			want:    Toitsu([2]Pai{s2P(t, "5p"), s2P(t, "5pr")}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewToitsu(tt.args.pais)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewToitsu() = %v, want %v", got, tt.want)
			}
		})
	}
}
