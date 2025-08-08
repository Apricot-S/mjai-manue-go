package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/go-json-experiment/json"
)

func TestNewEndKyoku(t *testing.T) {
	tests := []struct {
		name    string
		want    *EndKyoku
		wantErr bool
	}{
		{
			name: "test NewEndKyoku()",
			want: &EndKyoku{
				Message: Message{Type: TypeEndKyoku},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEndKyoku()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndKyoku() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndKyoku_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *EndKyoku
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			args: &EndKyoku{
				Message: Message{Type: TypeEndKyoku},
			},
			want:    `{"type":"end_kyoku"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &EndKyoku{
				Message: Message{Type: ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &EndKyoku{
				Message: Message{Type: TypeHello},
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal error = %v, want %v", err, tt.wantErr)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestEndKyoku_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    EndKyoku
		wantErr bool
	}{
		{
			name: "valid",
			args: `{"type":"end_kyoku"}`,
			want: EndKyoku{
				Message: Message{Type: TypeEndKyoku},
			},
			wantErr: false,
		},
		{
			name: "with metadata",
			args: `{
				"type":"end_kyoku",
				"metadata": {
					"foo": "bar"
				}
			}`,
			want: EndKyoku{
				Message: Message{Type: TypeEndKyoku},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: EndKyoku{
				Message: Message{Type: ""},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello"}`,
			want: EndKyoku{
				Message: Message{Type: TypeHello},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got EndKyoku
			err := json.Unmarshal([]byte(tt.args), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal error = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndKyoku_ToEvent(t *testing.T) {
	type fields struct {
		Message Message
	}
	tests := []struct {
		name   string
		fields fields
		want   *inbound.EndKyoku
	}{
		{
			name: "valid",
			fields: fields{
				Message: Message{TypeEndKyoku},
			},
			want: inbound.NewEndKyoku(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &EndKyoku{
				Message: tt.fields.Message,
			}
			if got := m.ToEvent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EndKyoku.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
