package game

import (
	"reflect"
	"testing"
)

func mustStrToPais(t *testing.T, str string) []Pai {
	pais, err := StrToPais(str)
	if err != nil {
		t.Fatalf("failed to parse pais: %v", err)
	}
	return pais
}

func TestNewMentsu(t *testing.T) {
	type args struct {
		t    MentsuType
		pais []Pai
	}
	tests := []struct {
		name    string
		args    args
		want    *Mentsu
		wantErr bool
	}{
		// Shuntsu valid cases
		{"shuntsu_valid_1", args{Shuntsu, mustStrToPais(t, "1m 2m 3m")}, &Mentsu{Shuntsu, mustStrToPais(t, "1m 2m 3m")}, false},
		{"shuntsu_valid_2", args{Shuntsu, mustStrToPais(t, "1m 3m 2m")}, &Mentsu{Shuntsu, mustStrToPais(t, "1m 3m 2m")}, false},
		{"shuntsu_valid_3", args{Shuntsu, mustStrToPais(t, "5mr 4m 6m")}, &Mentsu{Shuntsu, mustStrToPais(t, "5mr 4m 6m")}, false},
		// Shuntsu invalid cases
		{"shuntsu_invalid_1", args{Shuntsu, mustStrToPais(t, "2m 3m")}, nil, true},
		{"shuntsu_invalid_2", args{Shuntsu, mustStrToPais(t, "2m 3m 4m 5m")}, nil, true},

		// Kotsu valid cases
		{"kotsu_valid_1", args{Kotsu, mustStrToPais(t, "E E E")}, &Mentsu{Kotsu, mustStrToPais(t, "E E E")}, false},
		{"kotsu_valid_2", args{Kotsu, mustStrToPais(t, "5mr 5m 5m")}, &Mentsu{Kotsu, mustStrToPais(t, "5mr 5m 5m")}, false},
		{"kotsu_valid_3", args{Kotsu, mustStrToPais(t, "5p 5pr 5p")}, &Mentsu{Kotsu, mustStrToPais(t, "5p 5pr 5p")}, false},
		// Kotsu invalid cases
		{"kotsu_invalid_1", args{Kotsu, mustStrToPais(t, "1m 1m")}, nil, true},
		{"kotsu_invalid_2", args{Kotsu, mustStrToPais(t, "1m 1m 1m 1m")}, nil, true},

		// Kantsu valid cases
		{"kantsu_valid_1", args{Kantsu, mustStrToPais(t, "E E E E")}, &Mentsu{Kantsu, mustStrToPais(t, "E E E E")}, false},
		{"kantsu_valid_2", args{Kantsu, mustStrToPais(t, "5mr 5m 5m 5m")}, &Mentsu{Kantsu, mustStrToPais(t, "5mr 5m 5m 5m")}, false},
		{"kantsu_valid_3", args{Kantsu, mustStrToPais(t, "5s 5s 5sr 5s")}, &Mentsu{Kantsu, mustStrToPais(t, "5s 5s 5sr 5s")}, false},
		// Kantsu invalid cases
		{"kantsu_invalid_1", args{Kantsu, mustStrToPais(t, "1m 1m 1m")}, nil, true},
		{"kantsu_invalid_2", args{Kantsu, mustStrToPais(t, "1m 1m 1m 1m 1m")}, nil, true},

		// Toitsu valid cases
		{"toitsu_valid_1", args{Toitsu, mustStrToPais(t, "E E")}, &Mentsu{Toitsu, mustStrToPais(t, "E E")}, false},
		{"toitsu_valid_2", args{Toitsu, mustStrToPais(t, "5mr 5m")}, &Mentsu{Toitsu, mustStrToPais(t, "5mr 5m")}, false},
		// Toitsu invalid cases
		{"toitsu_invalid_1", args{Toitsu, mustStrToPais(t, "1m")}, nil, true},
		{"toitsu_invalid_2", args{Toitsu, mustStrToPais(t, "1m 1m 1m")}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMentsu(tt.args.t, tt.args.pais)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMentsu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMentsu() = %v, want %v", got, tt.want)
			}
		})
	}
}
