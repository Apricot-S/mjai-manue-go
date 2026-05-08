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
	maxNumMelds = 4
	// Reference: <https://note.com/daku_longyi/n/n51fe08566f1b>
	maxNumRiver = 24
	// Reference: <https://note.com/daku_longyi/n/n51fe08566f1b>
	maxNumDiscardedTiles = 27
)

type PlayerViewer interface {
	// Hand returns the hand (手牌).
	// It does not include the drawn tile (ツモ牌).
	Hand() (*hand.VisibleHand, bool)
	HandTiles() []tile.Tile
	// DrawnTile returns the drawn tile (ツモ牌).
	// It returns `nil` if the player does not have the drawn tile.
	DrawnTile() *tile.Tile
	// Melds returns the melds (副露).
	Melds() []meld.Meld

	// River returns the River (河).
	// It does not include the tiles that have been called.
	River() []tile.Tile
	// DiscardedTiles returns the discarded tiles (捨て牌).
	// It includes the tiles that have been called.
	DiscardedTiles() []tile.Tile
	// ExtraSafeTiles returns extra safe tiles (安全牌).
	// The tiles that are safe in the same turn and the tiles that are safe after riichi.
	ExtraSafeTiles() []tile.Tile
	// IsFuriten returns whether the player's actual visible hand is furiten.
	// Invisible players cannot update this from actual hand contents.
	IsFuriten() bool
	// CanRonBy returns whether ron by the winning tile is allowed by furiten constraints.
	// It does not validate yaku, score, or complete legal action availability.
	CanRonBy(winningTile *tile.Tile) bool

	// RiichiState returns the riichi state.
	RiichiState() RiichiState
	// RiichiRiverIndex returns the index of the riichi declaration tile in the river.
	//
	// If the riichi declaration tile is called (melded) by another player,
	// the index refers to the next tile in the river.
	// If that tile is also called, the index advances further until it points to
	// the first non-called tile after the riichi declaration.
	//
	// It returns -1 if the player has not declared riichi.
	RiichiRiverIndex() int
	// RiichiDiscardedTilesIndex returns the index of the riichi declaration tile in the discarded tiles.
	// It returns -1 if the player has not declared riichi.
	RiichiDiscardedTilesIndex() int

	// CanDiscard returns whether the player can discard a tile (打牌).
	CanDiscard() bool
	// CanChiiPonKan returns whether the player can chii/pon/called kan.
	CanChiiPonKan() bool
	// IsConcealed returns whether the player hand is concealed (門前).
	IsConcealed() bool
	// SwapCallTiles returns tiles forbidden as immediate discard after a call (喰い替え).
	SwapCallTiles() []tile.Tile
}

type PlayerActor interface {
	Draw(t tile.Tile) error
	Discard(t tile.Tile, tsumogiri bool) error

	Chii(chii meld.Chii) error
	Pon(pon meld.Pon) error
	CalledKan(kan meld.CalledKan) error
	ConcealedKan(kan meld.ConcealedKan) error
	PromotedKan(kan meld.PromotedKan) error

	Riichi() error
	RiichiAccepted() error

	AddExtraSafeTiles(t tile.Tile)
	TakeFromRiver(t tile.Tile) error
}

type Player interface {
	PlayerViewer
	PlayerActor
}
