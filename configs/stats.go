package configs

import (
	_ "embed"
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type YamitenStat struct {
	Total  int `json:"total"`
	Tenpai int `json:"tenpai"`
}

type RyukyokuTenpaiStat struct {
	Total                  int            `json:"total"`
	Tenpai                 int            `json:"tenpai"`
	Noten                  int            `json:"noten"`
	TenpaiTurnDistribution map[string]int `json:"tenpaiTurnDistribution"`
}

type LightGameStats struct {
	WinProbsMap map[string]map[string]float64 `json:"winProbsMap,omitempty"`
}

type GameStats struct {
	NumHoras             int                    `json:"numHoras"`
	NumTsumoHoras        int                    `json:"numTsumoHoras"`
	NumTurnsDistribution []float64              `json:"numTurnsDistribution"`
	RyukyokuRatio        float64                `json:"ryukyokuRatio"`
	AverageHoraPoints    float64                `json:"averageHoraPoints"`
	KoHoraPointsFreqs    map[string]int         `json:"koHoraPointsFreqs"`
	OyaHoraPointsFreqs   map[string]int         `json:"oyaHoraPointsFreqs"`
	YamitenStats         map[string]YamitenStat `json:"yamitenStats"`
	RyukyokuTenpaiStat   RyukyokuTenpaiStat     `json:"ryukyokuTenpaiStat"`
	LightGameStats
}

//go:embed game_stats.json
var rawGameStats []byte

//go:embed light_game_stats.json
var rawLightGameStats []byte

func LoadGameStats() (*GameStats, error) {
	var stats GameStats
	if err := json.Unmarshal(rawGameStats, &stats); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rawLightGameStats, &stats.LightGameStats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (s *GameStats) NumWins() int {
	return s.NumHoras
}

func (s *GameStats) NumSelfDrawWins() int {
	return s.NumTsumoHoras
}

func (s *GameStats) NonDealerWinPointFreqs() map[string]int {
	return s.KoHoraPointsFreqs
}

func (s *GameStats) DealerWinPointFreqs() map[string]int {
	return s.OyaHoraPointsFreqs
}

func (s *GameStats) TurnDistribution() []float64 {
	return s.NumTurnsDistribution
}

func (s *GameStats) ExhaustiveDrawRatio() float64 {
	return s.RyukyokuRatio
}

func (s *GameStats) AvgWinPts() float64 {
	return s.AverageHoraPoints
}

func (s *GameStats) ExhaustiveDrawNotenCount() int {
	return s.RyukyokuTenpaiStat.Noten
}

func (s *GameStats) ExhaustiveDrawTenpaiTurnFreq(turnKey string) (int, bool) {
	freq, ok := s.RyukyokuTenpaiStat.TenpaiTurnDistribution[turnKey]
	return freq, ok
}

func (s *GameStats) YamitenCounts(remainTurns int, numMelds int) (int, int, bool) {
	stat, ok := s.YamitenStats[fmt.Sprintf("%d,%d", remainTurns, numMelds)]
	if !ok {
		return 0, 0, false
	}
	return stat.Total, stat.Tenpai, true
}

func (s *GameStats) RelativeWinProbs(
	roundWind wind.Wind,
	roundNumber int,
	selfPosition int,
	otherPosition int,
) (map[string]float64, bool) {
	key := fmt.Sprintf("%s%d,%d,%d", roundWind, roundNumber, selfPosition, otherPosition)
	winProbs, ok := s.WinProbsMap[key]
	return winProbs, ok
}
