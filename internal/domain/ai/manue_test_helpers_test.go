package ai

type stubManueStats struct {
	numWins                       int
	numSelfDrawWins               int
	nonDealerWinPointFreqs        map[string]int
	dealerWinPointFreqs           map[string]int
	exhaustiveDrawNotenCount      int
	exhaustiveDrawTenpaiTurnFreqs map[string]int
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

func (s stubManueStats) ExhaustiveDrawNotenCount() int {
	return s.exhaustiveDrawNotenCount
}

func (s stubManueStats) ExhaustiveDrawTenpaiTurnFreq(turnKey string) (int, bool) {
	freq, ok := s.exhaustiveDrawTenpaiTurnFreqs[turnKey]
	return freq, ok
}
