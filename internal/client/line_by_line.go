package client

import (
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/bot"
)

type LineByLine struct {
	reader io.Reader
	writer io.Writer
	bot    bot.Bot
}
