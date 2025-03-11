package configs

import (
	_ "embed"

	"github.com/go-json-experiment/json"
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
	WinProbsMap map[string]map[string]float64 `json:"winProbsMap"`
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

var Stats GameStats

// Call only once throughout the program.
func InitializeStats() {
	if err := json.Unmarshal(rawGameStats, &Stats); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(rawLightGameStats, &Stats.LightGameStats); err != nil {
		panic(err)
	}
}
