package configs

import (
	"math"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

func TestLoadGameStats(t *testing.T) {
	epsilon := 1e-15

	got, err := LoadGameStats()
	if err != nil {
		t.Fatalf("LoadGameStats() error = %v", err)
		return
	}

	if got.NumHoras != 17793 {
		t.Errorf("LoadGameStats().NumHoras = %v, want %v", got.NumHoras, 17793)
	}
	if got.NumTsumoHoras != 7195 {
		t.Errorf("LoadGameStats().NumTsumoHoras = %v, want %v", got.NumTsumoHoras, 7195)
	}
	if math.Abs(got.NumTurnsDistribution[17]-0.16677791551444646) > epsilon {
		t.Errorf("LoadGameStats().NumTurnsDistribution[17] = %v, want %v", got.NumTurnsDistribution[8], 0.16677791551444646)
	}
	if math.Abs(got.RyukyokuRatio-0.15700390960236482) > epsilon {
		t.Errorf("LoadGameStats().RyukyokuRatio = %v, want %v", got.RyukyokuRatio, 0.15700390960236482)
	}
	if math.Abs(got.AverageHoraPoints-5533.648063845332) > epsilon {
		t.Errorf("LoadGameStats().AverageHoraPoints = %v, want %v", got.AverageHoraPoints, 5533.648063845332)
	}
	if got.KoHoraPointsFreqs["total"] != 12834 {
		t.Errorf("LoadGameStats().KoHoraPointsFreqs[\"total\"] = %v, want %v", got.KoHoraPointsFreqs["total"], 12834)
	}
	if got.OyaHoraPointsFreqs["1500"] != 555 {
		t.Errorf("LoadGameStats().OyaHoraPointsFreqs[\"1500\"] = %v, want %v", got.OyaHoraPointsFreqs["1500"], 555)
	}
	if got.YamitenStats["17,0"].Total != 41903 {
		t.Errorf("LoadGameStats().YamitenStats[\"17,0\"].Total = %v, want %v", got.YamitenStats["17,0"].Total, 41903)
	}
	if got.YamitenStats["12,4"].Tenpai != 4 {
		t.Errorf("LoadGameStats().YamitenStats[\"12,4\"].Tenpai = %v, want %v", got.YamitenStats["12,4"].Tenpai, 4)
	}
	if _, ok := got.YamitenStats["missing"]; ok {
		t.Errorf("LoadGameStats().YamitenStats[\"missing\"] exists unexpectedly")
	}
	if got.RyukyokuTenpaiStat.Total != 13172 {
		t.Errorf("LoadGameStats().RyukyokuTenpaiStat.Total = %v, want %v", got.RyukyokuTenpaiStat.Total, 13172)
	}
	if got.RyukyokuTenpaiStat.Tenpai != 5468 {
		t.Errorf("LoadGameStats().RyukyokuTenpaiStat.Tenpai = %v, want %v", got.RyukyokuTenpaiStat.Tenpai, 5468)
	}
	if got.RyukyokuTenpaiStat.Noten != 7704 {
		t.Errorf("LoadGameStats().RyukyokuTenpaiStat.Noten = %v, want %v", got.RyukyokuTenpaiStat.Noten, 7704)
	}
	if got.RyukyokuTenpaiStat.TenpaiTurnDistribution["0"] != 0 {
		t.Errorf("LoadGameStats().RyukyokuTenpaiStat.TenpaiTurnDistribution[\"0\"] = %v, want %v", got.RyukyokuTenpaiStat.TenpaiTurnDistribution["0"], 0)
	}
	// "null" is not referenced in normal program flow, so it is not a problem that it is 0 instead of nil.
	if got.RyukyokuTenpaiStat.TenpaiTurnDistribution["null"] != 0 {
		t.Errorf("LoadGameStats().RyukyokuTenpaiStat.TenpaiTurnDistribution[\"null\"] = %v, want %v", got.RyukyokuTenpaiStat.TenpaiTurnDistribution["null"], 0)
	}

	if math.Abs(got.WinProbsMap["E1,0,1"]["0"]-0.49478259990894036) > epsilon {
		t.Errorf("LoadGameStats().WinProbsMap[\"E1,0,1\"][\"0\"] = %v, want %v", got.WinProbsMap["E1,0,1"]["0"], 0.49478259990894036)
	}
}

func TestGameStats_WinScoreStats(t *testing.T) {
	got, err := LoadGameStats()
	if err != nil {
		t.Fatalf("LoadGameStats() error = %v", err)
		return
	}

	var winScoreStats ai.WinScoreStats = got
	if winScoreStats.NumWins() != got.NumHoras {
		t.Errorf("NumWins() = %v, want %v", winScoreStats.NumWins(), got.NumHoras)
	}
	if winScoreStats.NumSelfDrawWins() != got.NumTsumoHoras {
		t.Errorf("NumSelfDrawWins() = %v, want %v", winScoreStats.NumSelfDrawWins(), got.NumTsumoHoras)
	}
	if winScoreStats.NonDealerWinPointFreqs()["total"] != got.KoHoraPointsFreqs["total"] {
		t.Errorf("NonDealerWinPointFreqs()[\"total\"] = %v, want %v", winScoreStats.NonDealerWinPointFreqs()["total"], got.KoHoraPointsFreqs["total"])
	}
	if winScoreStats.DealerWinPointFreqs()["1500"] != got.OyaHoraPointsFreqs["1500"] {
		t.Errorf("DealerWinPointFreqs()[\"1500\"] = %v, want %v", winScoreStats.DealerWinPointFreqs()["1500"], got.OyaHoraPointsFreqs["1500"])
	}
}

func TestGameStats_DrawTenpaiStats(t *testing.T) {
	got, err := LoadGameStats()
	if err != nil {
		t.Fatalf("LoadGameStats() error = %v", err)
		return
	}

	var drawTenpaiStats ai.DrawTenpaiStats = got
	if drawTenpaiStats.ExhaustiveDrawNotenCount() != got.RyukyokuTenpaiStat.Noten {
		t.Errorf("ExhaustiveDrawNotenCount() = %v, want %v", drawTenpaiStats.ExhaustiveDrawNotenCount(), got.RyukyokuTenpaiStat.Noten)
	}
	freq, ok := drawTenpaiStats.ExhaustiveDrawTenpaiTurnFreq("17")
	if !ok {
		t.Errorf("ExhaustiveDrawTenpaiTurnFreq(\"17\") ok = false, want true")
	}
	if freq != got.RyukyokuTenpaiStat.TenpaiTurnDistribution["17"] {
		t.Errorf("ExhaustiveDrawTenpaiTurnFreq(\"17\") = %v, want %v", freq, got.RyukyokuTenpaiStat.TenpaiTurnDistribution["17"])
	}
	freq, ok = drawTenpaiStats.ExhaustiveDrawTenpaiTurnFreq("0.75")
	if !ok {
		t.Errorf("ExhaustiveDrawTenpaiTurnFreq(\"0.75\") ok = false, want true")
	}
	if freq != 0 {
		t.Errorf("ExhaustiveDrawTenpaiTurnFreq(\"0.75\") = %v, want 0", freq)
	}
	if _, ok := drawTenpaiStats.ExhaustiveDrawTenpaiTurnFreq("missing"); ok {
		t.Errorf("ExhaustiveDrawTenpaiTurnFreq(\"missing\") ok = true, want false")
	}
}
