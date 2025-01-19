package message

var NoneType = "none"

type None struct {
	Type string `json:"type" validate:"required,eq=none"`
}
