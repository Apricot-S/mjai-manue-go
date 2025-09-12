package ai

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/ai/core"
	"github.com/Apricot-S/mjai-manue-go/internal/ai/estimator"
	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func (a *ManueAI) getScoreChangesDistOnHora(
	state game.StateViewer,
	playerID int,
	horaPointsDist *core.ScalarProbDist,
) *core.VectorProbDist {
	tsumoHoraProb := float64(a.stats.NumTsumoHoras) / float64(a.stats.NumHoras)
	unitDistMap := core.NewHashMap[[4]float64]()

	for _, target := range state.Players() {
		var changes [4]float64
		if target.ID() != playerID {
			changes = [4]float64{0.0, 0.0, 0.0, 0.0}
			changes[playerID] = 1
			changes[target.ID()] = -1
		} else if playerID == state.Oya().ID() {
			changes = [4]float64{-1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0}
			changes[playerID] = 1
		} else {
			changes = [4]float64{-1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0}
			changes[playerID] = 1
			changes[state.Oya().ID()] = -1.0 / 2.0
		}

		var prob float64
		if target.ID() == playerID {
			prob = tsumoHoraProb
		} else {
			prob = (1.0 - tsumoHoraProb) / 3.0
		}
		unitDistMap.Set(changes, prob)
	}

	u := core.NewVectorProbDist(unitDistMap)
	return core.MultScalarVector(horaPointsDist, u)
}

func (a *ManueAI) getRyukyokuAveragePoints(
	state game.StateViewer,
	playerID int,
	selfTenpai bool,
) float64 {
	notenRyukyokuTenpaiProb := a.getNotenRyukyokuTenpaiProb(state)
	var ryukyokuTenpaiProbs [4]float64
	for i := range 4 {
		player := &state.Players()[i]
		var currentTenpaiProb = 0.0
		if player.ID() == playerID {
			if selfTenpai {
				currentTenpaiProb = 1.0
			} else {
				currentTenpaiProb = 0.0
			}
		} else {
			currentTenpaiProb = a.tenpaiProbEstimator.Estimate(player, state)
		}
		ryukyokuTenpaiProbs[i] = currentTenpaiProb + (1.0-currentTenpaiProb)*notenRyukyokuTenpaiProb
	}

	result := 0.0
	for i := range 1 << 4 {
		var tenpais [4]bool
		for j := range 4 {
			tenpais[j] = (i & (1 << j)) != 0
		}
		prob := 1.0
		numTenpais := 0
		for j := range 4 {
			if tenpais[j] {
				prob *= ryukyokuTenpaiProbs[j]
				numTenpais++
			} else {
				prob *= 1.0 - ryukyokuTenpaiProbs[j]
			}
		}
		if prob > 0.0 {
			var points float64
			if tenpais[playerID] {
				if numTenpais == 4 {
					points = 0.0
				} else {
					points = 3000.0 / float64(numTenpais)
				}
			} else {
				if numTenpais == 0 {
					points = 0.0
				} else {
					points = -3000.0 / float64(4-numTenpais)
				}
			}
			result += prob * points
		}
	}
	return result
}

// Distribution of score changes assuming the kyoku ends with ryukyoku.
func (a *ManueAI) getScoreChangesDistOnRyukyoku(
	state game.StateViewer,
	playerID int,
	selfTenpai bool,
) *core.VectorProbDist {
	notenRyukyokuTenpaiProb := a.getNotenRyukyokuTenpaiProb(state)
	hm1 := core.NewHashMap[[4]float64]()
	hm1.Set([4]float64{0.0, 0.0, 0.0, 0.0}, 1.0)
	tenpaisDist := core.NewVectorProbDist(hm1)

	for _, player := range state.Players() {
		var currentTenpaiProb float64
		if player.ID() == playerID {
			if selfTenpai {
				currentTenpaiProb = 1.0
			} else {
				currentTenpaiProb = 0.0
			}
		} else {
			currentTenpaiProb = a.tenpaiProbEstimator.Estimate(&player, state)
		}
		ryukyokuTenpaiProb := currentTenpaiProb + (1.0-currentTenpaiProb)*notenRyukyokuTenpaiProb

		var tenpais [4]float64
		for i := range 4 {
			if player.ID() == playerID {
				tenpais[i] = 1.0
			} else {
				tenpais[i] = 0.0
			}
		}

		hm2 := core.NewHashMap[[4]float64]()
		hm2.Set([4]float64{0.0, 0.0, 0.0, 0.0}, 1.0-ryukyokuTenpaiProb)
		hm2.Set(tenpais, ryukyokuTenpaiProb)
		dist := core.NewVectorProbDist(hm2)
		tenpaisDist = core.AddVectorVector(tenpaisDist, dist)
	}

	return tenpaisDist.MapValueVector(tenpaisToRyukyokuPointsFloat)
}

