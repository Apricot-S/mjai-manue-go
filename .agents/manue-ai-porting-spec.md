# Manue AI Porting Specification

この文書は CoffeeScript 版 Manue AI を Go へ移植するための設計仕様である。原仕様は `.agents/manue-ai-original-spec.md` を一次資料とし、本書では現行 `domain/game` と `application` へ接続するための責務分割だけを定める。

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
- リーチ宣言後・成立前の打牌は、`round.State.LegalActions` がテンパイ維持できる立直宣言牌だけを列挙する前提にする。AI candidate builder は追加のシャンテン再検証で候補を除外しない。候補の action と trace key は通常打牌のまま、評価だけ `scoreAsRiichi=true` として扱う。

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

### 5.1 現行 Manue 固有コードの扱い

現行 `internal/domain/ai` の Manue 固有コードは補助資料として扱う。再実装時は、次の分類を初期判断にする。流用は仕様に合う場合だけ行い、既存関数の都合で設計を曲げない。

| 分類 | 対象 | 扱い |
| --- | --- | --- |
| 維持 | `agent.go`, `tsumogiri_agent.go` | Manue 固有再実装の対象外。Agent interface と TsumogiriAgent は残す。 |
| 再採用候補 | `candidate_score.go` | 候補比較、score delta 合成、平均順位/期待点への変換は純粋関数として使える。命名と入力 struct は最終設計に合わせて調整する。 |
| 再採用候補 | `scalar_prob_dist.go`, `score_delta_prob_dist.go`, `ahead_vector_prob_dist.go` | 分布演算は純粋関数として使える。正規化や非正値 drop の仕様は characterization と合わせて維持する。 |
| 再採用候補 | `win_score.go`, `deal_in_score.go`, `round_end_prob.go`, `round_end_score.go`, `exhaustive_draw_score.go`, `rank.go` | score/rank model の部品として使える。stats interface との接続は新設計側で決める。 |
| 再採用候補 | `win_estimate_score.go`, `win_trials.go` | Monte Carlo trial 集計と wall/trial helper は純粋関数として使える。乱数列の完全互換は要求しない。 |
| 再採用候補 | `win_goals.go`, `win_estimator.go` | `service.CalculateFuHan` は CoffeeScript 版 `calculateFan` 相当の移植済み実装として使う。代表ケースは `service` 側と AI 側の goal scoring テストで固定する。 |
| 条件付き再採用 | `danger.go`, `danger_scene.go`, `danger_tiles.go` | danger tree interface と feature evaluator は使えるが、feature 名ごとの原仕様対応を確認しながら接続する。 |
| 再設計対象 | `manue_agent.go`, `self_turn_candidates.go`, `call_candidates.go`, `evaluator.go`, `trace.go` | orchestration と Agent 境界は新設計で組み直す。既存実装はテスト fixture と細部確認の補助に留める。 |
| 再設計対象 | `deps.go`, `stats.go` | deps/stats interface と validation は configs/CLI 接続時に再確認する。外側 concrete type に依存しない方針は維持する。 |

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
- 点数は現行 `service.CalculateFuHan` と `service.RonPoints` を使う。`CalculateFuHan` は CoffeeScript 版 `calculateFan` 相当の移植済み実装であり、正式な麻雀点数計算へ置き換えない。代表ケースは `service` 側の `TestCalculateFuHan` と、AI 側の goal scoring テストで固定する。
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
- Characterization は `.agents/manue-ai-original-spec.md` のケース候補を優先する。
- 最終確認は PowerShell で `$env:GOEXPERIMENT='jsonv2'; go test ./...` を実行する。

### 10.1 Characterization の移植先

`.agents/manue-ai-original-spec.md` の OC-* は、Go 側では次の粒度へ落とす。

