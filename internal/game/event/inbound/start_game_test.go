package inbound

import (
	"reflect"
	"testing"
)

func TestNewStartGame(t *testing.T) {
	type args struct {
		id    int
		names [4]string
	}
	tests := []struct {
		name    string
		args    args
		want    *StartGame
		wantErr bool
	}{
		{
			name: "valid id and names",
			args: args{
				id:    3,
				names: [4]string{"Player0", "Player1", "Player2", "Player3"},
			},
			want: &StartGame{
				ID:    3,
				Names: [4]string{"Player0", "Player1", "Player2", "Player3"},
			},
			wantErr: false,
		},
		{
			name: "invalid id min",
			args: args{
				id:    -1,
				names: [4]string{"Player0", "Player1", "Player2", "Player3"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid id max",
			args: args{
				id:    4,
				names: [4]string{"Player0", "Player1", "Player2", "Player3"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStartGame(tt.args.id, tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStartGame() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStartGame() = %v, want %v", got, tt.want)
			}
		})
	}
}
