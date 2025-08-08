package mjai

// Basic message structure.
type Message struct {
	Type Type `json:"type" validate:"required"`
}
