package message

type Action struct {
	BaseMessage
	Actor int `json:"actor" validate:"min=0,max=3"`
}
