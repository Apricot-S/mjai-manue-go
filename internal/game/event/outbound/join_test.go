package outbound

import (
	"reflect"
	"testing"
)

func TestNewJoin(t *testing.T) {
	type args struct {
		name string
		room string
	}
	tests := []struct {
		name string
		args args
		want *Join
	}{
		{
			name: "without name and room",
			args: args{
				name: "",
				room: "",
			},
			want: &Join{},
		},
		{
			name: "with name and room",
			args: args{
				name: "shanten",
				room: "default",
			},
			want: &Join{
				Name: "shanten",
				Room: "default",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewJoin(tt.args.name, tt.args.room)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}