func tenpaisToRyukyokuPointsFloat(tenpais [4]float64) [4]float64 {
	t := [4]bool{}
	for i := range tenpais {
		t[i] = tenpais[i] != 0.0
	}
	r := game.TenpaisToRyukyokuPoints(t)
	var ret [4]float64
	for i := range r {
		ret[i] = float64(r[i])
	}
	return ret
}

// Probability that the player is tenpai at the end of the kyoku if
// the player is currently noten and the kyoku ends with ryukyoku.
func (a *ManueAI) getNotenRyukyokuTenpaiProb(state game.StateViewer) float64 {
	notenFreq := float64(a.stats.RyukyokuTenpaiStat.Noten)
	tenpaiFreq := 0.0
	t := state.Turn() + 1.0/4.0
	for t <= game.FinalTurn {
		n := strconv.FormatFloat(t, 'f', -1, 64)
		tenpaiFreq += float64(a.stats.RyukyokuTenpaiStat.TenpaiTurnDistribution[n])
		t += 1.0 / 4.0
	}
	return tenpaiFreq / (tenpaiFreq + notenFreq)
}

func (a *ManueAI) getSafeProbs(
	state game.StateViewer,
	playerID int,
	dahaiCandidates []base.Pai,
) (map[string]float64, error) {
	safeProbs := make(map[string]float64, len(dahaiCandidates))
	for _, pai := range dahaiCandidates {
		var key string
		if pai.IsUnknown() {
			key = "none"
		} else {
			key = pai.ToString()
		}
		safeProbs[key] = 1.0
	}

	me := &state.Players()[playerID]
	for _, player := range state.Players() {
		if player.ID() == playerID {
			continue
		}

		scene, err := estimator.NewScene(state, me, &player)
		if err != nil {
			return nil, err
		}
		tenpaiProb := a.tenpaiProbEstimator.Estimate(&player, state)
		for _, pai := range dahaiCandidates {
			if pai.IsUnknown() {
				continue
			}

			isAnpai, err := scene.Evaluate("anpai", &pai)
			if err != nil {
				return nil, err
			}
			if isAnpai {
				continue
			}

			probInfo, err := a.dangerEstimator.EstimateProb(scene, &pai)
			if err != nil {
				return nil, err
			}
			safeProb := 1.0 - tenpaiProb*probInfo.Prob
			safeProbs[pai.ToString()] *= safeProb
		}
	}

	return safeProbs, nil
}

// Distribution of score changes which happen immediately, for each possible dahai.
// i.e., If this dahai causes hoju, score changes due to the hoju. Otherwise [0, 0, 0, 0].
func (a *ManueAI) getImmediateScoreChangesDists(
	state game.StateViewer,
	playerID int,
	dahaiCandidates []base.Pai,
) map[string]*core.VectorProbDist {
	scoreChangesDists := make(map[string]*core.VectorProbDist, len(dahaiCandidates))
	for _, pai := range dahaiCandidates {
		var key string
		if pai.IsUnknown() {
			key = "none"
		} else {
			key = pai.ToString()
		}
		hm := core.NewHashMap[[4]float64]()
		hm.Set(a.noChanges, 1.0)
		scoreChangesDists[key] = core.NewVectorProbDist(hm)
	}

	me := state.Players()[playerID]
	for _, horaPlayer := range state.Players() {
		if horaPlayer.ID() == playerID {
			continue
		}

		scene, err := estimator.NewScene(state, &me, &horaPlayer)
		if err != nil {
			panic(err)
		}

		tenpaiProb := a.tenpaiProbEstimator.Estimate(&horaPlayer, state)

		var horaPointsFreqs map[string]int
		if horaPlayer.ID() == state.Oya().ID() {
			horaPointsFreqs = a.stats.OyaHoraPointsFreqs
		} else {
			horaPointsFreqs = a.stats.KoHoraPointsFreqs
		}
		hm := core.NewHashMap[float64]()
		totalFreqs := float64(horaPointsFreqs["total"])
		for points, freq := range horaPointsFreqs {
			if points == "total" {
				continue
			}
			p, err := strconv.ParseFloat(points, 64)
			if err != nil {
				panic("Invalid stats file: failed to convert key of horaPointsFreqs to float64 (" + points + ").")
			}
			f := float64(freq) / totalFreqs
			hm.Set(p, f)
		}
		horaPointsDist := core.NewScalarProbDist(hm)

		hojuChanges := [4]float64{0.0, 0.0, 0.0, 0.0}
		hojuChanges[horaPlayer.ID()] = 1.0
		hojuChanges[playerID] = -1.0

		for _, pai := range dahaiCandidates {
			if pai.IsUnknown() {
				continue
			}
			key := pai.ToString()
			isAnpai, err := scene.Evaluate("anpai", &pai)
			if err != nil {
				panic(err)
			}

			var hojuProb float64
			if isAnpai {
				hojuProb = 0.0
			} else {
				probInfo, err := a.dangerEstimator.EstimateProb(scene, &pai)
				if err != nil {
					panic(err)
				}
				hojuProb = tenpaiProb * probInfo.Prob
			}

			hm := core.NewHashMap[[4]float64]()
			hm.Set(hojuChanges, hojuProb)
			hm.Set(a.noChanges, 1.0-hojuProb)
			unitDist := core.NewVectorProbDist(hm)
			// Considers only the first ron for double/triple ron to avoid too many combinations.
			new := core.MultScalarVector(horaPointsDist, unitDist)
			scoreChangesDists[key] = scoreChangesDists[key].Replace(a.noChanges, new)
		}
	}

	return scoreChangesDists
}

