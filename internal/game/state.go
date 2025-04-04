package game

import (
	"cmp"
	"fmt"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

const (
	numPlayers        = 4
	initScore         = 25_000
	maxNumDoraMarkers = 5
	numInitPipais     = NumIDs*4 - 13*numPlayers - 14
	finalTurn         = numInitPipais / numPlayers
)

var validPrevEventsMap = map[message.Type]map[message.Type]struct{}{
	message.TypeTsumo: {
		message.TypeDahai:         struct{}{},
		message.TypeDaiminkan:     struct{}{},
		message.TypeAnkan:         struct{}{},
		message.TypeKakan:         struct{}{},
		message.TypeDora:          struct{}{},
		message.TypeReachAccepted: struct{}{},
	},
	message.TypeDahai: {
		message.TypeTsumo: struct{}{},
		message.TypeChi:   struct{}{},
		message.TypePon:   struct{}{},
		message.TypeReach: struct{}{},
	},
	message.TypeChi: {
		message.TypeDahai:         struct{}{},
		message.TypeReachAccepted: struct{}{},
	},
	message.TypePon: {
		message.TypeDahai:         struct{}{},
		message.TypeReachAccepted: struct{}{},
	},
	message.TypeDaiminkan: {
		message.TypeDahai:         struct{}{},
		message.TypeReachAccepted: struct{}{},
	},
	message.TypeAnkan: {
		message.TypeTsumo: struct{}{},
	},
	message.TypeKakan: {
		message.TypeTsumo: struct{}{},
	},
	message.TypeDora: {
		message.TypeTsumo: struct{}{},
		message.TypeDahai: struct{}{},
		message.TypeAnkan: struct{}{},
	},
	message.TypeReach: {
		message.TypeTsumo: struct{}{},
	},
	message.TypeReachAccepted: {
		message.TypeDahai: struct{}{},
	},
	message.TypeHora: {
		message.TypeTsumo: struct{}{},
		message.TypeDahai: struct{}{},
		message.TypeKakan: struct{}{},
	},
	message.TypeRyukyoku: {
		message.TypeTsumo:         struct{}{}, // Nine Different Terminals and Honors (九種九牌)
		message.TypeDahai:         struct{}{},
		message.TypeReachAccepted: struct{}{}, // Four-Player Riichi (四人立直)
	},
}

func validateCurrentEvent(current, prev message.Type) error {
	validPrevs, exists := validPrevEventsMap[current]
	if !exists {
		return fmt.Errorf("invalid current event: %s", current)
	}
	if _, ok := validPrevs[prev]; !ok {
		return fmt.Errorf("%s is invalid after %s", current, prev)
	}
	return nil
}

func GetPlayerDistance(p1 *Player, p2 *Player) int {
	return (numPlayers + p1.ID() - p2.ID()) % numPlayers
}

func GetNextKyoku(bakaze *Pai, kyokuNum int) (*Pai, int) {
	if kyokuNum == 4 {
		return bakaze.NextForDora(), 1
	}
	return bakaze, kyokuNum + 1
}

type State interface {
	Players() *[numPlayers]Player
	Bakaze() *Pai
	KyokuNum() int
	Honba() int
	Oya() *Player
	Chicha() *Player
	DoraMarkers() []Pai
	NumPipais() int

	Anpais(player *Player) []Pai
	VisiblePais(player *Player) []Pai
	Doras() []Pai
	Jikaze(player *Player) *Pai
	YakuhaiFan(pai *Pai, player *Player) int
	Turn() int
	RankedPlayers() [numPlayers]Player

	OnStartGame(event *message.StartGame) error
	Update(event jsontext.Value) error
	Print()
}

type StateImpl struct {
	players     [numPlayers]Player
	bakaze      Pai
	kyokuNum    int
	honba       int
	oya         *Player
	chicha      *Player
	doraMarkers []Pai
	numPipais   int

	prevActionType message.Type
	// -1 if prev action is not dahai
	prevDahaiActor    int
	prevDahaiPai      *Pai
	currentActionType message.Type
}

func (s *StateImpl) Players() *[numPlayers]Player {
	return &s.players
}

func (s *StateImpl) Bakaze() *Pai {
	return &s.bakaze
}

func (s *StateImpl) KyokuNum() int {
	return s.kyokuNum
}

func (s *StateImpl) Honba() int {
	return s.honba
}

func (s *StateImpl) Oya() *Player {
	return s.oya
}

func (s *StateImpl) Chicha() *Player {
	return s.chicha
}

func (s *StateImpl) DoraMarkers() []Pai {
	return s.doraMarkers
}

func (s *StateImpl) NumPipais() int {
	return s.numPipais
}

func (s *StateImpl) Anpais(player *Player) []Pai {
	return slices.Concat(player.sutehais, player.extraAnpais)
}

func (s *StateImpl) VisiblePais(player *Player) []Pai {
	visiblePais := []Pai{}

	for _, p := range s.players {
		visiblePais = slices.Concat(visiblePais, p.ho)
		for _, furo := range p.furos {
			visiblePais = slices.Concat(visiblePais, furo.Pais())
		}
	}

	return slices.Concat(visiblePais, s.doraMarkers, player.tehais)
}

func (s *StateImpl) Doras() []Pai {
	doras := make([]Pai, len(s.doraMarkers))
	for i, d := range s.doraMarkers {
		doras[i] = *d.NextForDora()
	}
	return doras
}

func (s *StateImpl) Jikaze(player *Player) *Pai {
	j := 1 + (4+player.id-s.oya.id)%4
	p, _ := NewPaiWithDetail(tsupaiType, uint8(j), false)
	return p
}

func (s *StateImpl) YakuhaiFan(pai *Pai, player *Player) int {
	if !pai.IsTsupai() {
		// Suhai
		return 0
	}

	// Jihai
	n := pai.Number()
	if 5 <= n && n <= 7 {
		// Sangenpai
		return 1
	}

	// Kazehai
	fan := 0
	if pai.HasSameSymbol(&s.bakaze) {
		fan++
	}
	if pai.HasSameSymbol(s.Jikaze(player)) {
		fan++
	}
	return fan
}

func (s *StateImpl) Turn() int {
	return (numInitPipais - s.numPipais) / numPlayers
}

func (s *StateImpl) RankedPlayers() [numPlayers]Player {
	players := s.players
	slices.SortStableFunc(players[:], func(p1, p2 Player) int {
		if c := cmp.Compare(p1.score, p2.score); c != 0 {
			// Sort descending.
			return -c
		}
		// In case of a tie, sort by closest to chicha.
		return GetPlayerDistance(&p1, s.chicha) - GetPlayerDistance(&p2, s.chicha)
	})

	return players
}

func (s *StateImpl) OnStartGame(event *message.StartGame) error {
	if event == nil {
		return fmt.Errorf("start_game message is nil")
	}

	names := []string{"", "", "", ""}
	if event.Names != nil {
		names = slices.Clone(event.Names)
	}
	if len(names) != numPlayers {
		return fmt.Errorf("number of players must be 4, but got %d", len(names))
	}

	east, err := NewPaiWithName("E")
	if err != nil {
		return err
	}

	var players [numPlayers]Player
	for i, name := range names {
		p, err := NewPlayer(i, name, initScore)
		if err != nil {
			return err
		}
		players[i] = *p
	}

	s.players = players
	s.bakaze = *east
	s.kyokuNum = 1
	s.honba = 0
	s.oya = &s.players[0]
	s.chicha = &s.players[0]
	s.doraMarkers = make([]Pai, 0, maxNumDoraMarkers)
	s.numPipais = numInitPipais

	s.prevActionType = ""
	s.prevDahaiActor = -1
	s.prevDahaiPai = nil
	s.currentActionType = ""

	return nil
}

func (s *StateImpl) Update(event jsontext.Value) error {
	s.prevActionType = s.currentActionType

	var msg message.Message
	if err := json.Unmarshal(event, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if err := validateCurrentEvent(msg.Type, s.prevActionType); err != nil {
		return err
	}

	// This is specially handled here because it's not an anpai if the dahai is followed by a hora.
	if msg.Type != message.TypeHora && s.prevActionType == message.TypeDahai {
		for _, p := range s.players {
			if p.ID() != s.prevDahaiActor {
				p.AddExtraAnpais(*s.prevDahaiPai)
			}
		}
	}

	switch msg.Type {
	case message.TypeStartKyoku:
		var e message.StartKyoku
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal start_kyoku: %w", err)
		}
		s.currentActionType = message.TypeStartKyoku
		s.onStartKyoku(&e)
	case message.TypeTsumo:
		var e message.Tsumo
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal tsumo: %w", err)
		}
		s.currentActionType = message.TypeTsumo
		s.onTsumo(&e)
	case message.TypeDahai:
		var e message.Dahai
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal dahai: %w", err)
		}
		s.currentActionType = message.TypeDahai
		s.onDahai(&e)
	case message.TypeChi:
		var e message.Chi
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal chi: %w", err)
		}
		s.currentActionType = message.TypeChi
		s.onChi(&e)
	case message.TypePon:
		var e message.Pon
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal pon: %w", err)
		}
		s.currentActionType = message.TypePon
		s.onPon(&e)
	case message.TypeDaiminkan:
		var e message.Daiminkan
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal daiminkan: %w", err)
		}
		s.currentActionType = message.TypeDaiminkan
		s.onDaiminkan(&e)
	case message.TypeAnkan:
		var e message.Ankan
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal ankan: %w", err)
		}
		s.currentActionType = message.TypeAnkan
		s.onAnkan(&e)
	case message.TypeKakan:
		var e message.Kakan
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal kakan: %w", err)
		}
		s.currentActionType = message.TypeKakan
		s.onKakan(&e)
	case message.TypeDora:
		var e message.Dora
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal dora: %w", err)
		}
		s.currentActionType = message.TypeDora
		s.onDora(&e)
	case message.TypeReach:
		var e message.Reach
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal reach: %w", err)
		}
		s.currentActionType = message.TypeReach
		s.onReach(&e)
	case message.TypeReachAccepted:
		var e message.ReachAccepted
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal reach_accepted: %w", err)
		}
		s.currentActionType = message.TypeReachAccepted
		s.onReachAccepted(&e)
	case message.TypeHora:
		var e message.Hora
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal hora: %w", err)
		}
		s.currentActionType = message.TypeHora
		s.onHora(&e)
	case message.TypeRyukyoku:
		var e message.Ryukyoku
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal ryukyoku: %w", err)
		}
		s.currentActionType = message.TypeRyukyoku
		s.onRyukyoku(&e)
	default:
		return fmt.Errorf("unknown event type: %v", event)
	}

	return nil
}

