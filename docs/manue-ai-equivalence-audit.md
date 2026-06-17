# Manue AI Semantic Equivalence Audit

対象: CoffeeScript 版 `mjai-manue` を正とした、Go `develop` ブランチの Manue AI 判断ロジック監査。

この文書は実装変更ではなく監査レポートである。コード形状や関数名ではなく、候補生成、評価式、条件分岐、優先順位、境界条件、設定値の意味を比較する。

## 1. 監査時点

- 監査ブランチ: `develop`
- 監査時点の作業ツリー: `?? zzz_list.txt`
- CoffeeScript 一次資料: `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee`
- CoffeeScript 補助資料: `ai.coffee`, `game.coffee`, `pai.coffee`, `furo.coffee`, `prob_dist.coffee`, `danger_estimator.coffee`, `tenpai_prob_estimator.coffee`, `shanten_analysis.coffee`
- Go 一次対象: `internal/domain/ai/*.go` の production file 全体
- Go 補助対象: `internal/domain/game/round` の `LegalActions` / viewer、`internal/domain/game/round/service` の向聴・点数・役判定、`cmd/mjai-manue` の seed/config 接続
- 補助参照: `reference/repositories/mjai-manue-go-main/`。判断根拠ではなく、過去移植の差分理解に限る。
- 検証コマンド: `$env:GOEXPERIMENT='jsonv2'; go test ./...`
- 検証結果: pass

## 2. CoffeeScript 原仕様の分解

`ManueAI.respondToAction(action)` は次の意味単位に分解できる。

```text
respondToAction
  ├─ 自分の tsumo/chi/pon/reach
  │  ├─ categorizeActions(possibleActions)
  │  ├─ hora があれば即 createAction(hora)
  │  ├─ reach accepted 中の tsumo は評価せずツモ切り
  │  └─ decideDahai
  │     ├─ getMetrics
  │     │  ├─ 手牌から重複除去済み candDahais を作る
  │     │  ├─ canReach あり: reachMode="now" の聴牌候補 + reachMode="never" の通常候補
  │     │  └─ canReach なし: reachDeclared なら reachMode="now"、通常は "default"
  │     ├─ printMetrics / printTenpaiProbs
  │     └─ chooseBestMetric(preferBlack=true)
  └─ 他家の dahai/kakan
     ├─ categorizeActions(possibleActions)
     ├─ hora があれば即 createAction(hora)
     └─ decideFuro
        ├─ none 候補を getMetricsInternal で評価
        ├─ 各副露候補を仮適用し、喰い替え牌を除いた打牌候補を評価
        ├─ printMetrics / printTenpaiProbs
        └─ chooseBestMetric(preferBlack=false)
```

中心評価は `getMetricsInternal(tehais, furos, candDahais, reachMode)` である。候補ごとに危険度、直後放銃分布、自分和了推定、流局、他家和了、順位期待値を合成し、最終的に `averageRank` 昇順、`expectedPoints` 降順、必要な場合だけ黒牌優先、最後は走査順で選ぶ。

## 3. 全体判定

総合判定: **概ね設計上の意味等価を目指した実装だが、CoffeeScript 版との完全な action 等価は未証明**。

- `respondToAction` 相当の dispatch、和了優先、リーチ accepted ツモ切り、自己手番候補、副露反応候補、評価値合成、比較優先順位は Go 側に対応実装がある。
- `possible_actions` 非依存、`LegalActions` への責務移動、乱数列非互換、stdout/stderr 分離は既存設計で許容された差分であり、単体では不一致扱いしない。
- 動的に CoffeeScript と同一局面を実行して action を比較する characterization fixture は、この監査時点では不足している。特に danger feature 全量、点数計算の全役組み合わせ、Monte Carlo wall/乱数差の action 影響は未証明である。

## 4. 対応表

判定値:

- `Equivalent`: 意味が原仕様と一致しているとコードとテストから判断できる。
- `Equivalent by design difference`: 形は異なるが、設計で許容された差分として意味が維持される。
- `Different`: CoffeeScript と意味が異なり、action 不一致の可能性がある。
- `Missing in Go`: CoffeeScript の意味単位に対応する Go 実装がない。
- `Go-only`: Go 側だけの境界、検証、補助機能。
- `Needs dynamic check`: 静的読解と単体テストだけでは等価性を断定できない。

