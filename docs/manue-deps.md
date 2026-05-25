# Manue AI Dependency Interfaces

`ManueAgent` は最終的に stats / estimator / decision tree などを参照して評価値を計算する。ただし、`configs.GameStats` や danger tree を AI 内部構造へ丸ごとコピーすると、大きい設定値を二重保持しやすく、責務も曖昧になる。

そのため、AI 側は必要な読み口だけを小さい interface として定義し、`ManueAgentDeps` ではそれらを束ねた複合 interface を受け取る方針にする。`configs` の値は adapter 的な実装で interface を満たし、AI ロジックは JSON schema や embed の都合に直接依存しない。

## Target Interfaces

Manue 本体を `cmd/mjai-manue` から動かす段階では、stats と危険度推定に必要な依存を deps として渡す。AI package は `configs` を直接 import せず、外側が `configs.LoadGameStats()` / `configs.LoadDangerTree()` の結果をこの interface へ適合させる。

```go
func NewManueAgentWithDeps(seed uint64, deps ManueAgentDeps) (*ManueAgent, error)

type ManueAgentDeps struct {
    Stats  ManueStats
    Danger DangerEstimator
}

// ManueStats provides read-only access to immutable statistical data used by
// ManueAgent. Implementations must return stable values for the lifetime of the
// agent after validation.
type ManueStats interface {
    WinScoreStats
    RoundEndStats
    DrawTenpaiStats
    TenpaiEstimatorStats
    RankStats
    DealInStats
}

type WinScoreStats interface {
    NumWins() int
    NumSelfDrawWins() int
    NonDealerWinPointFreqs() map[string]int
    DealerWinPointFreqs() map[string]int
}

type RoundEndStats interface {
    TurnDistribution() []float64
    ExhaustiveDrawRatio() float64
}

type DrawTenpaiStats interface {
    ExhaustiveDrawNotenCount() int
    ExhaustiveDrawTenpaiTurnFreq(turnKey string) (freq int, ok bool)
}

type TenpaiEstimatorStats interface {
    YamitenCounts(remainTurns int, numMelds int) (total int, tenpai int, ok bool)
}

type DealInStats interface {
    AvgWinPts() float64
}

type RankStats interface {
    RelativeWinProbs(roundWind wind.Wind, roundNumber int, selfPosition int, otherPosition int) (map[string]float64, bool)
}

type DangerEstimator interface {
    DealInProb(scene DangerScene, discard tile.Tile) (float64, error)
}
```

`configs.GameStats` は上記の `WinScoreStats` / `RoundEndStats` / `DrawTenpaiStats` / `TenpaiEstimatorStats` / `RankStats` / `DealInStats` を構造的に満たす getter を持つ。これにより、AI package は `configs` を import せず、外側の組み立て側だけが `configs.LoadGameStats()` の戻り値を deps に渡せる。順位期待値用の `winProbsMap` key 形式（例: `E1,0,1`）は config schema 側の都合なので、AI 側の interface 境界には出さず、`configs.GameStats.RelativeWinProbs` の内部で組み立てる。

## Future Split

`ManueStats` は大きな単一 interface にせず、用途別の小さい interface を埋め込んで育てる。stats 以外の危険度推定は `ManueStats` に混ぜず、`DangerEstimator` として分ける。

`YamitenStat` や `RyukyokuTenpaiStat` のような config schema 型をそのまま AI interface で返すことは避ける。AI が使う値に寄せた getter にすることで、config の JSON 構造変更や分割の影響を狭くする。

CoffeeScript 版には `ManueAI#getTenpaiProb` と `TenpaiProbEstimator#estimate` の重複実装があるが、Go 版では `TenpaiEstimatorStats` を通した共通関数へ一本化する。差し替えは stats interface の実装差し替えで行い、別の estimator deps は特徴量ベース推定が必要になった時点で追加する。

## Stats Validation

stats はゲーム中に変化しない静的データなので、値の正当性は評価計算のたびに確認するのではなく、deps を組み立てる段階で一度検証する方針にする。これにより、壊れた設定値は早期に「設定ロード/組み立ての失敗」として検出でき、候補評価中の重複チェックも後から減らせる。

deps 付き constructor は `validateManueStats(stats ManueStats) error` と danger estimator の nil check を実行し、壊れた設定値を Agent 作成時に検出する。`NewManueAgent(seed)` はテストや簡易利用用に残してよいが、完成した `cmd/mjai-manue` では使わない。

validation は、`NumWins() > 0`、自摸和了数が和了数の範囲内であること、和了点頻度の `total` が正で点数 key が parse 可能であること、流局時聴牌確率に必要な turn key が存在すること、流局ノーテン数が非負であることなど、静的 stats の構造・範囲を対象にする。一方で actor/dealer ID や `currentTurn` のような呼び出しごとの引数 validation は、必要に応じて計算関数側に残す。

計算関数内に既にある stats 関連チェックは、validation 関数を追加した直後には無理に削らない。stats が Agent 作成時 validation を通る経路に一本化された後で、重複しているチェックを整理する。

## Candidate Score Helper Cleanup

`apply*ToCandidateScore` 系の小さい helper は、移植中に CoffeeScript 版 `getMetricsInternal` の各 metric 代入と対応を取りやすくし、estimator 未接続の段階でも純粋関数として検証するための足場である。

deps / estimator 接続後に候補評価の入力と責務が固まったら、単なるフィールド代入だけの helper は `evaluateCandidate...` 系の orchestration へ畳み込み、意味のある境界だけを残す。これは stats validation を Agent 作成時へ寄せた後に計算関数内の重複 validation を削る作業と並ぶ、後続リファクタリング項目として扱う。

## LightGameStats

`LightGameStats` は `game_stats.json` とは別の `light_game_stats.json` から読み込まれ、主用途は順位推定の `WinProbsMap` である。そのため、`configs` 側でも `GameStats` へ埋め込む必然性は薄く、将来的には分離したまま扱ってよい。

AI interface としては `RankStats` に対応させる。`GameStats` と `LightGameStats` の両方を持つ adapter が `ManueStats` を満たす形にすれば、AI は読み込み元 JSON の違いを意識しない。

## Danger Tree

`DecisionNode` / danger tree は stats とは別の設定値なので、`ManueStats` に混ぜない。AI は config schema ではなく、危険度推定を行う estimator interface に依存する方針にする。

```go
type DangerTreeNode interface {
    LeafProb() (float64, bool)
    Feature() (string, bool)
    NegativeNode() DangerTreeNode
    PositiveNode() DangerTreeNode
}
```

実装時に必要なら `configs.DecisionNode` に上記を満たす accessor を追加する。これにより、巨大な decision tree を AI 内部へコピーせず、estimator 実装が必要な参照だけを保持できる。
