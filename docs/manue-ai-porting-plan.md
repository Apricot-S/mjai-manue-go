# Manue AI Porting Plan

この文書は、`internal/domain/ai` 配下で CoffeeScript 版 `manue_ai.coffee` のロジックを移植し直すための作業計画である。実装時は `docs/design.md` を一次設計として読み、本書は Manue AI 移植に閉じた作業指示として使う。

## 1. 参照元と再利用方針

- ロジックの一次資料は `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee` とする。
- 過去 Go 実装 `reference/repositories/mjai-manue-go-main/internal/ai/manue_ai.go` は、`DecideAction` から outbound action へつながる動線と trace/log 形式を確認する補助資料として使う。
- `domain/game` 側に実装済みのルール判定は再実装しない。特に `round.State.LegalActions`、`round.ActionStateViewer`、`service.AnalyzeShanten`、`service.CalculateFuHan`、`service.RonPoints`、`service.RyukyokuPoints`、`service.IsTenpai*` を使う。
- `internal/domain/ai` は `configs` を直接 import しない。stats / danger tree は `cmd/mjai-manue` など外側で load し、AI 側の小さい interface へ渡す。
- 局面 observation は `domain/game` 側の `round.ActionStateViewer` / `round.StateViewer` を使う。AI package 内で同等の state viewer interface を重複定義せず、純粋関数の小さい viewer interface は局所的な入力制約としてだけ残す。
- 移植元の乱数列完全一致は目標にしない。Go stdlib PCG と `--seed` による再現性を優先する。

## 2. オリジナル実装の分解

CoffeeScript 版 `ManueAI` は大きく次の責務を 1 class に混在させている。

- `respondToAction` / `DecideAction`: イベント反応、和了優先、打牌判断、副露判断への dispatch。
- `decideDahai`: 打牌候補と立直+打牌候補を評価し、最善 key を action に戻す。
- `decideFuro`: `none` と各副露候補を評価し、鳴くか見送るかを選ぶ。鳴く候補は「鳴いた後に何を切るか」まで評価している。
- `getMetrics` / `getMetricsInternal`: 候補ごとの評価値を組み立てる中心処理。
- `getSafeProbs` / `getImmediateScoreChangesDists`: 危険度推定、聴牌確率、deal-in 時 score delta。
- `getHoraEstimation`: 未見牌から Monte Carlo で自分の win probability / win-points distribution / shanten を推定する。
- `getRyukyoku*` / `getRandomHoraScoreChangesDist`: exhaustive draw と other-player win の score delta 分布。
- `getAverageRank` / `getWinProb`: score delta 分布から順位期待値を計算する。
- `printMetrics` / `decidedKey`: protocol output とは別の診断 trace。

再構築では、上記を `ManueAgent.Decide` から呼ばれる小さい orchestration と純粋関数へ分ける。移植元の `metric` は Go 側では `candidateEvaluation` のような構造にまとめ、表示用 trace と選択ロジックの入力にだけ使う。

## 3. 現行 `internal/domain/ai` の扱い

残すもの:

- `agent.go` と `tsumogiri_agent.go`: 既存 interface と最小 AI。
- `manue_prob_dist.go`: scalar / score delta / ahead vector の確率分布演算。移植元 `ProbDist` 相当として使える。
- `manue_score.go`: score delta、和了点分布、流局点、期待点など、移植元の式と対応が明確な純粋関数。
- `manue_rank.go`: score delta 分布から順位期待値を作る純粋関数。
- `manue_trace.go`: CoffeeScript 版の metric table と `decidedKey` に対応する trace 整形。
- `manue_candidate.go` 内の壁生成、試行ツモ、候補別 accumulator、Goal scoring など、`getHoraEstimation` 相当へ直接使える関数。

削除または畳み込み済みのもの:

- `apply*ToCandidateScore` のような単なるフィールド代入 helper。
- `candidateScoreDeltaDists` のように、未接続 wrapper としてしか働いていない構造。
- `Phase 1` / `Phase 2 scaffold` コメント付きの暫定 fallback コメント。
- 評価の入力と出力を曖昧にする helper。必要な式は、通常手番評価または副露反応評価の orchestration へ直接組み込む。

## 4. 再構築後の動線

`ManueAgent.Decide` は薄い dispatcher にする。