| CoffeeScript location | Original behavior | Go location | Go behavior | 判定 | 根拠 |
| --- | --- | --- | --- | --- | --- |
| `manue_ai.coffee:respondToAction` | 自分の `tsumo`/`chi`/`pon`/`reach` と他家 `dahai`/`kakan` だけを判断入口にする。 | `manue_agent.go:Decide` | `LegalActions` がある局面で、和了、自己手番、他家反応へ dispatch する。 | Equivalent by design difference | Go は入力 event 種別ではなく `round.State` と合法手から行動機会を判断する。 |
| `categorizeActions` | `hora` と `reach` は最後に見つかったもの、その他は `furos`。 | `firstActionOfType`, `LegalActions` | 型付き action から先頭の `Win` / `Riichi` / `Pass` などを選ぶ。 | Equivalent by design difference | Go の合法手生成側が重複や順序を制御する前提。複数和了/立直候補の順序差は動的確認対象。 |
| `respondToAction` hora branch | 和了可能なら評価せず即選択。 | `manue_agent.go:Decide` | `Win` があれば即 `Decision{Action: win}`。 | Equivalent | OC-001 相当が `TestManueAgent_DecideActionSkeleton` で固定済み。 |
| `respondToAction` reach accepted branch | リーチ成立中の自摸は評価せず `tsumogiri: true`。 | `decideSelfTurn`, `tsumogiriDiscard` | `RiichiAccepted` ならツモ切り discard のみ返し、暗槓は無視。 | Equivalent | OC-002/003 相当が固定済み。 |
| `getMetrics` candDahais | 手牌から `Pai.equal` 重複と `cannotDahai` を除外。 | `buildSelfTurnCandidates`, `normalizedSelfTurnDiscards`, `round.LegalActions` | 合法 discard を入力にし、同一 `tile.Tile` はツモ切り側へ正規化。 | Equivalent by design difference | `cannotDahai` は Go では合法手列挙境界へ移動。赤牌は `tile.Tile` として別候補。 |
| `getMetrics` reach branch | `canReach` 時は `0.<pai>` reach 候補と `-1.<pai>` 通常候補を併存。 | `buildSelfTurnCandidates` | `Riichi` action があれば `riichi=true` 候補と通常 discard 候補を作る。 | Equivalent | OC-006/007 相当が固定済み。 |
| `selectTenpaiMetrics` | `reachMode="now"` は `shanten <= 0` の候補だけ残す。 | `candidateShanten`, `filteredWinEstimateGoals`, `buildSelfTurnCandidates` | `scoreAsRiichi` かつ `shanten <= 0` の候補だけ reach 候補化。 | Equivalent | `ThrowableVector` 由来の shanten テストあり。 |
| `reachDeclared` | 立直宣言後は `reachMode="now"` で聴牌候補だけ残す。 | `buildSelfTurnCandidates` | `LegalActions` が立直宣言牌を絞る前提で、AI は `scoreAsRiichi=true` の評価だけ行う。 | Equivalent by design difference | `domain/game` 境界へ責務移動。AI 側では再検証しない設計。 |
| `decideDahai` | 最善 key が `0.*` なら reach、そうでなければ discard。 | `buildSelfTurnCandidate`, `buildCandidateDecision` | selected candidate の immediate action を返す。 | Equivalent | 候補 action が `Riichi` または `Discard` に固定される。 |
| `decideDahai` tsumogiri flag | 元 event が `tsumo`/`reach` かつ選択牌が手牌末尾と等しければ `tsumogiri=true`。 | `normalizedSelfTurnDiscards` | 同一牌に手出し/ツモ切りがあればツモ切り action を残す。 | Equivalent by design difference | Go は出力時ではなく候補正規化時に決める。 |
| `decideFuro` none | `none` を副露候補と同じ比較器で評価。 | `buildReactionCandidates` | `Pass` を `traceKey="none"` の評価候補にする。 | Equivalent | OC-009 相当が固定済み。 |
| `decideFuro` call apply | `consumed` を手牌から削除し `Furo` を追加して評価。 | `actionCallMeld`, `baseHand.Call`, `buildCallReactionCandidates` | domain の meld/hand API で副露後手牌と melds を作る。 | Equivalent by design difference | Go は domain invariant を使う。 |
| `isKuikae` | consumed+dahai が刻子同種または同色連続なら除外。 | `callSwapTiles`, `isSwapCallTile`, `meld.SwapCallTiles` | meld 側が算出した喰い替え牌と同種牌を除外。 | Needs dynamic check | テストはあるが、CoffeeScript の `isKuikae` 全入力との網羅比較は未実施。 |
| `decideFuro` selected action | 副露候補が勝っても後続打牌は返さず副露 action だけ返す。 | `buildCallReactionCandidates` | 評価 candidate は打牌 tile を持つが、返却 action は call action。 | Equivalent | OC-010 相当が固定済み。 |
| `getMetricsInternal` | safe、直後放銃、自分和了、流局、他家和了を合成。 | `candidateEvaluator`, `evaluateCandidateFromComponents` | 同じ構成要素を context 化し、候補ごとに score を作る。 | Equivalent | 部品単位のテストあり。 |
| `getSafeProbs` | 相手ごとに `1 - tenpaiProb * dangerProb` を掛ける。 | `dealInEstimates`, `safeProb` | 相手ごとの `tenpai * rawProb` から安全確率を積算。 | Equivalent | OC-012/013 と `TestSafeProb` で固定済み。 |
| `scene.anpai` shortcut | 現物なら対象相手への放銃確率 0。 | `DecisionTreeDangerEstimator.EstimateDealInProb` | `state.SafeTiles(winner)` に同種牌があれば 0。 | Equivalent | `TestDecisionTreeDangerEstimator_SafeTileSkipsSceneBuild` あり。 |
| `danger_estimator.coffee` features | danger tree feature を scene から評価。 | `danger_scene.go`, `danger_tiles.go` | feature 名を switch/parse して評価。 | Needs dynamic check | 代表テストはあるが、config 内 feature 全量と CoffeeScript の scene 構築全量比較は未完了。 |
| `getImmediateScoreChangesDists` | 相手ごとに無変化 branch を置換し、最初のロンだけ考慮。 | `immediateScoreDeltaDist` | `scoreDelta{}` branch を順に replace。 | Equivalent | OC-014 相当が固定済み。 |
| `getHoraEstimation` goal pruning | `reachMode="now"` と `analysis.shanten()>3` の pruning。 | `filteredWinEstimateGoals` | `scoreAsRiichi` と `baseShanten > 3` で同等 pruning。 | Equivalent | `TestFilteredWinEstimateGoals` あり。 |
| `calculateFan` | 簡易役/飜/符/ロン点を算出。 | `scoredWinEstimateGoals`, `service.CalculateFuHan`, `service.RonPoints` | domain service の移植済み簡易点数を使う。 | Needs dynamic check | 代表テストはあるが `calculateFan` 全役組み合わせとの直接比較は未完了。 |
| `getHoraEstimation` visible wall | 自分から見えている牌を全牌集合から除いて未見牌を作る。 | `winEstimatesFromState`, `unseenWallFromVisibleTiles` | `state.VisibleTiles(self)` から 34 種 4 枚の wall を作る。 | Equivalent by design difference | Go は赤を `RemoveRed` して 34 種へ集約し、過剰 visible を error にする。 |
| `getHoraEstimation` random | 呼び出しごとに `seedRandom("")` 固定で 1000 試行。 | `candidateEvaluator`, `winEstimatesFromShuffledWall` | Agent seed の PCG を使い 1000 試行。 | Equivalent by design difference | 設計で乱数列完全一致は非要求。ただし action 影響は動的確認対象。 |
| `getHoraEstimation` best goal | 同一候補で複数 goal 達成時は最大 points。 | `trialWinPts`, `candidateTrialWinPts` | 達成 goal の最大 points を採用。 | Equivalent | OC-016 相当が固定済み。 |
| `getHoraEstimation` zero wins | `averageHoraPoints` は 0 除算特別扱いなし。 | `newWinEstimate` | 無和了なら平均点 0、空分布を返す。 | Different | Go は NaN/Infinity を避ける防御的仕様。action 影響は要確認。 |
| `getRyukyokuProb` | 残り turn 分布で `ryukyokuRatio / den`。 | `exhaustiveDrawProb` | turn 分布の残り質量で `ExhaustiveDrawRatio()/remainingProb`。 | Equivalent | OC-017 相当が固定済み。 |
| `getRyukyokuProbOnMyNoHora` | `getRyukyokuProb()^(3/4)`。 | `exhaustiveDrawProbOnSelfNoWin` | `common.NumPlayers` から `3/4` を算出。 | Equivalent | テストあり。 |
| `getNotenRyukyokuTenpaiProb` | `turn + 1/4` から最終巡目までの tenpai freq。 | `notenExhaustiveDrawTenpaiProb` | 現在 turn より未来の freq と noten count で算出。 | Equivalent | OC-018 相当が固定済み。 |
| `getScoreChangesDistOnRyukyoku` | 4 人の聴牌 Bernoulli を合成して 3000 点授受分布。 | `exhaustiveDrawScoreDeltaDist` | `aheadVectorProbDist` 経由で同等分布を作る。 | Equivalent | `TestExhaustiveDrawScoreDeltaDist` あり。 |
| `getRandomHoraScoreChangesDist` / `getHoraFactorsDist` | 自摸/ロン比率と親子点数 freq から他家和了分布。 | `randomWinScoreDeltaDist`, `winScoreFactorDist` | stats の `NumSelfDrawWins/NumWins` と点数 freq から分布化。 | Equivalent | win score tests あり。 |
| `getAverageRank` / `getWinProb` | 点差分布を次局勝率表に写し、独立 Bernoulli で平均順位。 | `rank.go` | relative score dist と `RelativeWinProbs` から平均順位を計算。 | Equivalent | OC-019 相当が固定済み。 |
| `getWinProbFromRelativeScore` | stats 欠損時は起家から近い方が同点有利。 | `winProbFromRelativeScore` | stats 欠損時に starting dealer からの位置で fallback。 | Equivalent | fallback テストあり。 |
| `chooseBestMetric` / `compareMetric` | 平均順位、期待点、黒牌、走査順。 | `chooseBestCandidate`, `compareCandidates` | 同じ比較順。 | Equivalent | OC-020 相当が固定済み。 |
| `printMetrics` | 常に `preferBlack=true` で並べた表を `@log`。 | `formatCandidateTrace` | 常に `sortedCandidates(..., true)` で表を作る。 | Equivalent | trace tests あり。 |
| `printTenpaiProbs` | 相手 ID と聴牌確率を `@log`。 | `formatTenpaiProbsTrace` | self 以外を ID 順で出力。 | Equivalent | trace tests あり。 |
| `console.log("goals")`, `console.log("decidedKey")` | stdout に診断出力。 | `formatDecisionTrace`, `Decision.Trace` | stdout へ直接書かず trace として返す。 | Equivalent by design difference | I/O 安全性の設計差分。 |
| CoffeeScript constructor | JSON stats と danger tree を直接 load。 | `deps.go`, `NewManueAgent`, `cmd/mjai-manue` | stats/danger は外側で load し deps として注入。 | Equivalent by design difference | Clean Architecture の依存方向を維持。 |
| `TenpaiProbEstimator.estimate` | リーチなら 1、yamiten stats 欠損なら 1。 | `tenpaiProb` | 同じ fallback。 | Equivalent | OC-013 相当が固定済み。 |
| `ProbDist` | scalar/vector 分布の add/mult/merge/replace/expected。 | `scalar_prob_dist.go`, `score_delta_prob_dist.go`, `ahead_vector_prob_dist.go` | 型別の分布演算として実装。 | Equivalent | 各分布演算テストあり。 |
| `ShantenAnalysis(... allowedExtraPais: 1)` | 候補別 shanten と goal を作る。 | `service.AnalyzeShanten(... AllowedExtraTiles(1))`, `candidate_builder.go` | domain service を使用し、候補 shanten は throwable vector 由来。 | Equivalent by design difference | shanten service の完全移植性は本監査の補助対象。 |
| CoffeeScript 対象外 action | 暗槓、加槓、九種九牌は主要評価経路にない。 | `buildSelfTurnCandidates` | concealed/promoted kan と kyushukyuhai を選ばない。 | Equivalent | `TestBuildSelfTurnCandidates_IgnoresConcealedAndPromotedKan` あり。 |
| Go Agent interface | なし。 | `agent.go` | `Request`, `Decision`, `Agent` の境界。 | Go-only | Go アーキテクチャ上の境界。 |
| Go tsumogiri agent | CoffeeScript には別 `tsumogiri_ai.coffee`。 | `tsumogiri_agent.go` | Manue 監査対象外の最小 AI。 | Go-only | Manue 判断ロジックではない。 |
| Go stats validation | CoffeeScript は load 後に構造 validation しない。 | `stats.go` | stats の構造不整合を起動時 error にする。 | Go-only | 防御的境界。正常 stats では action 意味に影響しない。 |

