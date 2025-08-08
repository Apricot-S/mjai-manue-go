package outbound

import (
	"reflect"
	"testing"
)

func TestNewSkip(t *testing.T) {
	type args struct {
		actor int
		log   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Skip
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor: 0,
				log:   "",
			},
			want: &Skip{
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
			want: &Skip{
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
			got, err := NewSkip(tt.args.actor, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSkip() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}