```text
ManueAgent.Decide
  ├─ legal actions from request.Round.LegalActions(request.Self)
  ├─ Win action selection
  ├─ decideSelfTurn
  │   ├─ riichi accepted: tsumogiri discard only
  │   ├─ buildSelfTurnCandidates
  │   ├─ evaluateActionCandidates
  │   └─ chooseBestCandidate(preferBlack=true)
  └─ decideOtherDiscardReaction
      ├─ buildReactionCandidates
      ├─ evaluateActionCandidates
      └─ chooseBestCandidate(preferBlack=false)
```

通常手番:

- 合法手から `Discard` を候補化する。
- 合法手に同一牌の手出し discard とツモ切り discard が両方ある場合は、候補生成時点で手出し側を省き、ツモ切り側だけを残す。CoffeeScript 版は `tehais` から候補 key を作り、選択牌が最後のツモ牌と `Pai.equal` で同一なら action 出力時に `tsumogiri=true` としていたため、現行 Go の action 差分はこの正規化で吸収する。赤牌差分は `Pai.equal` と同じく同一扱いしない。
- 合法手に `Riichi` がある場合、CoffeeScript 版と同様に通常手の向聴数が 0 以下になる打牌だけ `Riichi + Discard` 候補を追加する。
- リーチ宣言後・成立前の打牌は、action と trace key は通常打牌のまま扱い、候補評価だけ `reachMode="now"` 相当としてリーチあり点数・聴牌 goal filtering を使う。
- 暗槓、加槓、九種九牌は今回の評価対象外とし、合法でも候補にしない。和了だけは常に即選択する。

他家打牌反応:

- `Pass` を `none` 候補として評価対象に含める。
- `Chii` / `Pon` / `CalledKan` は、副露後の手牌と副露リストを作り、喰い替えになる打牌を除外したうえで「鳴いた後の最善打牌」を評価する。
- 最終選択が `none` なら `Pass` を返し、それ以外なら対応する鳴き action を返す。鳴いた直後の打牌は次の action opportunity で改めて判断する。

## 5. 候補評価

候補評価は CoffeeScript 版 `getMetricsInternal` と同じ意味の値を作る。ただし移植実装の構造体・関数・trace 以外の内部名は、現行 domain の語彙に合わせて原則 English にする。CoffeeScript 由来の `hoju` は `dealIn`、`hora` は `win`、`ryukyoku` は `exhaustiveDraw` として表す。

- `safeProb`: 全相手に対して deal in しない確率。
- `dealInProb`: `1 - safeProb`。
- `winProb`: その候補後に自分が win する確率。
- `averageWinPoints`: 自分が win した場合の平均点。
- `winPointsDist`: 自分が win した場合の点数分布。
- `expectedWinPoints`: Monte Carlo 全試行に対する win points 期待値。
- `exhaustiveDrawProb`: 自分が win しなかった場合を含めた exhaustive draw 確率。
- `exhaustiveDrawAveragePoints`: exhaustive draw 時の自分の平均点差。
- `immediateScoreDeltaDist`: 直後 deal-in の score delta 分布。
- `futureScoreDeltaDist`: self win、exhaustive draw、other-player win の score delta 分布。
- `scoreDeltaDist`: 即時分布の no-change branch を future 分布で置き換えた最終 score delta 分布。
- `expectedPoints`: `scoreDeltaDist` における自分の期待点。
- `averageRank`: `scoreDeltaDist` と順位統計から計算した平均順位。
- `shanten` / `red`: trace と tie-break 用。

選択規則は CoffeeScript 版と同じく、平均順位が小さい候補、期待点が大きい候補、赤牌を切らない候補の順で比較する。副露判断では `preferBlack=false` とし、赤牌 tie-break は使わない。

## 6. 危険度推定

危険度推定は `internal/domain/ai` に移植する。ただし `configs.DecisionNode` へ直接依存せず、AI 側で以下のような読み口を定義する。

```go
type DangerTreeNode interface {
    LeafProb() (float64, bool)
    Feature() (string, bool)
    NegativeNode() DangerTreeNode
    PositiveNode() DangerTreeNode
}
```

実装時に必要なら `configs.DecisionNode` に上記を満たす accessor を追加する。`configs.DecisionNode` は accessor の戻り値型として `ai.DangerTreeNode` を返してよい。依存方向は `configs -> internal/domain/ai` なので循環 import にはならない。AI 側の `DangerEstimator` は tree traversal と feature evaluator だけを持つ。

`Scene` は現行 `round.StateViewer` から組み立てる。

