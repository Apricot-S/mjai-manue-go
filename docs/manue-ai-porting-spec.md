# Manue AI Porting Specification

この文書は CoffeeScript 版 Manue AI を Go へ移植するための設計仕様である。原仕様は `docs/manue-ai-original-spec.md` を一次資料とし、本書では現行 `domain/game` と `application` へ接続するための責務分割だけを定める。

## 1. 境界と方針

- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee` をロジックの一次資料にする。
- 過去 Go 実装と現行 `internal/domain/ai` の Manue 固有コードは補助資料に留める。流用は仕様に合う純粋関数だけに限定する。
- `agent.go` と `tsumogiri_agent.go` は再実装対象外。Manue 固有コードだけを新規設計で作り直す。
- `internal/domain/ai` は `configs` を直接 import しない。stats、danger tree、estimator は外側で load し、小さい interface または plain data として渡す。
- 局面 observation は `round.ActionStateViewer` / `round.StateViewer` を使う。AI 側に同等の state viewer を重複定義しない。
- `round.State`、`EventApplier`、`LegalActions` は現状の責務を維持し、AI 再実装のために状態遷移や合法手列挙を別 service 化しない。

## 2. Agent 動線

`ManueAgent.Decide` は薄い dispatcher にする。

```text
ManueAgent.Decide
  ├─ legal actions from request.Round
  ├─ win action selection
  ├─ decideSelfTurn
  │   ├─ riichi accepted: tsumogiri discard only
  │   ├─ buildSelfTurnCandidates
  │   ├─ evaluateCandidates
  │   └─ chooseBestCandidate(preferBlack=true)
  └─ decideOtherDiscardReaction
      ├─ buildReactionCandidates
      ├─ evaluateCandidates
      └─ chooseBestCandidate(preferBlack=false)
