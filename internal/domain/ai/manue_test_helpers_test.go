package ai

import "fmt"

type stubManueStats struct {
	numWins                       int
	numSelfDrawWins               int
	nonDealerWinPointFreqs        map[string]int
	dealerWinPointFreqs           map[string]int
	turnDistribution              []float64
	exhaustiveDrawRatio           float64
	exhaustiveDrawNotenCount      int
	exhaustiveDrawTenpaiTurnFreqs map[string]int
	yamitenCounts                 map[string]yamitenCount
}

type yamitenCount struct {
	total  int
	tenpai int
}

func (s stubManueStats) NumWins() int {
	return s.numWins
}

func (s stubManueStats) NumSelfDrawWins() int {
	return s.numSelfDrawWins
}

func (s stubManueStats) NonDealerWinPointFreqs() map[string]int {
	return s.nonDealerWinPointFreqs
}

func (s stubManueStats) DealerWinPointFreqs() map[string]int {
	return s.dealerWinPointFreqs
}

func (s stubManueStats) TurnDistribution() []float64 {
	return s.turnDistribution
}

func (s stubManueStats) ExhaustiveDrawRatio() float64 {
	return s.exhaustiveDrawRatio
}

func (s stubManueStats) ExhaustiveDrawNotenCount() int {
	return s.exhaustiveDrawNotenCount
}

func (s stubManueStats) ExhaustiveDrawTenpaiTurnFreq(turnKey string) (int, bool) {
	freq, ok := s.exhaustiveDrawTenpaiTurnFreqs[turnKey]
	return freq, ok
}

func (s stubManueStats) YamitenCounts(remainTurns int, numMelds int) (int, int, bool) {
	count, ok := s.yamitenCounts[fmt.Sprintf("%d,%d", remainTurns, numMelds)]
	if !ok {
		return 0, 0, false
	}
	return count.total, count.tenpai, true
}
