package event

type OutboundEventType int

const (
	TypeNone OutboundEventType = iota + 1
	TypeJoin
	TypeTsumoAction
	TypeDahaiAction
	TypeSkipAction
	TypeChiAction
	TypePonAction
	TypeDaiminkanAction
	TypeKakanAction
	TypeAnkanAction
	TypeReachAction
	TypeHoraAction
)