- 自分の手牌: `Player(request.Self).Hand()`
- 安全牌: `round.State.SafeTiles(targetSeat)`
- 見えている牌: `round.State.VisibleTiles(selfSeat)`
- ドラ: `round.State.Doras()`
- 場風 / 自風: `RoundWind()` / `SeatWind(targetSeat)`
- リーチ前捨て牌、リーチ宣言牌、早い/遅い捨て牌: `PlayerViewer` の河と riichi state から導出する。

feature evaluator 群は旧 Go 実装 `estimator/danger_estimator.go` を現行 `tile` / `hand` / `player` 型へ置き換える。`anpai` は shortcut として probability 0 を返し、その他は decision tree の feature 評価に従う。

## 7. 和了推定

和了推定は既存の `manue_candidate.go` の使える関数を整理して接続する。

- 見えている牌から未見牌の壁を作る。
- `expectedRemainingTurns` で試行ツモ枚数を決める。
- `service.AnalyzeShanten` の `Goal` と `ThrowableVector` を使い、候補ごとの達成可能性を判定する。
- 候補ごとの `shanten` は打牌後の手牌を再解析するのではなく、CoffeeScript 版の `shantenVector` と同じく、打牌前/副露後の `Goal.ThrowableVector` からその牌を捨てられる goal の最小 `Shanten` を採用する。`none` は base shanten を使う。
- 和了推定 goal の pruning は候補ごとの `shanten` ではなく、CoffeeScript 版の `analysis.shanten()` と同じ打牌前/副露後の base shanten を使う。
- `service.CalculateFuHan` と `service.RonPoints` で点数を計算する。`CalculateFuHan` の hand 引数は `handBlocks` と同じ構成牌を持つ concealed hand を前提とするため、候補 scoring では `Goal.Blocks` から scoring hand を組み立てて渡す。
- 赤ドラは CoffeeScript 版の `pai.red() && any(allPais, sameSymbol)` と同じく、候補 hand / meld 内の赤牌のうち goal の構成牌に同種牌が含まれるものだけを scoring hand / meld に残して `adr` として数える。
- 候補別 accumulator で `winProb`、`averageWinPoints`、`winPointsDist`、`expectedWinPoints` を作る。

Monte Carlo はまず直列で実装する。並列化が必要になった場合は、worker ごとに候補別 accumulator を作って merge する。

## 8. Configs と main.go 動線

`cmd/mjai-manue/main.go` から Manue 本体を動かすため、外側で設定を読み込んで deps に渡す。

- `configs.LoadGameStats()` で stats を読む。
- `configs.LoadDangerTree()` で danger tree を読む。
- `ai.NewManueAgent(seed, deps)` は deps validation を行い、失敗したら error を返す。
- `cmd/mjai-manue` は deps 構築エラーを stderr に出し、runtime error として終了する。
- runtime / application / adapter の既存動線は変更しない。

旧 `NewManueAgent(seed)` と `NewManueAgentWithDeps(seed, deps)` の二重入口は廃止し、`NewManueAgent(seed, deps)` に一本化する。

## 9. 実装順

1. `docs/design.md` と本書に従い、足場 helper を削除して評価構造を `actionCandidate` と `candidateScore` 中心に整理する。単なる field 代入 helper と未接続 wrapper は削除済み。
2. `ManueAgent.Decide` を薄い dispatcher に戻し、通常手番と他家打牌反応の orchestration を分ける。
3. stats validation と `cmd/mjai-manue` の deps 構築を実装する。
4. 和了推定を通常手番評価に接続する。
5. 流局・他家和了・順位期待値を接続する。
6. 危険度推定を移植し、`safeProb` / `dealInProb` / immediate score delta distribution を実値化する。
7. 副露候補の「鳴いた後の最善打牌」評価を接続する。副露後 meld は候補ごとの和了点推定に反映する。
8. trace table を CoffeeScript 版の列と key 形式に揃える。

## 10. テスト方針

- 純粋関数は `internal/domain/ai` でテーブル駆動テストにする。
- 危険度 feature は feature ごとに小さい scene fixture を作る。
- Agent 判断は、private field を直接壊さず、mjai JSON Lines から `round.State` を構築して action を比較する。
- runtime golden は action のみ比較する。trace / log は別 fixture にする。
- 最終確認は PowerShell で `$env:GOEXPERIMENT='jsonv2'; go test ./...` を実行する。
