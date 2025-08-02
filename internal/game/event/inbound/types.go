package event

type InboundEventType int

const (
	TypeHello InboundEventType = iota + 1
	TypeStartGame
	TypeStartKyoku
	TypeTsumo
	TypeDahai
	TypeChi
	TypePon
	TypeDaiminkan
	TypeKakan
	TypeAnkan
	TypeDora
	TypeReach
	TypeReachAccepted
	TypeHora
	TypeRyukyoku
	TypeEndKyoku
	TypeEndGame
	TypeError
)