## 5. Go production file 対応一覧

| Go file | 主な対応先 | 判定 |
| --- | --- | --- |
| `agent.go` | Agent 境界 | Go-only |
| `manue_agent.go` | `respondToAction`, `decideDahai`, `decideFuro` orchestration | Equivalent by design difference |
| `candidate.go` | CoffeeScript metric key と候補の暗黙 state | Equivalent by design difference |
| `candidate_builder.go` | 候補別 shanten / `ShantenAnalysis` goal usage | Equivalent |
| `self_turn_candidates.go` | `getMetrics`, `mergeMetrics`, `selectTenpaiMetrics` | Equivalent |
| `call_candidates.go` | `decideFuro`, `isKuikae` | Needs dynamic check |
| `candidate_score.go` | `getMetricsInternal`, `chooseBestMetric`, `compareMetric` | Equivalent |
| `evaluator.go` | `getMetricsInternal` の合成全体 | Equivalent |
| `deal_in_score.go` | `getSafeProbs`, `getImmediateScoreChangesDists` | Equivalent |
| `danger.go` | `DangerEstimator.estimateProb`, `scene.anpai` shortcut | Equivalent |
| `danger_scene.go` | `DangerEstimator.Scene` feature evaluation | Needs dynamic check |
| `danger_tiles.go` | danger feature helper | Needs dynamic check |
| `win_goals.go` | `getHoraEstimation` goal filtering / `calculateFan` 接続 | Needs dynamic check |
| `win_trials.go` | Monte Carlo trial goal achievement | Equivalent |
| `win_estimator.go` | `getHoraEstimation` trial loop | Equivalent by design difference |
| `win_estimate_score.go` | `totalHoraVector`, `totalPointsVector`, points dist 集計 | Different |
| `win_score.go` | `getScoreChangesDistOnHora`, `getRandomHoraScoreChangesDist` | Equivalent |
| `round_end_prob.go` | `getRyukyokuProb`, `getNumExpectedRemainingTurns` | Equivalent |
| `round_end_score.go` | `futureScoreChangesDist` merge | Equivalent |
| `exhaustive_draw_score.go` | `getTenpaiProb`, `getScoreChangesDistOnRyukyoku` | Equivalent |
| `rank.go` | `getAverageRank`, `getWinProb` | Equivalent |
| `scalar_prob_dist.go` | `ProbDist` scalar distribution | Equivalent |
| `score_delta_prob_dist.go` | `ProbDist.replace`, vector distribution | Equivalent |
| `ahead_vector_prob_dist.go` | 流局/順位用 Bernoulli 合成 | Equivalent |
| `trace.go` | `printMetrics`, `printTenpaiProbs`, `console.log` 相当 | Equivalent by design difference |
| `deps.go` | constructor の stats/danger deps | Equivalent by design difference |
| `stats.go` | stats validation | Go-only |
| `tsumogiri_agent.go` | Manue 外の最小 AI | Go-only |

