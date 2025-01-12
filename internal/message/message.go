package message

type Message struct {
	Type string `json:"type" validate:"required"`
}
