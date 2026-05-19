# Manue AI Dependency Interfaces

`ManueAgent` は最終的に stats / estimator / decision tree などを参照して評価値を計算する。ただし、`configs.GameStats` や danger tree を AI 内部構造へ丸ごとコピーすると、大きい設定値を二重保持しやすく、責務も曖昧になる。

そのため、AI 側は必要な読み口だけを小さい interface として定義し、`ManueAgentDeps` ではそれらを束ねた複合 interface を受け取る方針にする。`configs` の値は adapter 的な実装で interface を満たし、AI ロジックは JSON schema や embed の都合に直接依存しない。

## Current Interfaces

現時点で実装している読み口は、ランダム和了時の score change 分布に必要なものだけ。

```go
func NewManueAgentWithDeps(seed uint64, deps ManueAgentDeps) *ManueAgent

type ManueAgentDeps struct {
    Stats ManueStats
}

type ManueStats interface {
    WinScoreStats
}

type WinScoreStats interface {
    NumWins() int
    NumSelfDrawWins() int
    NonDealerWinPointFreqs() map[string]int
    DealerWinPointFreqs() map[string]int
}
```

`configs.GameStats` は上記の `WinScoreStats` を構造的に満たす getter を持つ。これにより、AI package は `configs` を import せず、外側の組み立て側だけが `configs.LoadGameStats()` の戻り値を deps に渡せる。

## Planned Split

`ManueStats` は大きな単一 interface にせず、用途別の小さい interface を埋め込んで育てる。

```go
type ManueStats interface {
    WinScoreStats
    RoundEndStats
    TenpaiStats
    RankStats
    DealInStats
}
```

候補:

```go
type RoundEndStats interface {
    NumTurnsDistribution() []float64
    ExhaustiveDrawRatio() float64
}

type TenpaiStats interface {
    YamitenCounts(remainTurns int, numMelds int) (total int, tenpai int, ok bool)
    ExhaustiveDrawNotenCount() int
    ExhaustiveDrawTenpaiTurnCount(turnKey string) int
}

type RankStats interface {
    RelativeWinProbTable(key string) (relativeWinProbTable, bool)
}

type DealInStats interface {
    AverageWinPoints() float64
}
```

`YamitenStat` や `RyukyokuTenpaiStat` のような config schema 型をそのまま AI interface で返すことは避ける。AI が使う値に寄せた getter にすることで、config の JSON 構造変更や分割の影響を狭くする。

## LightGameStats

`LightGameStats` は `game_stats.json` とは別の `light_game_stats.json` から読み込まれ、主用途は順位推定の `WinProbsMap` である。そのため、`configs` 側でも `GameStats` へ埋め込む必然性は薄く、将来的には分離したまま扱ってよい。

AI interface としては `RankStats` に対応させる。`GameStats` と `LightGameStats` の両方を持つ adapter が `ManueStats` を満たす形にすれば、AI は読み込み元 JSON の違いを意識しない。

## Danger Tree

`DecisionNode` / danger tree は stats とは別の設定値なので、`ManueStats` に混ぜない。AI は tree そのものではなく、危険度推定を行う estimator interface に依存する方針にする。

```go
type DangerEstimator interface {
    // Details will be added when danger estimation is migrated.
}
```

これにより、巨大な decision tree を AI 内部へコピーせず、estimator 実装が必要な参照だけを保持できる。