| Original ID | Go 側テスト粒度 | 比較対象 |
| --- | --- | --- |
| OC-001, OC-002, OC-003 | Agent dispatcher | action と error |
| OC-004, OC-005, OC-006, OC-007 | candidate builder | candidate key、action、riichi/scoreAsRiichi、discard tile |
| OC-008 | domain/game LegalActions boundary | `TestState_LegalActions_RiichiDeclarationTileKeepsTenpai` で検証済み。AI 側では再検証しない。 |
| OC-009, OC-010, OC-011 | reaction candidate builder + Agent 判断 | candidate key、返却 action |
| OC-012, OC-013, OC-014 | danger/deal-in evaluator | `safeProb`、`dealInProb`、`immediateScoreDeltaDist` |
| OC-015, OC-016 | win estimator | candidate shanten、win probability、points distribution |
| OC-017, OC-018 | round-end model | exhaustive draw probability、tenpai probability、score delta distribution |
| OC-019 | rank model | average rank |
| OC-020 | candidate comparator | selected candidate key |
| OC-021 | trace formatter | trace/log string |

Agent 判断の characterization は、最終的に mjai JSON Lines から `round.State` を作って action を比較する。純粋関数で十分に固定できるものは、CoffeeScript 実行 fixture へ寄せず、Go の table-driven test で期待値を明示する。

### 10.2 仕様固定フェーズの coverage

実装前の最低ラインは次の状態とする。`Covered` は unit test で Go 側の責務境界を固定済みであることを表す。

| ID | Status | 既存/追加予定テスト | 備考 |
| --- | --- | --- | --- |
| OC-001 | Covered | `TestManueAgent_DecideActionSkeleton` | 和了 action が評価前に優先される入口を固定する。 |
| OC-002 | Covered | `TestManueAgent_DecideActionSkeleton` | リーチ accepted 中の自摸はツモ切り discard を返す。 |
| OC-006 | Covered | `TestBuildSelfTurnCandidates_BuildsRiichiCandidate`, `TestBuildSelfTurnCandidates_FiltersNonTenpaiRiichiCandidate` | reach action ありの候補生成を固定する。 |
| OC-009 | Covered | `TestBuildReactionCandidates_BuildsPassAndCallDiscards`, `TestManueAgent_decideOtherDiscardReaction_EvaluatesCallCandidates` | `none` と副露後打牌候補を同じ評価経路へ載せる。 |
| OC-011 | Covered | `TestBuildReactionCandidates_BuildsPassAndCallDiscards` | 喰い替え候補が除外されることを副露候補生成で固定する。 |
| OC-012 | Covered | `TestCandidateEvaluator_dealInEstimates_SafeTileHasZeroDealInProb` | danger scene 単体ではなく、deal-in evaluator 経由で安全牌 shortcut を固定する。 |
| OC-015 | Covered | `TestCandidateShantenUsesThrowableVector`, `TestCandidateShantenReturnsBaseForNone`, `TestCandidateShantenReturnsInfinityWhenTileIsNotThrowable` | 候補別 shanten が打牌後再解析ではなく throwable vector 由来であることを固定する。 |
| OC-017 | Covered | `TestExhaustiveDrawProb`, `TestExhaustiveDrawProbOnSelfNoWin` | 流局確率と自分非和了条件付き補正を固定する。 |
| OC-020 | Covered | `TestCompareCandidateScore`, `TestChooseBestCandidate_PrefersBlackTileOnTie`, `TestChooseBestCandidate_DoesNotPreferRiichiOnScoreTie` | 平均順位、期待点、赤牌、走査順の tie-break を固定する。 |

OC-008 は `domain/game` の `TestState_LegalActions_RiichiDeclarationTileKeepsTenpai` で境界を固定済みのため、AI 側では追加しない。

最低ライン外の OC は次の純粋関数/trace テストで固定済みとして扱う。これらは実装順の各段階で回帰防止に使うが、Agent action golden の必須条件にはしない。