func (a *ManueAI) getRyukyokuProbOnMyNoHora(state game.StateViewer) float64 {
	return math.Pow(a.getRyukyokuProb(state), 3.0/4.0)
}

func (a *ManueAI) getRandomHoraScoreChangesDist(
	state game.StateViewer,
	playerID int,
	actor *base.Player,
) *core.VectorProbDist {
	var horaPointsFreqs map[string]int
	if actor.ID() == state.Oya().ID() {
		horaPointsFreqs = a.stats.OyaHoraPointsFreqs
	} else {
		horaPointsFreqs = a.stats.KoHoraPointsFreqs
	}

	hm := core.NewHashMap[float64]()
	totalFreqs := float64(horaPointsFreqs["total"])
	for points, freq := range horaPointsFreqs {
		if points == "total" {
			continue
		}
		p, err := strconv.ParseFloat(points, 64)
		if err != nil {
			panic("Invalid stats file: failed to convert key of horaPointsFreqs to float64 (" + points + ").")
		}
		f := float64(freq) / totalFreqs
		hm.Set(p, f)
	}

	horaPointsDist := core.NewScalarProbDist(hm)
	horaFactorsDist := a.getHoraFactorsDist(state, playerID, actor)
	return core.MultScalarVector(horaPointsDist, horaFactorsDist)
}

func (a *ManueAI) getHoraFactorsDist(
	state game.StateViewer,
	playerID int,
	actor *base.Player,
) *core.VectorProbDist {
	tsumoHoraProb := float64(a.stats.NumTsumoHoras) / float64(a.stats.NumHoras)
	m := core.NewHashMap[[4]float64]()
	for _, target := range state.Players() {
		var prob float64
		if target.ID() == playerID {
			prob = tsumoHoraProb
		} else {
			prob = (1.0 - tsumoHoraProb) / 3.0
		}
		m.Set(a.getHoraFactors(state, actor, &target), prob)
	}
	return core.NewVectorProbDist(m)
}

func (a *ManueAI) getHoraFactors(state game.StateViewer, actor, target *base.Player) [4]float64 {
	actorID := actor.ID()
	targetID := target.ID()
	if actorID != targetID {
		// Ron hora
		horaFactors := [4]float64{0.0, 0.0, 0.0, 0.0}
		for i, p := range state.Players() {
			switch p.ID() {
			case actorID:
				horaFactors[i] = 1.0
			case targetID:
				horaFactors[i] = -1.0
			}
		}
		return horaFactors
	}

	oyaID := state.Oya().ID()
	// Tsumo hora
	if actorID == oyaID {
		// Oya tsumo hora
		horaFactors := [4]float64{-1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0}
		for i, p := range state.Players() {
			if p.ID() == actorID {
				horaFactors[i] = 1.0
			}
		}
		return horaFactors
	}

	// Ko tsumo hora
	horaFactors := [4]float64{-1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0}
	for i, p := range state.Players() {
		switch p.ID() {
		case actorID:
			horaFactors[i] = 1.0
		case oyaID:
			horaFactors[i] = -1.0 / 2.0
		}
	}
	return horaFactors
}

