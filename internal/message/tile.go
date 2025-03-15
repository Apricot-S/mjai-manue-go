package message

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

var tiles = []string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
	"E", "S", "W", "N", "P", "F", "C",
	"5mr", "5pr", "5sr",
	"?",
}

func isValidTile(fl validator.FieldLevel) bool {
	return slices.Contains(tiles, fl.Field().String())
}

var winds = []string{"E", "S", "W", "N"}

func isValidWind(fl validator.FieldLevel) bool {
	return slices.Contains(winds, fl.Field().String())
}
