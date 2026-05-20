# Manue AI Dependency Interfaces

`ManueAgent` は最終的に stats / estimator / decision tree などを参照して評価値を計算する。ただし、`configs.GameStats` や danger tree を AI 内部構造へ丸ごとコピーすると、大きい設定値を二重保持しやすく、責務も曖昧になる。

そのため、AI 側は必要な読み口だけを小さい interface として定義し、`ManueAgentDeps` ではそれらを束ねた複合 interface を受け取る方針にする。`configs` の値は adapter 的な実装で interface を満たし、AI ロジックは JSON schema や embed の都合に直接依存しない。

## Current Interfaces

現時点で実装している読み口は、ランダム和了時の score change 分布と流局時聴牌確率の基礎計算に必要なものだけ。

```go
func NewManueAgentWithDeps(seed uint64, deps ManueAgentDeps) *ManueAgent

type ManueAgentDeps struct {
    Stats ManueStats
}

// ManueStats provides read-only access to immutable statistical data used by
// ManueAgent. Implementations must return stable values for the lifetime of the
// agent after validation.
type ManueStats interface {
    WinScoreStats
    DrawTenpaiStats
}

type WinScoreStats interface {
    NumWins() int
    NumSelfDrawWins() int
    NonDealerWinPointFreqs() map[string]int
    DealerWinPointFreqs() map[string]int
}

type DrawTenpaiStats interface {
    ExhaustiveDrawNotenCount() int
    ExhaustiveDrawTenpaiTurnFreq(turnKey string) (freq int, ok bool)
}
```

`configs.GameStats` は上記の `WinScoreStats` / `DrawTenpaiStats` を構造的に満たす getter を持つ。これにより、AI package は `configs` を import せず、外側の組み立て側だけが `configs.LoadGameStats()` の戻り値を deps に渡せる。

## Planned Split

`ManueStats` は大きな単一 interface にせず、用途別の小さい interface を埋め込んで育てる。

```go
type ManueStats interface {
    WinScoreStats
    RoundEndStats
    DrawTenpaiStats
    TenpaiEstimatorStats
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

type TenpaiEstimatorStats interface {
    YamitenCounts(remainTurns int, numMelds int) (total int, tenpai int, ok bool)
}

type RankStats interface {
    RelativeWinProbTable(key string) (relativeWinProbTable, bool)
}

type DealInStats interface {
    AverageWinPoints() float64
}
```

`YamitenStat` や `RyukyokuTenpaiStat` のような config schema 型をそのまま AI interface で返すことは避ける。AI が使う値に寄せた getter にすることで、config の JSON 構造変更や分割の影響を狭くする。

## Stats Validation

stats はゲーム中に変化しない静的データなので、値の正当性は評価計算のたびに確認するのではなく、deps を組み立てる段階で一度検証する方針にする。これにより、壊れた設定値は早期に「設定ロード/組み立ての失敗」として検出でき、候補評価中の重複チェックも後から減らせる。

現時点では private な `validateManueStats(stats ManueStats) error` を純粋関数として追加し、`NewManueAgentWithDeps` にはまだ組み込まない。`NewManueAgent(seed)` は configs 非依存の CLI 起動経路として残しており、stats が実際に必須になった段階で deps 付き constructor の error 返却や validation 組み込みを検討する。

validation は、`NumWins() > 0`、自摸和了数が和了数の範囲内であること、和了点頻度の `total` が正で点数 key が parse 可能であること、流局時聴牌確率に必要な turn key が存在すること、流局ノーテン数が非負であることなど、静的 stats の構造・範囲を対象にする。一方で actor/dealer ID や `currentTurn` のような呼び出しごとの引数 validation は、必要に応じて計算関数側に残す。

計算関数内に既にある stats 関連チェックは、validation 関数を追加した直後には無理に削らない。stats が Agent 作成時 validation を通る経路に一本化された後で、重複しているチェックを整理する。

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
