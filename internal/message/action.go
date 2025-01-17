package message

type Action struct {
	Message
	Actor int `json:"actor" validate:"min=0,max=3"`
}
