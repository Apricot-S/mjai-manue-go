package game

import (
	"cmp"
	"fmt"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
)

const (
	numPlayers        = 4
	initScore         = 25_000
	maxNumDoraMarkers = 5
	numInitPipais     = NumIDs*4 - 13*numPlayers - 14
	finalTurn         = numInitPipais / numPlayers
)

func getDistance(p1 *Player, p2 *Player) int {
	return (numPlayers + p1.ID() - p2.ID()) % numPlayers
}

func getNextKyoku(bakaze *Pai, kyokuNum int) (*Pai, int) {
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
	Update(event any) error
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
		return getDistance(&p1, s.chicha) - getDistance(&p2, s.chicha)
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
	s.oya = &players[0]
	s.chicha = &players[0]
	s.doraMarkers = make([]Pai, 0, maxNumDoraMarkers)
	s.numPipais = numInitPipais

	s.prevActionType = ""
	s.prevDahaiActor = -1
	s.prevDahaiPai = nil
	s.currentActionType = ""

	return nil
}

func (s *StateImpl) Update(event any) error {
	s.prevActionType = s.currentActionType

	// This is specially handled here because it's not an anpai if the dahai is followed by a hora.
	if _, isHora := event.(*message.Hora); !isHora && s.prevActionType == message.TypeDahai {
		for _, p := range s.players {
			if p.ID() != s.prevDahaiActor {
				p.AddExtraAnpais(*s.prevDahaiPai)
			}
		}
	}

	switch e := event.(type) {
	case *message.StartKyoku:
		s.currentActionType = message.TypeStartKyoku
		s.onStartKyoku(e)
	case *message.Tsumo:
		s.currentActionType = message.TypeTsumo
		s.onTsumo(e)
	case *message.Dahai:
		s.currentActionType = message.TypeDahai
		s.onDahai(e)
	case *message.Chi:
		s.currentActionType = message.TypeChi
		s.onChi(e)
	case *message.Pon:
		s.currentActionType = message.TypePon
		s.onPon(e)
	case *message.Daiminkan:
		s.currentActionType = message.TypeDaiminkan
		s.onDaiminkan(e)
	case *message.Ankan:
		s.currentActionType = message.TypeAnkan
		s.onAnkan(e)
	case *message.Kakan:
		s.currentActionType = message.TypeKakan
		s.onKakan(e)
	case *message.Reach:
		s.currentActionType = message.TypeReach
		s.onReach(e)
	case *message.ReachAccepted:
		s.currentActionType = message.TypeReachAccepted
		s.onReachAccepted(e)
	case *message.Hora:
		s.currentActionType = message.TypeHora
		s.onHora(e)
	case *message.Ryukyoku:
		s.currentActionType = message.TypeRyukyoku
		s.onRyukyoku(e)
	default:
		return fmt.Errorf("unknown event type: %T", e)
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

		s.players[i].onStartKyoku(tehais, nil)

		if event.Scores != nil {
			s.players[i].SetScore(event.Scores[i])
		}
	}

	return nil
}

