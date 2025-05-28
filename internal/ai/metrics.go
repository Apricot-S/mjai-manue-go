package ai

import (
	"fmt"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/ai/core"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

type goal struct {
	game.Goal
	requiredBitVectors [4]core.BitVector
	furos              []game.Furo
	points             int
}

type metric struct {
	horaProb                   float64
	safeProb                   float64
	hojuProb                   float64
	ryukyokuProb               float64
	othersHoraProb             float64
	averageHoraPoints          float64
	ryukyokuAveragePoints      float64
	horaPointsDist             *core.ProbDist[[]float64]
	immediateScoreChangesDist  *core.ProbDist[[]float64]
	futureScoreChangesDist     *core.ProbDist[[]float64]
	scoreChangesDist           *core.ProbDist[[]float64]
	scoreChangesDistOnHora     *core.ProbDist[[]float64]
	scoreChangesDistOnRyukyoku *core.ProbDist[[]float64]
	expectedPoints             float64
	expectedHoraPoints         float64
	safeExpectedPoints         float64
	unsafeExpectedPoints       float64
	ryukyokuExpectedPoints     float64
	averageRank                float64
	shanten                    int
	red                        bool
}

type metrics map[string]metric

// numTries is the number of Monte Carlo simulations.
const numTries = 1_000
const numTriesFloat = float64(numTries)

func (a *ManueAI) getMetrics(
	state game.StateViewer,
	playerID int,
	dahaiCandidates []game.Pai,
	reachDahaiCandidates []game.Pai,
	forbiddenDahais []game.Pai,
) (pai *game.Pai, isReach bool, err error) {
	// TODO: Implement logic.
	return &dahaiCandidates[len(dahaiCandidates)-1], false, nil
}

// horaProb: P(hora | this dahai doesn't cause hoju)
// averageHoraPoints: Average hora points assuming I hora
// horaPointsDist: Distribution of hora points assuming I hora
// expectedHoraPoints: Expected hora points assuming this dahai doesn't cause hoju
// shanten: Shanten number
func (a *ManueAI) getMetricsInternal(
	state game.StateViewer,
	playerID int,
	tehais game.Pais,
	furos []game.Furo,
	dahaiCandidates []game.Pai,
	reach bool,
) (metrics, error) {
	ps, err := game.NewPaiSetWithPais(tehais)
	if err != nil {
		return nil, err
	}
	shanten, goals, err := game.AnalyzeShantenWithOption(ps, 1, 8)
	if err != nil {
		return nil, err
	}

	safeProbs, err := a.getSafeProbs(state, playerID, dahaiCandidates)
	if err != nil {
		return nil, err
	}
	immediateScoreChangesDists := a.getImmediateScoreChangesDists(state, playerID, dahaiCandidates)
	ms, err := a.getHoraEstimation(state, playerID, dahaiCandidates, shanten, goals, reach)
	if err != nil {
		return nil, err
	}

	tenpaiRyukyokuAveragePoints := a.getRyukyokuAveragePoints(state, playerID, true)
	notenRyukyokuAveragePoints := a.getRyukyokuAveragePoints(state, playerID, false)
	ryukyokuProb := a.getRyukyokuProb(state)
	ryukyokuProbOnMyNoHora := a.getRyukyokuProbOnMyNoHora(state)

	scoreChangesDistOnRyukyokuIfTenpaiNow := a.getScoreChangesDistOnRyukyoku(state, playerID, true)
	scoreChangesDistOnRyukyokuIfNotenNow := a.getScoreChangesDistOnRyukyoku(state, playerID, false)
	scoreChangesDistsOnOtherHora := make([]*core.ProbDist[[]float64], 0, 3)
	for _, p := range state.Players() {
		if p.ID() == playerID {
			continue
		}
		d := a.getRandomHoraScoreChangesDist(state, playerID, &p)
		scoreChangesDistsOnOtherHora = append(scoreChangesDistsOnOtherHora, d)
	}

	for _, pai := range dahaiCandidates {
		var key string
		if pai.IsUnknown() {
			key = "none"
		} else {
			key = pai.ToString()
		}
		m := ms[key]
		m.red = pai.IsRed()
		m.safeProb = safeProbs[key]
		m.hojuProb = 1.0 - m.safeProb
		m.safeExpectedPoints = m.safeProb * m.expectedHoraPoints
		m.unsafeExpectedPoints = -(1.0 - m.safeProb) * a.stats.AverageHoraPoints
		m.ryukyokuProb = ryukyokuProb
		if m.shanten <= 0 {
			m.ryukyokuAveragePoints = tenpaiRyukyokuAveragePoints
		} else {
			m.ryukyokuAveragePoints = notenRyukyokuAveragePoints
		}
		m.ryukyokuExpectedPoints = m.safeProb * ryukyokuProb * m.ryukyokuAveragePoints

		m.immediateScoreChangesDist = immediateScoreChangesDists[key]
		if m.shanten <= 0 {
			m.scoreChangesDistOnRyukyoku = scoreChangesDistOnRyukyokuIfTenpaiNow
		} else {
			m.scoreChangesDistOnRyukyoku = scoreChangesDistOnRyukyokuIfNotenNow
		}
		m.scoreChangesDistOnHora = a.getScoreChangesDistOnHora(state, playerID, m.horaPointsDist)

		m.ryukyokuProb = (1.0 - m.horaProb) * ryukyokuProbOnMyNoHora
		m.othersHoraProb = (1.0 - m.horaProb) * (1.0 - ryukyokuProbOnMyNoHora)

		myHoraItem := core.WeightedProbDist[[]float64]{Pd: m.scoreChangesDistOnHora, Prob: m.horaProb}
		ryukyokuItem := core.WeightedProbDist[[]float64]{Pd: m.scoreChangesDistOnRyukyoku, Prob: m.ryukyokuProb}
		var otherHoraItems [3]core.WeightedProbDist[[]float64]
		for i, d := range scoreChangesDistsOnOtherHora {
			otherHoraItems[i] = core.WeightedProbDist[[]float64]{Pd: d, Prob: m.othersHoraProb / 3.0}
		}
		items := []core.WeightedProbDist[[]float64]{
			myHoraItem,
			ryukyokuItem,
			otherHoraItems[0],
			otherHoraItems[1],
			otherHoraItems[2],
		}

		m.futureScoreChangesDist = core.Merge(items)
		m.scoreChangesDist = m.immediateScoreChangesDist.Replace(a.noChanges[:], m.futureScoreChangesDist)
		m.expectedPoints = m.scoreChangesDist.Expected()[playerID]
		m.averageRank = a.getAverageRank(state, playerID, m.scoreChangesDist)

		ms[key] = m
	}

	return ms, nil
}

func (a *ManueAI) getHoraEstimation(
	state game.StateViewer,
	playerID int,
	dahaiCandidates []game.Pai,
	shanten int,
	goals []game.Goal,
	reach bool,
) (metrics, error) {
	gs := make([]goal, 0, len(goals))
	for _, g := range goals {
		if reach && g.Shanten > 0 {
			continue
		}
		if shanten > 3 && g.Shanten > shanten {
			// If shanten > 3, including goals with extra pais is too slow.
			continue
		}
		requiredBitVectors := core.CountVectorToBitVectors(&g.RequiredVector)
		// TODO: Implement calculateFan()
		points := 1_000
		if points > 0 {
			gg := goal{
				Goal:               g,
				requiredBitVectors: requiredBitVectors,
				furos:              nil,
				points:             points,
			}
			gs = append(gs, gg)
		}
	}
	gs = slices.Clip(gs)
	fmt.Fprintf(os.Stderr, "goals %d\n", len(gs))

	visiblePaiSet, err := game.NewPaiSetWithPais(state.VisiblePais(&state.Players()[playerID]))
	if err != nil {
		return nil, err
	}
	invisiblePaiSet := game.GetAll()
	invisiblePaiSet.RemovePaiSet(visiblePaiSet)
	invisiblePais := invisiblePaiSet.ToPais()

	numTsumos := a.getNumExpectedRemainingTurns(state)
	// Uses a fixed seed to get a reproducable result, and to make the result comparable
	// e.g., with and without reach.
	rng := core.CreateRNG()
	totalHoraVector := [game.NumIDs + 1]int{}
	totalPointsVector := [game.NumIDs + 1]int{}
	totalPointsFreqsVector := [game.NumIDs + 1]map[int]int{}
	totalYakuToFanVector := [game.NumIDs + 1]map[string]int{}
	for range numTries {
		core.ShuffleWall(rng, &invisiblePais)
		tsumoPais := make(game.Pais, numTsumos)
		copy(tsumoPais, invisiblePais[:numTsumos])
		tsumoVector, err := game.NewPaiSetWithPais(tsumoPais)
		if err != nil {
			return nil, err
		}
		tsumoBitVectors := core.CountVectorToBitVectors(tsumoVector)
		horaVector := [game.NumIDs + 1]int{}
		pointsVector := [game.NumIDs + 1]int{}
		yakuToFanVector := [game.NumIDs + 1]map[string]int{}
		for _, g := range gs {
			achieved := true
			for i := range len(tsumoBitVectors) {
				if !g.requiredBitVectors[i].IsSubsetOf(tsumoBitVectors[i]) {
					achieved = false
					break
				}
			}
			if achieved {
				for pid := range game.NumIDs + 1 {
					if pid == game.NumIDs || g.ThrowableVector[pid] > 0 {
						horaVector[pid] = 1
						if g.points > pointsVector[pid] {
							pointsVector[pid] = g.points
							yakuToFanVector[pid] = nil
						}
					}
				}
			}
		}

		for pid := range game.NumIDs + 1 {
			if horaVector[pid] != 1 {
				continue
			}
			totalHoraVector[pid]++
			points := pointsVector[pid]
			totalPointsVector[pid] += points
			if _, ok := totalPointsFreqsVector[pid][points]; !ok {
				totalPointsFreqsVector[pid][points] = 0
			}
			totalPointsFreqsVector[pid][points]++
			for name, fan := range yakuToFanVector[pid] {
				totalYakuToFanVector[pid][name] += fan
			}
		}
	}

	shantenVector := [game.NumIDs + 1]int{}
	for i := range len(shantenVector) {
		shantenVector[i] = game.InfinityShanten
	}
	shantenVector[game.NumIDs] = shanten
	for _, g := range goals {
		for pid := range game.NumIDs {
			if g.ThrowableVector[pid] > 0 && g.Shanten < shantenVector[pid] {
				shantenVector[pid] = g.Shanten
			}
		}
	}

	ms := make(metrics, len(dahaiCandidates))
	for _, pai := range dahaiCandidates {
		pid := min(pai.RemoveRed().ID(), game.NumIDs)
		var key string
		if pai.IsUnknown() {
			key = "none"
		} else {
			key = pai.ToString()
		}

		hm := core.NewHashMap[[]float64]()
		m := metric{
			horaProb:           float64(totalHoraVector[pid]) / numTriesFloat,
			averageHoraPoints:  float64(totalPointsVector[pid]) / float64(totalHoraVector[pid]),
			horaPointsDist:     core.NewProbDist(hm),
			expectedHoraPoints: float64(totalPointsVector[pid]) / numTriesFloat,
			shanten:            shantenVector[pid],
		}
		ms[key] = m
	}

	return ms, nil
}
