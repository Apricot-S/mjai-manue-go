package inbound

// Message is a single mjai protocol message decoded from JSON.
//
// Some messages can be converted into domain events via ToEvent(), but not all
// mjai messages correspond to domain concepts (e.g. hello/start_game/end_game).
//
// The interface is intentionally sealed to keep the set of messages controlled
// by this package.
type Message interface {
	inboundMessage()
}