```

和了 action は CoffeeScript 版同様に常に最優先する。暗槓、加槓、九種九牌など CoffeeScript 版の主要評価経路にない action は、再実装初期では評価対象外にし、合法でも選ばない。後で扱う場合は別途仕様を追加する。

## 3. 候補生成

### 3.1 通常手番

- 合法手から discard 候補を作る。
- 同一牌の手出し discard とツモ切り discard が両方ある場合は、CoffeeScript 版の「選択牌が手牌末尾と `Pai.equal` なら `tsumogiri=true`」に合わせ、候補生成時にツモ切り側へ正規化する。
- 赤牌差分は CoffeeScript 版の `Pai.equal` と同じ扱いにする。赤を黒 5 と勝手に同一視しない。
- legal action に `Riichi` がある場合、通常打牌候補とは別に reach+discard 候補を作る。ただし原仕様どおり、`reachMode="now"` 評価で `shanten <= 0` になる候補だけを残す。
- リーチ宣言後・成立前の打牌は action と trace key は通常打牌のまま、評価だけ `reachMode="now"` 相当にする。

### 3.2 他家打牌反応

- `Pass` を `none` 候補として必ず評価対象に含める。
- `Chii`、`Pon`、`CalledKan` は副露後の手牌と副露リストを作り、喰い替えになる打牌を除いた「鳴いた後の打牌候補」を評価する。
- 最終選択が `none` なら `Pass` を返す。副露候補が勝った場合は副露 action だけを返し、鳴いた後の打牌は次の action opportunity で判断する。

## 4. 候補評価モデル

候補評価は CoffeeScript 版 `getMetricsInternal` と同じ意味の値を作る。Go 内部名は domain 寄りの English にする。

- `safeProb`: 全相手に対して直後放銃しない確率。
- `dealInProb`: `1 - safeProb`。CoffeeScript の `hojuProb`。
- `winProb`: 候補後に自分が和了する確率。CoffeeScript の `horaProb`。
- `averageWinPoints`: 和了時の条件付き平均点。
- `winPointsDist`: 和了点分布。
- `expectedWinPoints`: Monte Carlo 全試行あたりの和了点期待値。
- `exhaustiveDrawProb`: 自分が和了しなかった branch における流局確率を含む最終流局確率。
- `exhaustiveDrawAveragePoints`: 流局時の平均点差。
- `immediateScoreDeltaDist`: 直後放銃または無変化の score delta 分布。
- `selfWinScoreDeltaDist`: 自分和了時の score delta 分布。
- `exhaustiveDrawScoreDeltaDist`: 流局時の score delta 分布。
- `futureScoreDeltaDist`: 自分和了、流局、他家和了を混合した score delta 分布。
- `scoreDeltaDist`: 直後無変化 branch を future 分布で置換した最終 score delta 分布。
- `expectedPoints`: `scoreDeltaDist` における自分の期待点。
- `averageRank`: `scoreDeltaDist` と順位統計から計算した平均順位。
- `shanten`: 候補別の向聴。打牌後再解析ではなく、原仕様の `Goal.ThrowableVector` 由来に合わせる。
- `red`: 赤牌 tie-break 用。

選択規則は平均順位が小さい候補、期待点が大きい候補、赤牌を切らない候補の順にする。副露判断では `preferBlack=false` とし、赤牌 tie-break を使わない。

## 5. 純粋ロジックと依存データ

次の単位へ分ける。

- candidate builder: legal actions と observation から候補を作る。評価値は持たない。
- candidate evaluator: 候補、局面、stats、danger estimator、random source を受け、評価値を返す。
- danger estimator: `round.StateViewer` から scene を作り、danger tree と feature evaluator で放銃確率を返す。
- win estimator: 見えている牌、手牌、副露、`service.AnalyzeShanten` の goal、残り巡目、random source から候補別 win metrics を返す。
- score/rank model: 和了、流局、他家和了、順位期待値の score delta 分布を作る。
- trace formatter: evaluation の表と decided key を文字列化する。I/O は行わない。

既存 `domain/game` のルール判定を再実装しない。特に合法手、向聴、和了形、役、点数、聴牌判定は `round` / `service` の既存 API を使う。AI 側の不足は「評価のための変換」だけに留める。

## 6. 危険度推定

AI package は configs の concrete type に依存せず、danger tree を読むための小さい interface を定義する。

```go
type DangerTreeNode interface {
    LeafProb() (float64, bool)
    Feature() (string, bool)
    NegativeNode() DangerTreeNode
    PositiveNode() DangerTreeNode
}
```

必要なら外側の `configs.DecisionNode` に accessor を追加してこの interface を満たす。依存方向は `configs -> internal/domain/ai` までに留める。

Scene は `round.StateViewer` から構築する。

- 自分の手牌、対象者の河/副露/リーチ状態。
- 安全牌、見えている牌、ドラ。
- 場風、自風。
- リーチ前捨て牌、リーチ宣言牌、早い/遅い捨て牌。

安全牌は shortcut として放銃確率 0 を返す。それ以外は decision tree の feature 評価に従う。相手の聴牌確率は原仕様と同じくリーチ者 1、非リーチ者は stats の `yamitenStats`、欠損時 1 とする。

## 7. 和了推定

- 見えている牌から未見牌 wall を作る。
- stats の残り巡目分布から期待ツモ回数を求める。
- Monte Carlo はまず直列実装にする。並列化が必要になった場合は、worker ごとに候補 accumulator を持ち、最後に merge する。
- `service.AnalyzeShanten` の goal と throwable vector を使い、候補ごとの達成可能性を判定する。
- goal pruning は候補別 shanten ではなく、打牌前/副露後の base shanten を基準にする。
- 候補別 shanten は原仕様どおり、打牌後の再解析ではなく、その牌を捨てられる goal の最小 shanten とする。`none` は base shanten。
- 点数は現行 `service.CalculateFuHan` と `service.RonPoints` を使う。原仕様の簡易役計算との差が出る場合は、差分を characterization で確認し、意図的な Go 側改善か互換性問題かを記録する。
- 赤ドラは原仕様と同様、候補 hand / meld 内の赤牌のうち goal 構成牌に同種牌が含まれるものだけを scoring 対象にする。

## 8. Configs と CLI 接続

- `cmd/mjai-manue` など外側で `configs.LoadGameStats()` と `configs.LoadDangerTree()` を呼ぶ。
- `ai.NewManueAgent(seed, deps)` は deps validation を行い、不足や不整合があれば error を返す。
- `cmd/mjai-manue` は deps 構築エラーを stderr に出し、runtime error として終了する。
- runtime、application、adapter の既存動線は変更しない。
- `--seed` は Agent の再現性に使う。CoffeeScript の `seedRandom("")` 乱数列との完全一致は狙わない。

## 9. Trace とログ

Agent は stdout/stderr へ直接書かない。評価表、tenpai probabilities、`decidedKey` 相当は `Decision.Trace` や application の decision metadata として返す。

trace table は原仕様の列名に寄せる。

- `action`
- `avgRank`
- `expPt`
- `hojuProb`
- `myHoraProb`
- `ryukyokuProb`
- `otherHoraProb`
- `avgHoraPt`
- `ryukyokuAvgPt`
- `shanten`

wire JSON の `log` field と stderr trace は混ぜない。action golden test では trace/log を比較対象にしない。

## 10. テスト仕様

- 純粋関数は `internal/domain/ai` でテーブル駆動テストにする。
- 候補比較は平均順位、期待点、赤牌 tie-break の順序を固定する。
- 危険度 feature は feature ごとに小さい scene fixture を作る。
- 確率分布、流局点、順位期待値は小さい人工 stats で期待値を固定する。
- Agent 判断は private field を直接壊さず、mjai JSON Lines から `round.State` を構築して action を比較する。
- runtime golden は action のみ比較する。trace/log は別 fixture にする。
- Characterization は `docs/manue-ai-original-spec.md` のケース候補を優先する。
- 最終確認は PowerShell で `$env:GOEXPERIMENT='jsonv2'; go test ./...` を実行する。
