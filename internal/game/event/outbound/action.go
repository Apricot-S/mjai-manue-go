package outbound

type action struct {
	Actor int `validate:"min=0,max=3"`
	Log   string
}

func (a *action) isOutboundEvent() {}
