package outbound

import (
	"reflect"
	"testing"
)

func TestNewReach(t *testing.T) {
	type args struct {
		actor int
		log   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Reach
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor: 0,
				log:   "",
			},
			want: &Reach{
				action: action{
					Actor: 0,
					Log:   "",
				},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor: 3,
				log:   "test",
			},
			want: &Reach{
				action: action{
					Actor: 3,
					Log:   "test",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor: -1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor: 4,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReach(tt.args.actor, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReach() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReach() = %v, want %v", got, tt.want)
			}
		})
	}
}