func (s *StateImpl) onTsumo(event *message.Tsumo) error {
	if event == nil {
		return fmt.Errorf("tsumo message is nil")
	}

	s.numPipais--
	if s.numPipais < 0 {
		return fmt.Errorf("numPipais is negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	actor := event.Actor
	player := &s.players[actor]
	player.onTsumo(*pai)

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
	player.onDahai(*pai)

	s.prevDahaiActor = actor
	s.prevDahaiPai = pai

	return nil
}

func (s *StateImpl) onChi(event *message.Chi) error {
	if event == nil {
		return fmt.Errorf("chi message is nil")
	}
	panic("unimplemented!")
}

func (s *StateImpl) onPon(event *message.Pon) error {
	if event == nil {
		return fmt.Errorf("pon message is nil")
	}
	panic("unimplemented!")
}

func (s *StateImpl) onDaiminkan(event *message.Daiminkan) error {
	if event == nil {
		return fmt.Errorf("daiminkan message is nil")
	}
	panic("unimplemented!")
}

func (s *StateImpl) onAnkan(event *message.Ankan) error {
	if event == nil {
		return fmt.Errorf("ankan message is nil")
	}
	panic("unimplemented!")
}

func (s *StateImpl) onKakan(event *message.Kakan) error {
	if event == nil {
		return fmt.Errorf("kakan message is nil")
	}
	panic("unimplemented!")
}

func (s *StateImpl) onReach(event *message.Reach) error {
	if event == nil {
		return fmt.Errorf("reach message is nil")
	}
	panic("unimplemented!")
}

func (s *StateImpl) onReachAccepted(event *message.ReachAccepted) error {
	if event == nil {
		return fmt.Errorf("reach_accepted message is nil")
	}
	panic("unimplemented!")
}

func (s *StateImpl) onHora(event *message.Hora) error {
	if event == nil {
		return fmt.Errorf("hora message is nil")
	}
	panic("unimplemented!")
}

func (s *StateImpl) onRyukyoku(event *message.Ryukyoku) error {
	if event == nil {
		return fmt.Errorf("ryukyoku message is nil")
	}
	panic("unimplemented!")
}

// Java version

// private void onTsumo(final Tsumo tsumo) {
// 	this.numPipais--;

// 	final var actor = tsumo.getActor();
// 	var player = this.players.get(actor);
// 	player.onTsumo(new Pai(tsumo.getPai()));

// 	final var analysis = new ShantenAnalysis(new PaiSet(player.getTehais()));
// 	this.tenpais[actor] = analysis.getShanten() <= 0;
// }

// private void onDahai(final Dahai dahai) {
// 	final var actor = dahai.getActor();
// 	final var pai = new Pai(dahai.getPai());
// 	var player = this.players.get(actor);
// 	player.onDahai(pai);

// 	final var analysis = new ShantenAnalysis(new PaiSet(player.getTehais()));
// 	this.tenpais[actor] = analysis.getShanten() <= 0;

// 	this.previousDahaiActor = OptionalInt.of(actor);
// 	this.previousDahaiPai = Optional.of(pai);
// }

// private void onChi(final Chi chi) {
// 	final var actor = chi.getActor();
// 	final var consumed = chi.getConsumed().stream().map(Pai::new).toList();
// 	final var furo =
// 			new Furo(
// 					Furo.Type.CHI,
// 					Optional.of(new Pai(chi.getPai())),
// 					consumed,
// 					OptionalInt.of(chi.getTarget()));
// 	this.onChiPonKan(actor, furo);
// }

// private void onPon(final Pon pon) {
// 	final var actor = pon.getActor();
// 	final var consumed = pon.getConsumed().stream().map(Pai::new).toList();
// 	final var furo =
// 			new Furo(
// 					Furo.Type.PON,
// 					Optional.of(new Pai(pon.getPai())),
// 					consumed,
// 					OptionalInt.of(pon.getTarget()));
// 	this.onChiPonKan(actor, furo);
// }

// private void onDaiminkan(final Daiminkan daiminkan) {
// 	final var actor = daiminkan.getActor();
// 	final var consumed = daiminkan.getConsumed().stream().map(Pai::new).toList();
// 	final var furo =
// 			new Furo(
// 					Furo.Type.DAIMINKAN,
// 					Optional.of(new Pai(daiminkan.getPai())),
// 					consumed,
// 					OptionalInt.of(daiminkan.getTarget()));
// 	this.onChiPonKan(actor, furo);
// }

// private void onChiPonKan(final int actor, final Furo furo) {
// 	var player = this.players.get(actor);
// 	player.onChiPonKan(furo);
// 	this.players.get(furo.getTarget().orElseThrow()).onTargeted(furo);

// 	final var analysis = new ShantenAnalysis(new PaiSet(player.getTehais()));
// 	this.tenpais[actor] = analysis.getShanten() <= 0;
// }

// private void onAnkan(final Ankan ankan) {
// 	final var actor = ankan.getActor();
// 	final var consumed = ankan.getConsumed().stream().map(Pai::new).toList();
// 	final var furo = new Furo(Furo.Type.ANKAN, Optional.empty(), consumed, OptionalInt.empty());
// 	this.players.get(actor).onAnkan(furo);
// }

// private void onKakan(final Kakan kakan) {
// 	final var actor = kakan.getActor();
// 	final var consumed = kakan.getConsumed().stream().map(Pai::new).toList();
// 	final var furo =
// 			new Furo(
// 					Furo.Type.KAKAN,
// 					Optional.of(new Pai(kakan.getPai())),
// 					consumed,
// 					OptionalInt.empty());
// 	this.players.get(actor).onKakan(furo);
// }

// private void onDora(final Dora dora) throws IllegalArgumentException {
// 	if (this.doraMarkers.size() >= Game.MAX_NUM_DORA_MARKER) {
// 		throw new IllegalArgumentException("A 6th dora cannot be added.");
// 	}
// 	this.doraMarkers.add(new Pai(dora.getDoraMarker()));
// }

// private void onReach(final Reach reach) {
// 	this.players.get(reach.getActor()).onReach();
// }

// private void onReachAccepted(final ReachAccepted reachAccepted) {
// 	final var scores = reachAccepted.getScores();
// 	final var score =
// 			scores != null
// 					? OptionalInt.of(scores[reachAccepted.getActor()])
// 					: OptionalInt.empty();
// 	this.players.get(reachAccepted.getActor()).onReachAccepted(score);
// }

// private void onHora(final Hora hora) {
// 	Optional.ofNullable(hora.getScores())
// 			.ifPresent(
// 					scores ->
// 							IntStream.range(0, scores.length)
// 									.forEach(i -> this.players.get(i).setScore(scores[i])));
// }

// private void onRyukyoku(final Ryukyoku ryukyoku) {
// 	Optional.ofNullable(ryukyoku.getScores())
// 			.ifPresent(
// 					scores ->
// 							IntStream.range(0, scores.length)
// 									.forEach(i -> this.players.get(i).setScore(scores[i])));
// }
// }
