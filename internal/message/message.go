package message

import "github.com/go-playground/validator/v10"

// Basic message structure.
type Message struct {
	Type Type `json:"type" validate:"required"`
}

// Validator used throughout the program.
var messageValidator = func() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("tile", isValidTile)
	v.RegisterValidation("wind", isValidWind)
	return v
}()
