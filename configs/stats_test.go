package configs

import (
	"math"
	"testing"
)

func TestGetStats(t *testing.T) {
	t.Run("stats test", func(t *testing.T) {
		epsilon := 1e-15

		got, err := GetStats()
		if err != nil {
			t.Errorf("GetStats() error = %v", err)
			return
		}

		if got.NumHoras != 17793 {
			t.Errorf("GetStats().NumHoras = %v, want %v", got.NumHoras, 17793)
		}
		if got.NumTsumoHoras != 7195 {
			t.Errorf("GetStats().NumTsumoHoras = %v, want %v", got.NumTsumoHoras, 7195)
		}
		if math.Abs(got.NumTurnsDistribution[17]-0.16677791551444646) > epsilon {
			t.Errorf("GetStats().NumTurnsDistribution[17] = %v, want %v", got.NumTurnsDistribution[8], 0.16677791551444646)
		}
		if math.Abs(got.RyukyokuRatio-0.15700390960236482) > epsilon {
			t.Errorf("GetStats().RyukyokuRatio = %v, want %v", got.RyukyokuRatio, 0.15700390960236482)
		}
		if math.Abs(got.AverageHoraPoints-5533.648063845332) > epsilon {
			t.Errorf("GetStats().AverageHoraPoints = %v, want %v", got.AverageHoraPoints, 5533.648063845332)
		}
		if got.KoHoraPointsFreqs["total"] != 12834 {
			t.Errorf("GetStats().KoHoraPointsFreqs[\"total\"] = %v, want %v", got.KoHoraPointsFreqs["total"], 12834)
		}
		if got.OyaHoraPointsFreqs["1500"] != 555 {
			t.Errorf("GetStats().OyaHoraPointsFreqs[\"1500\"] = %v, want %v", got.OyaHoraPointsFreqs["1500"], 555)
		}
		if got.YamitenStats["17,0"].Total != 41903 {
			t.Errorf("GetStats().YamitenStats[\"17,0\"].Total = %v, want %v", got.YamitenStats["17,0"].Total, 41903)
		}
		if got.YamitenStats["12,4"].Tenpai != 4 {
			t.Errorf("GetStats().YamitenStats[\"12,4\"].Tenpai = %v, want %v", got.YamitenStats["12,4"].Tenpai, 4)
		}
		if got.RyukyokuTenpaiStat.Total != 13172 {
			t.Errorf("GetStats().RyukyokuTenpaiStat.Total = %v, want %v", got.RyukyokuTenpaiStat.Total, 13172)
		}
		if got.RyukyokuTenpaiStat.Tenpai != 5468 {
			t.Errorf("GetStats().RyukyokuTenpaiStat.Tenpai = %v, want %v", got.RyukyokuTenpaiStat.Tenpai, 5468)
		}
		if got.RyukyokuTenpaiStat.Noten != 7704 {
			t.Errorf("GetStats().RyukyokuTenpaiStat.Noten = %v, want %v", got.RyukyokuTenpaiStat.Noten, 7704)
		}
		if got.RyukyokuTenpaiStat.TenpaiTurnDistribution["0"] != 0 {
			t.Errorf("GetStats().RyukyokuTenpaiStat.TenpaiTurnDistribution[\"0\"] = %v, want %v", got.RyukyokuTenpaiStat.TenpaiTurnDistribution["0"], 0)
		}
		// "null" is not referenced in normal program flow, so it is not a problem that it is 0 instead of nil.
		if got.RyukyokuTenpaiStat.TenpaiTurnDistribution["null"] != 0 {
			t.Errorf("GetStats().RyukyokuTenpaiStat.TenpaiTurnDistribution[\"null\"] = %v, want %v", got.RyukyokuTenpaiStat.TenpaiTurnDistribution["null"], 0)
		}

		if math.Abs(got.WinProbsMap["E1,0,1"]["0"]-0.49478259990894036) > epsilon {
			t.Errorf("GetStats().WinProbsMap[\"E1,0,1\"][\"0\"] = %v, want %v", got.WinProbsMap["E1,0,1"]["0"], 0.49478259990894036)
		}
	})
}