func (s *StateImpl) Print() {
	for _, p := range s.players {
		fmt.Fprintf(
			os.Stderr,
			`[%d] tehai: %s
       ho: %s

`,
			p.id,
			PaisToStr(p.tehais),
			PaisToStr(p.ho))
	}
}

func (s *StateImpl) onStartKyoku(event *message.StartKyoku) error {
	if event == nil {
		return fmt.Errorf("start_kyoku message is nil")
	}

	bakaze, err := NewPaiWithName(event.Bakaze)
	if err != nil {
		return err
	}
	s.bakaze = *bakaze
	s.kyokuNum = event.Kyoku
	s.honba = event.Honba
	s.oya = &s.players[event.Oya]
	s.doraMarkers = make([]Pai, 0, maxNumDoraMarkers)
	doraMarker, err := NewPaiWithName(event.DoraMarker)
	if err != nil {
		return err
	}
	s.doraMarkers = append(s.doraMarkers, *doraMarker)
	s.numPipais = numInitPipais

	for i := range numPlayers {
		tehais := make([]Pai, initTehaisSize)
		for j := range tehais {
			tehai, err := NewPaiWithName(event.Tehais[i][j])
			if err != nil {
				return err
			}
			tehais[j] = *tehai
		}

		err := s.players[i].onStartKyoku(tehais, nil)
		if err != nil {
			return err
		}

		if event.Scores != nil {
			s.players[i].SetScore(event.Scores[i])
		}
	}

	s.prevActionType = ""
	s.prevDahaiActor = -1
	s.prevDahaiPai = nil
	s.currentActionType = ""

	return nil
}