func (a *ManueAI) printTenpaiProbs(state game.StateViewer, playerID int) {
	var output strings.Builder
	output.WriteString("tenpaiProbs:  ")
	for _, p := range state.Players() {
		if p.ID() != playerID {
			fmt.Fprintf(&output, "%d: %.3f  ", p.ID(), a.tenpaiProbEstimator.Estimate(&p, state))
		}
	}
	output.WriteString("\n")
	a.log(output.String())
}

func (a *ManueAI) getNumExpectedRemainingTurns(state game.StateViewer) int {
	currentTurn := math.Round(state.Turn())
	num := 0.0
	den := 0.0
	ct := int(currentTurn)
	for i := ct; i < len(a.stats.NumTurnsDistribution); i++ {
		prob := a.stats.NumTurnsDistribution[i]
		num += prob * (float64(i) - currentTurn + 0.5)
		den += prob
	}
	if den == 0.0 {
		return 0
	} else {
		return int(math.Round(num / den))
	}
}

func (a *ManueAI) getRyukyokuProb(state game.StateViewer) float64 {
	currentTurn := int(state.Turn())
	den := 0.0
	for _, prob := range a.stats.NumTurnsDistribution[currentTurn:] {
		den += prob
	}
	return a.stats.RyukyokuRatio / den
}

func (a *ManueAI) getAverageRank(
	state game.StateViewer,
	playerID int,
	scoreChangesDist *core.VectorProbDist,
) float64 {
	hm1 := core.NewHashMap[[4]float64]()
	hm1.Set([4]float64{0.0, 0.0, 0.0, 0.0}, 1.0)
	winsDist := core.NewVectorProbDist(hm1)
	for _, other := range state.Players() {
		if other.ID() == playerID {
			continue
		}
		winProb := a.getWinProb(state, playerID, scoreChangesDist, &other)
		hm2 := core.NewHashMap[[4]float64]()
		hm2.Set([4]float64{0.0, 0.0, 0.0, 0.0}, 1.0-winProb)
		w := [4]float64{0.0, 0.0, 0.0, 0.0}
		w[other.ID()] = 1.0
		hm2.Set(w, winProb)
		d := core.NewVectorProbDist(hm2)
		winsDist = core.AddVectorVector(winsDist, d)
	}

	rankDist := winsDist.MapValueScalar(func(wins [4]float64) float64 {
		c, _ := core.Count(wins[:], func(w float64) (bool, error) {
			// Since w == 1.0 is problematic, a threshold is tentatively set
			return math.Abs(w-1.0) < 1e-5, nil
		})
		return float64(4 - c)
	})
	return rankDist.Expected()
}

func (a *ManueAI) getWinProb(
	state game.StateViewer,
	playerID int,
	scoreChangesDist *core.VectorProbDist,
	other *base.Player,
) float64 {
	me := &state.Players()[playerID]
	// TODO Change this considering renchan.
	nexttKyokuBakaze, nextKyokuNum := state.NextKyoku()
	myPos := game.GetPlayerDistance(me, state.Chicha())
	otherPos := game.GetPlayerDistance(other, state.Chicha())
	key := fmt.Sprintf("%s%d,%d,%d", nexttKyokuBakaze.ToString(), nextKyokuNum, myPos, otherPos)
	winProbs := a.stats.WinProbsMap[key]
	relativeScoreDist := scoreChangesDist.MapValueScalar(func(scoreChanges [4]float64) float64 {
		return (float64(me.Score()) + scoreChanges[playerID]) - (float64(other.Score()) + scoreChanges[other.ID()])
	})
	winProb := 0.0
	relativeScoreDist.Dist().ForEach(func(relativeScore, prob float64) {
		winProb += prob * a.getWinProbFromRelativeScore(relativeScore, winProbs, myPos, otherPos)
	})
	return winProb
}

func (a *ManueAI) getWinProbFromRelativeScore(
	relativeScore float64,
	winProbs map[string]float64,
	myPos int,
	otherPos int,
) float64 {
	if winProbs != nil {
		key := fmt.Sprintf("%.0f", relativeScore)
		if prob, ok := winProbs[key]; ok {
			return prob
		}
	}
	// abs(relativeScore) is so big that statistics are missing,
	// or the current kyoku is S-4 (orasu).
	if myPos < otherPos {
		if relativeScore >= 0.0 {
			return 1.0
		}
		return 0.0
	}
	if relativeScore > 0.0 {
		return 1.0
	}
	return 0.0
}
