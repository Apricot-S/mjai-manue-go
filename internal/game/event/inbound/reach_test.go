package inbound

import (
	"reflect"
	"testing"
)

func TestNewReach(t *testing.T) {
	type args struct {
		actor int
	}
	tests := []struct {
		name    string
		args    args
		want    *Reach
		wantErr bool
	}{
		{
			name: "valid actor",
			args: args{
				actor: 0,
			},
			want: &Reach{
				Actor: 0,
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
			got, err := NewReach(tt.args.actor)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReach() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReach() = %v, want %v", got, tt.want)
			}
		})
	}
}