func (s *StateImpl) onTsumo(event *message.Tsumo) error {
	if event == nil {
		return fmt.Errorf("tsumo message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("tsumo is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}
	s.numPipais--

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	actor := event.Actor
	player := &s.players[actor]
	err = player.onTsumo(*pai)
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onDahai(event *message.Dahai) error {
	if event == nil {
		return fmt.Errorf("dahai message is nil")
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	actor := event.Actor
	player := &s.players[actor]
	err = player.onDahai(*pai)
	if err != nil {
		return err
	}

	s.prevDahaiActor = actor
	s.prevDahaiPai = pai

	return nil
}

func (s *StateImpl) onChi(event *message.Chi) error {
	if event == nil {
		return fmt.Errorf("chi message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("chi is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	var consumed [2]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewChi(*pai, consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onChiPonKan(furo)
	if err != nil {
		return err
	}

	target := event.Target
	err = s.players[target].onTargeted(furo)
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onPon(event *message.Pon) error {
	if event == nil {
		return fmt.Errorf("pon message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("pon is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	var consumed [2]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewPon(*pai, consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onChiPonKan(furo)
	if err != nil {
		return err
	}

	target := event.Target
	err = s.players[target].onTargeted(furo)
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onDaiminkan(event *message.Daiminkan) error {
	if event == nil {
		return fmt.Errorf("daiminkan message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("daiminkan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	var consumed [3]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewDaiminkan(*pai, consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onChiPonKan(furo)
	if err != nil {
		return err
	}

	target := event.Target
	err = s.players[target].onTargeted(furo)
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onAnkan(event *message.Ankan) error {
	if event == nil {
		return fmt.Errorf("ankan message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("ankan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	var consumed [4]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewAnkan(consumed)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onAnkan(furo)
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onKakan(event *message.Kakan) error {
	if event == nil {
		return fmt.Errorf("kakan message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("kakan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	var consumed [3]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewKakan(*pai, consumed, nil)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onKakan(furo)
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onDora(event *message.Dora) error {
	if event == nil {
		return fmt.Errorf("dora message is nil")
	}

	if len(s.doraMarkers) >= maxNumDoraMarkers {
		return fmt.Errorf("a 6th dora cannot be added")
	}

	pai, err := NewPaiWithName(event.DoraMarker)
	if err != nil {
		return err
	}
	s.doraMarkers = append(s.doraMarkers, *pai)

	return nil
}

func (s *StateImpl) onReach(event *message.Reach) error {
	if event == nil {
		return fmt.Errorf("reach message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("reach is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	actor := event.Actor
	player := &s.players[actor]
	err := player.onReach()
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onReachAccepted(event *message.ReachAccepted) error {
	if event == nil {
		return fmt.Errorf("reach_accepted message is nil")
	}

	actor := event.Actor
	player := &s.players[actor]
	var err error
	if event.Scores != nil {
		err = player.onReachAccepted(&event.Scores[actor])
	} else {
		err = player.onReachAccepted(nil)
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onHora(event *message.Hora) error {
	if event == nil {
		return fmt.Errorf("hora message is nil")
	}

	if event.Scores != nil {
		for i, score := range event.Scores {
			s.players[i].SetScore(score)
		}
	}

	return nil
}

func (s *StateImpl) onRyukyoku(event *message.Ryukyoku) error {
	if event == nil {
		return fmt.Errorf("ryukyoku message is nil")
	}

	if event.Scores != nil {
		for i, score := range event.Scores {
			s.players[i].SetScore(score)
		}
	}

	return nil
}