## 6. 意味差分一覧

| 重要度 | 判定 | 差分 | action 不一致リスク | 備考 |
| --- | --- | --- | --- | --- |
| High | Different | Go の `newWinEstimate` は無和了試行時に平均点 0 と空分布を返すが、CoffeeScript は `0 / 0` 相当を特別扱いしない。 | 中 | NaN が CoffeeScript の比較や trace にどう効くか、実局面で要動的確認。Go の方が防御的。 |
| High | Needs dynamic check | danger feature 全量が CoffeeScript の `danger_estimator.coffee` と完全一致するか未確認。 | 高 | danger tree feature 名が 1 つでも逆なら候補選択が変わる。config 内 feature 全列挙と比較が必要。 |
| High | Needs dynamic check | `service.CalculateFuHan` が CoffeeScript `calculateFan` の全役・赤ドラ・ドラ条件と完全一致するか未確認。 | 高 | 和了推定点数が変わると期待点/順位が変わる。 |
| Medium | Equivalent by design difference | Go は `possible_actions` を読まず `LegalActions` で合法手を算出する。 | 中 | 設計上正しい差分だが、CoffeeScript がサーバ由来 `possibleActions` の順序/重複に依存する局面では差が出る可能性がある。 |
| Medium | Equivalent by design difference | Go は PCG seed を Agent lifetime で使い、CoffeeScript は `getHoraEstimation` 呼び出しごとに `seedRandom("")` を使う。 | 中 | 乱数列一致は非要求。ただし reach/通常など比較条件間の乱数相関は CoffeeScript と異なる可能性がある。 |
| Medium | Needs dynamic check | 副露喰い替え判定は Go では meld 側 `SwapCallTiles` に委譲する。 | 中 | 代表テストはあるが、CoffeeScript `isKuikae` の 3 枚判定との全パターン比較は未実施。 |
| Low | Go-only | stats validation が不正 stats を起動時 error にする。 | 低 | 正常 config では影響なし。 |
| Low | Equivalent by design difference | trace/log は stdout へ直接出さず `Decision.Trace` / `Log` に分離。 | 低 | protocol output 安全性のための設計差分。 |

