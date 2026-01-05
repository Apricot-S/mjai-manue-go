package player

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type RiichiState int

const (
	NotRiichi RiichiState = iota + 1
	RiichiDeclared
	RiichiAccepted
)

const (
	initHandSize = 13
	maxNumMelds  = 4
	// Reference: <https://note.com/daku_longyi/n/n51fe08566f1b>
	maxNumRiver = 24
	// Reference: <https://note.com/daku_longyi/n/n51fe08566f1b>
	maxNumSutehais = 27
)

type PlayerViewer interface {
	// Hand (手牌)
	// It does not include the drawn tile (ツモ牌).
	Hand() (*hand.VisibleHand, bool)
	HandTiles() []tile.Tile
	// Drawn tile (ツモ牌)
	// It is `nil` if the player does not have the drawn tile.
	DrawnTile() *tile.Tile
	// Melds (副露)
	Melds() []meld.Meld
	// River (河)
	// It does not include the tiles that have been called.
	River() []tile.Tile
	// Discarded tiles (捨て牌)
	// It includes the tiles that have been called.
	DiscardedTiles() []tile.Tile
	// Extra safe tiles (安全牌)
	// The tiles that are safe in the same turn and the tiles that are safe after riichi.
	ExtraSafeTiles() []tile.Tile
	// Riichi state
	RiichiState() RiichiState
	// The index of the tile that was declared as riichi in the river.
	// It is -1 if the player has not declared riichi.
	RiichiRiverIndex() int
	// The index of the tile that was declared as riichi in the discarded tiles.
	// It is -1 if the player has not declared riichi.
	RiichiDiscardedTilesIndex() int
	// Whether the player can discard a tile (打牌)
	CanDiscard() bool
	// Whether the player hand is concealed (門前)
	IsConcealed() bool
}

type PlayerActor interface {
	Draw(t tile.Tile) error
	Discard(t tile.Tile, tsumogiri bool) error

	Chii(c meld.Chii) error
	Pon(p meld.Pon) error
	CalledKan(k meld.CalledKan) error
	ConcealedKan(k meld.ConcealedKan) error
	PromotedKan(k meld.PromotedKan) error

	Riichi() error
	RiichiAccepted() error

	AddExtraSafeTiles(t tile.Tile)
	TakeFromRiver(t tile.Tile) error
}

type Player interface {
	PlayerViewer
	PlayerActor
}
