package message

type Action struct {
	Type  string `json:"type" validate:"required"`
	Actor int    `json:"actor" validate:"min=0,max=3"`
}
