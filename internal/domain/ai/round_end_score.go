package ai

// futureScoreDeltaDist mixes possible future round-ending outcomes after a
// discard did not immediately deal in.
func futureScoreDeltaDist(
	selfWinDist scoreDeltaProbDist,
	selfWinProb float64,
	exhaustiveDrawDist scoreDeltaProbDist,
	exhaustiveDrawProb float64,
	otherWinDists []scoreDeltaProbDist,
	otherWinProb float64,
) scoreDeltaProbDist {
	items := []weightedScoreDeltaProbDist{
		{dist: selfWinDist, prob: selfWinProb},
		{dist: exhaustiveDrawDist, prob: exhaustiveDrawProb},
	}
	perOtherWinProb := 0.0
	if len(otherWinDists) > 0 {
		perOtherWinProb = otherWinProb / float64(len(otherWinDists))
	}
	for _, dist := range otherWinDists {
		items = append(items, weightedScoreDeltaProbDist{dist: dist, prob: perOtherWinProb})
	}
	return mergeScoreDeltaProbDists(items)
}
