package inbound

import (
	"encoding/json/v2"
	"fmt"
)

func parseAs[M Message](b []byte) (Message, error) {
	var msg M
	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}

var parseMessageByType = map[string]func([]byte) (Message, error){
	"hello":      parseAs[*Hello],
	"start_game": parseAs[*StartGame],
	"end_game":   parseAs[*EndGame],

	"start_kyoku": parseAs[*StartKyoku],
	"tsumo":       parseAs[*Tsumo],
	"dahai":       parseAs[*Dahai],
}

// ParseMessage decodes a single mjai inbound JSON message into an inbound.Message.
//
// The returned message may or may not be convertible into a domain event.
func ParseMessage(b []byte) (Message, error) {
	var header struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(b, &header); err != nil {
		return nil, err
	}

	parser, ok := parseMessageByType[header.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported message type: %q", header.Type)
	}
	return parser(b)
}
