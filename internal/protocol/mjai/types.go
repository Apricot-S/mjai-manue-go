package mjai

// Type represents the message type in Mjai protocol.
type Type string

const (
	TypeNone          Type = "none"
	TypeHello         Type = "hello"
	TypeJoin          Type = "join"
	TypeStartGame     Type = "start_game"
	TypeStartKyoku    Type = "start_kyoku"
	TypeTsumo         Type = "tsumo"
	TypeDahai         Type = "dahai"
	TypeChi           Type = "chi"
	TypePon           Type = "pon"
	TypeDaiminkan     Type = "daiminkan"
	TypeKakan         Type = "kakan"
	TypeAnkan         Type = "ankan"
	TypeDora          Type = "dora"
	TypeReach         Type = "reach"
	TypeReachAccepted Type = "reach_accepted"
	TypeHora          Type = "hora"
	TypeRyukyoku      Type = "ryukyoku"
	TypeEndKyoku      Type = "end_kyoku"
	TypeEndGame       Type = "end_game"
	TypeError         Type = "error"
)