## 7. 未判定 / 追加調査一覧

- CoffeeScript と Go を同一 mjson 局面で実行し、final action を比較する characterization fixture が不足している。
- `configs/danger_tree.all.json` など実際の danger tree に含まれる feature 名を全列挙し、Go `dangerScene.evaluate` がすべて処理し、CoffeeScript と真偽一致するか確認する必要がある。
- `service.CalculateFuHan` と CoffeeScript `calculateFan` について、全対象役、鳴き/門前差、ドラ、赤ドラ、符 30/40、満貫以上丸めの直接比較が必要である。
- Monte Carlo の乱数差が action に影響しやすい近接候補局面では、Go の deterministic seed で期待 action を別途固定する必要がある。
- `LegalActions` への責務移動により、CoffeeScript の `possibleActions` 順序や「最後に見つかった hora/reach」挙動との差が出ないか、複数候補局面で確認する必要がある。
- Go の無和了時平均点 0 化が、CoffeeScript の NaN/Infinity 的挙動と比較器上どの程度違うか、該当局面を作って確認する必要がある。

## 8. テスト対応表

| Original ID | Go 側の主な固定箇所 | 状態 |
| --- | --- | --- |
| OC-001 | `TestManueAgent_DecideActionSkeleton` | 固定済み |
| OC-002 | `TestManueAgent_DecideActionSkeleton` | 固定済み |
| OC-003 | `TestManueAgent_decideSelfTurn_ReturnsErrorWithoutTsumogiriAfterRiichiAccepted` | 固定済み |
| OC-004 | `TestNormalizedSelfTurnDiscards_PrefersTsumogiriForSameTile` | 固定済み |
| OC-005 | `TestNormalizedSelfTurnDiscards_PrefersTsumogiriForSameTile`, `TestChooseBestCandidate_PrefersBlackTileOnTie` | 固定済み |
| OC-006 | `TestBuildSelfTurnCandidates_BuildsRiichiCandidate`, `TestBuildSelfTurnCandidates_FiltersNonTenpaiRiichiCandidate` | 固定済み |
| OC-007 | `TestBuildSelfTurnCandidates_IncludesRiichiAndDiscardCandidates`, `TestFilteredWinEstimateGoals` | 固定済み |
| OC-008 | `domain/game` の legal action テスト | AI 側では設計上対象外 |
| OC-009 | `TestBuildReactionCandidates_BuildsPassAndCallDiscards`, `TestManueAgent_decideOtherDiscardReaction_EvaluatesCallCandidates` | 固定済み |
| OC-010 | `TestBuildCandidateDecision_ReturnsCallActionForWinningReactionCandidate` | 固定済み |
| OC-011 | `TestBuildReactionCandidates_BuildsPassAndCallDiscards` | 代表固定済み、全パターンは未確認 |
| OC-012 | `TestCandidateEvaluator_dealInEstimates_SafeTileHasZeroDealInProb` | 固定済み |
| OC-013 | `TestTenpaiProb_ReturnsOneWithRiichi`, `TestTenpaiProb_ReturnsYamitenRatio`, `TestTenpaiProb_ReturnsOneWithoutStats` | 固定済み |
| OC-014 | `TestImmediateScoreDeltaDist`, `TestImmediateScoreDeltaDistFromStats` | 固定済み |
| OC-015 | `TestCandidateShantenUsesThrowableVector`, `TestCandidateShantenReturnsBaseForNone`, `TestCandidateShantenReturnsInfinityWhenTileIsNotThrowable` | 固定済み |
| OC-016 | `TestTrialWinPts`, `TestCandidateTrialWinPts` | 固定済み |
| OC-017 | `TestExhaustiveDrawProb`, `TestExhaustiveDrawProbOnSelfNoWin` | 固定済み |
| OC-018 | `TestExhaustiveDrawTenpaiProbs`, `TestExhaustiveDrawScoreDeltaDist` | 固定済み |
| OC-019 | `TestAverageRank`, `TestWinProbAgainst`, `TestWinProbFromRelativeScore_UsesStatsWhenAvailable`, `TestWinProbFromRelativeScore_FallsBackToStartingDealerOrder` | 固定済み |
| OC-020 | `TestCompareCandidates`, `TestChooseBestCandidate_PrefersBlackTileOnTie`, `TestChooseBestCandidate_DoesNotPreferRiichiOnScoreTie` | 固定済み |
| OC-021 | `TestFormatCandidateTrace`, `TestFormatDecisionTrace_AppendsDecidedKey` | 固定済み |

