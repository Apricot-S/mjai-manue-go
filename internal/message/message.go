package message

import "github.com/go-playground/validator/v10"

// Basic message structure.
type Message struct {
	Type Type `json:"type" validate:"required"`
}

// Validator used throughout the program.
var messageValidator = validator.New()
