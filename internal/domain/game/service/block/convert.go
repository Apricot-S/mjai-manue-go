package block

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/player/meld"
)

func NewBlockFromMeld(m meld.Meld) Block {
	switch m.(type) {
	case *meld.Chii:
		return MustSequence(*m.ToTiles()[0].RemoveRed())
	case *meld.Pon:
		// Red five is sorted after normal, so RemoveRed() is not necessary.
		return MustTriplet(m.ToTiles()[0])
	case *meld.CalledKan, *meld.ConcealedKan, *meld.PromotedKan:
		// Red five is sorted after normal, so RemoveRed() is not necessary.
		return MustQuad(m.ToTiles()[0])
	default:
		panic("unknown meld type")
	}
}