## 9. テスト不足

- CoffeeScript 実行結果を正とする action golden がない。既存テストは Go 内部責務の固定が中心である。
- danger tree feature の全量一致テストがない。
- CoffeeScript `calculateFan` と Go `service.CalculateFuHan` の直接比較テストがない。
- `possibleActions` 順序/重複に由来する CoffeeScript の選択挙動と、Go `LegalActions` の順序/正規化の差を検出するテストがない。
- Monte Carlo の乱数差を許容したうえで、近接評価局面の action 安定性を見るテストがない。
- 無和了候補の `averageHoraPoints` / points distribution が最終比較に与える影響を固定するテストがない。

## 10. 監査結論

Go develop の AI 判断ロジックは、既存設計で許容された差分を除けば、CoffeeScript 版 `manue_ai.coffee` の主要な意味単位へ対応する実装を持っている。特に dispatch、候補生成、評価値合成、比較優先順位は Go 側の unit test でかなり固定されている。

ただし、CoffeeScript 版を正とする「意味的等価」を完了判定するには、次の追加監査が必要である。

1. danger feature 全量の CoffeeScript/Go 真偽比較。
2. `calculateFan` / `CalculateFuHan` の直接比較。
3. CoffeeScript と Go を同一局面で実行する action characterization。
4. 無和了候補と乱数差が final action に与える影響の確認。

そのため現時点の最終判定は **「主要構造は意味等価に近いが、High/Medium の未確認点が残るため完全等価とは判定しない」** とする。