| ID | Status | 既存テスト | 備考 |
| --- | --- | --- | --- |
| OC-003 | Covered | `TestManueAgent_decideSelfTurn_ReturnsErrorWithoutTsumogiriAfterRiichiAccepted` | リーチ accepted 中にツモ切りがない異常系を固定する。 |
| OC-004 | Covered | `TestNormalizedSelfTurnDiscards_PrefersTsumogiriForSameTile` | 同一牌候補の正規化を固定する。 |
| OC-005 | Covered | `TestNormalizedSelfTurnDiscards_PrefersTsumogiriForSameTile`, `TestChooseBestCandidate_PrefersBlackTileOnTie` | 赤 5 と黒 5 を別候補として残し、比較では黒優先を固定する。 |
| OC-007 | Covered | `TestBuildSelfTurnCandidates_IncludesRiichiAndDiscardCandidates`, `TestFilteredWinEstimateGoals` | reach 候補と通常 discard 候補が別系統で評価されることを固定する。 |
| OC-010 | Covered | `TestBuildReactionCandidates_BuildsPassAndCallDiscards`, `TestManueAgent_decideOtherDiscardReaction_EvaluatesCallCandidates` | 副露後打牌は評価用候補であり、選択 action は副露 action のままになることを固定する。 |
| OC-013 | Covered | `TestTenpaiProb_ReturnsOneWithRiichi`, `TestTenpaiProb_ReturnsYamitenRatio`, `TestTenpaiProb_ReturnsOneWithoutStats` | リーチ者、統計あり非リーチ者、統計欠損の聴牌確率を固定する。 |
| OC-014 | Covered | `TestImmediateScoreDeltaDist`, `TestImmediateScoreDeltaDistFromStats`, `TestImmediateScoreDeltaDist_ReturnsNoChangeWithoutDealInDists` | 直後放銃分布と無変化 branch 置換を固定する。 |
| OC-016 | Covered | `TestTrialWinPts`, `TestCandidateTrialWinPts` | 複数 goal 達成時に候補ごとの最大 points を採用することを固定する。 |
| OC-018 | Covered | `TestExhaustiveDrawTenpaiProbs`, `TestExhaustiveDrawScoreDeltaDist` | 流局時聴牌確率と 3000 点授受分布を固定する。 |
| OC-019 | Covered | `TestAverageRank`, `TestWinProbAgainst`, `TestWinProbFromRelativeScore_UsesStatsWhenAvailable`, `TestWinProbFromRelativeScore_FallsBackToStartingDealerOrder` | 点差分布から順位期待値へ変換する rank model を固定する。 |
| OC-021 | Covered | `TestFormatCandidateTrace`, `TestFormatCandidateTrace_FormatsInfinityShanten`, `TestFormatDecisionTrace_AppendsDecidedKey` | trace/log は action golden とは分けて文字列形を固定する。 |

### 10.3 原仕様との差分記録が必要な箇所

次は移植時に差分が出やすいため、実装前または最初のテスト追加時に意図を明記する。

- 点数計算: Go 側の `service.CalculateFuHan` は CoffeeScript 版 `calculateFan` 相当の移植済み実装として扱う。正式な麻雀点数計算との差ではなく、オリジナル互換性を `TestCalculateFuHan` と AI 側 goal scoring テストで固定する。
- 乱数: CoffeeScript 版は `seedRandom("")` で呼び出しごとに同じ乱数列を使う。Go 側は `--seed` と `math/rand/v2` による決定性を優先し、乱数列の完全一致は要求しない。
- trace/log: CoffeeScript 版は `console.log` と `@log` が混在する。Go 側は stdout へ直接出さず、`Decision.Trace` / metadata として返す。
- legal actions: CoffeeScript 版の対象外 action は初期実装では選ばない。暗槓、加槓、九種九牌などを扱う場合は別仕様を追加する。
- state 管理: CoffeeScript 版は AI 内部の game/player 参照が広い。Go 側では `round.ActionStateViewer` / `round.StateViewer` を observation として使い、AI 側に重複 state を持たない。
- 立直宣言牌: CoffeeScript 版は AI 側の `shanten <= 0` filtering で宣言牌を絞るが、Go 側は `domain/game` の `LegalActions` が合法な立直宣言牌だけを列挙する。AI 側で同じシャンテン計算を防御的に繰り返さない。
