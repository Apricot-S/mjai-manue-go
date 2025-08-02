package outbound

import (
	"github.com/go-playground/validator/v10"
)

// Validator used throughout the program.
var eventValidator = validator.New()
