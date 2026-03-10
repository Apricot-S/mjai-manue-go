package event

type RequestAction struct {
}

func NewRequestAction() *RequestAction {
	return &RequestAction{}
}

func (e RequestAction) EventType() string {
	return "request_action"
}
