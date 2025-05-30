package ai

import (
	"fmt"
	"os"
	"slices"
	"strings"

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
	horaPointsDist             *core.ProbDist[float64]
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
) (metrics, error) {
	player := state.Players()[playerID]
	tehais := player.Tehais()
	furos := player.Furos()
	ms := make(metrics)

	canReach := len(dahaiCandidates) != 0 && len(reachDahaiCandidates) != 0
	if canReach {
		nowMetrics, err := a.getMetricsInternal(state, playerID, tehais, furos, reachDahaiCandidates, true)
		if err != nil {
			return nil, err
		}
		ms = mergeMetrics(ms, 0, nowMetrics)

		neverMetrics, err := a.getMetricsInternal(state, playerID, tehais, furos, dahaiCandidates, false)
		if err != nil {
			return nil, err
		}
		ms = mergeMetrics(ms, -1, neverMetrics)
		return ms, nil
	}

	reachDeclared := player.ReachState() == game.Declared
	if reachDeclared {
		defaultMetrics, err := a.getMetricsInternal(state, playerID, tehais, furos, reachDahaiCandidates, true)
		if err != nil {
			return nil, err
		}
		ms = mergeMetrics(ms, -1, defaultMetrics)
		return ms, nil
	}

	defaultMetrics, err := a.getMetricsInternal(state, playerID, tehais, furos, dahaiCandidates, false)
	if err != nil {
		return nil, err
	}
	ms = mergeMetrics(ms, -1, defaultMetrics)
	return ms, nil
}

func mergeMetrics(ms metrics, prefix int, otherMetrics metrics) metrics {
	for key, metric := range otherMetrics {
		ms[fmt.Sprintf("%d.%s", prefix, key)] = metric
	}
	return ms
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

func (a *ManueAI) chooseBestMetric(ms metrics, preferBlack bool) string {
	bestKey := ""
	var bestMetric metric
	for key, m := range ms {
		if bestKey == "" || a.compareMetric(&m, &bestMetric, preferBlack) < 0 {
			bestKey = key
			bestMetric = m
		}
	}
	return bestKey
}

func (a *ManueAI) compareMetric(lhs, rhs *metric, preferBlack bool) int {
	if lhs.averageRank < rhs.averageRank {
		return -1
	}
	if lhs.averageRank > rhs.averageRank {
		return 1
	}
	if lhs.expectedPoints > rhs.expectedPoints {
		return -1
	}
	if lhs.expectedPoints < rhs.expectedPoints {
		return 1
	}
	if preferBlack {
		if !lhs.red && rhs.red {
			return -1
		}
		if lhs.red && !rhs.red {
			return 1
		}
	}
	return 0
}

func (a *ManueAI) printMetrics(ms metrics) {
	if len(ms) == 0 {
		return
	}

	type keyValue struct {
		key string
		m   metric
	}
	sortedMetrics := make([]keyValue, 0, len(ms))
	for k, m := range ms {
		sortedMetrics = append(sortedMetrics, keyValue{k, m})
	}
	slices.SortFunc(sortedMetrics, func(kv1, kv2 keyValue) int {
		return a.compareMetric(&kv1.m, &kv2.m, true)
	})

	columns := []struct {
		Header string
		Format string
	}{
		{"action", "%s"},
		{"avgRank", "%.4f"},
		{"expPt", "%.0f"},
		{"hojuProb", "%.3f"},
		{"myHoraProb", "%.3f"},
		{"ryukyokuProb", "%.3f"},
		{"otherHoraProb", "%.3f"},
		{"avgHoraPt", "%.0f"},
		{"ryukyokuAvgPt", "%.0f"},
		{"shanten", "%d"},
	}

	table := make([][]string, 0, len(sortedMetrics)+1)
	header := make([]string, len(columns))
	for i, col := range columns {
		header[i] = col.Header
	}
	table = append(table, header)

	for _, kv := range sortedMetrics {
		row := make([]string, len(columns))
		row[0] = fmt.Sprintf(columns[0].Format, kv.key)
		row[1] = fmt.Sprintf(columns[1].Format, kv.m.averageRank)
		row[2] = fmt.Sprintf(columns[2].Format, kv.m.expectedPoints)
		row[3] = fmt.Sprintf(columns[3].Format, kv.m.hojuProb)
		row[4] = fmt.Sprintf(columns[4].Format, kv.m.horaProb)
		row[5] = fmt.Sprintf(columns[5].Format, kv.m.ryukyokuProb)
		row[6] = fmt.Sprintf(columns[6].Format, kv.m.othersHoraProb)
		row[7] = fmt.Sprintf(columns[7].Format, kv.m.averageHoraPoints)
		row[8] = fmt.Sprintf(columns[8].Format, kv.m.ryukyokuAveragePoints)
		row[9] = fmt.Sprintf(columns[9].Format, kv.m.shanten)
		table = append(table, row)
	}

	a.log(formatArraysAsTable(table) + "\n")
}

func formatArraysAsTable(arrays [][]string) string {
	if len(arrays) == 0 || len(arrays[0]) == 0 {
		return ""
	}
	widths := make([]int, len(arrays[0]))
	for _, array := range arrays {
		for i, val := range array {
			widths[i] = max(widths[i], len(val))
		}
	}
	var sb strings.Builder
	for _, array := range arrays {
		sb.WriteString("| ")
		for i, val := range array {
			w := widths[i]
			sb.WriteString(fmt.Sprintf("%*s | ", w, val))
		}
		sb.WriteString("\n")
	}
	return sb.String()
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
	for pid := range game.NumIDs + 1 {
		totalPointsFreqsVector[pid] = make(map[int]int)
	}

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

		for _, g := range gs {
			achieved := true
			for i := range len(tsumoBitVectors) {
				if !g.requiredBitVectors[i].IsSubsetOf(tsumoBitVectors[i]) {
					achieved = false
					break
				}
			}
			if !achieved {
				continue
			}

			for pid := range game.NumIDs + 1 {
				if pid == game.NumIDs || g.ThrowableVector[pid] > 0 {
					horaVector[pid] = 1
					if g.points > pointsVector[pid] {
						pointsVector[pid] = g.points
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

		hm := core.NewHashMap[float64]()
		for points, freq := range totalPointsFreqsVector[pid] {
			hm.Set(float64(points), float64(freq)/float64(totalHoraVector[pid]))
		}
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
