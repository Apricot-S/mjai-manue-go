# Original-vs-Port 差分調査レポート

日付: 2026-07-02

対象:

- `test/original_vs_port/compare`
- `internal/domain/ai`
- `internal/domain/game/round`
- `reference/repositories/mjai-manue-original/coffee`

## 前提と制約

リポジトリ内には比較に使った `.mjson` ログが存在しないため、報告された「20% くらい打牌が異なる」結果そのものは再実行できなかった。
本レポートは、compare 実装、Go port の Manue AI、CoffeeScript original の Manue AI を静的に突き合わせた調査結果である。

実行確認:

```sh
GOEXPERIMENT=jsonv2 go test ./test/original_vs_port/compare
GOEXPERIMENT=jsonv2 go test ./internal/domain/ai
GOEXPERIMENT=jsonv2 go test ./...
GOEXPERIMENT=jsonv2 go fix ./...
GOEXPERIMENT=jsonv2 go vet ./...
```

結果: pass。

## 結論

追加調査で、立直評価の扱いに移植ミスを確認した。これは乱数差では説明しにくい「単純な何切るで、効率より手役を過剰に見る」症状と整合する。

1. **Go port は original の `reachMode=default` を `reachMode=never` 相当に移植していた。**
   - original は通常打牌評価でも、門前なら将来立直できる前提で立直 1 翻を加える。
   - ただし `reachMode=default` では、現在テンパイ形だけに絞り込まない。
   - Go port は `scoreAsRiichi` を「立直を加点する」と「テンパイ以外を捨てる」の両方に使い、通常打牌では `false` にしていた。
   - このため、門前・非テンパイ・役なしの素直な受け入れ形が 0 点扱いになり、目に見える手役を持つ候補を過大評価しやすかった。
   - 対応として `scoreAsRiichi` と `pruneToTenpai` を分離した。default は `scoreAsRiichi=true, pruneToTenpai=false`、即立直/立直済みは `true, true`、立直可能時の通常打牌は `false, false` とする。

2. **対子 goal の列挙差分は、残差分の主因ではなさそう。**
   - original は、手元に 1 枚もない牌を 2 枚引いて対子にする goal を列挙しない。
   - Go port はこの goal を列挙するため、仮説として `myHoraProb` を一方向に押し上げる可能性を疑った。
   - しかし compare ではほとんど変化がなく、残った 13% 程度の不一致の主因ではなさそう。
   - 最近接和了形集合を「現在手牌から必要牌を補って到達できる和了形集合」と定義するなら、手元にない牌を 2 枚引く対子も含める Go port 側のほうが自然。ここは original の過小列挙を直したバグ修正扱いとし、Go port の挙動を維持する。

3. **他家和了時の点棒変化分布で、original 側に target 確率の割り当てバグらしい挙動を確認した。**
   - `zzz_before_fix_bug.txt` の mismatch ログを同一 action 同士で集計すると、`myHoraProb` は正負に割れる一方、`expPt` は 813 行中 812 行で Go port が高い。
   - ブロック平均でも 74/74 ブロックで Go port の `expPt` が高く、候補固有の和了確率ではなく、候補共通の将来局結果分布が疑わしい。
   - original の `getHoraFactorsDist(actor)` は、他家 `actor` の和了分布を作るとき、ツモ確率を `target == actor` ではなく `target == @player()` に割り当てている。
   - その結果、他家ツモと「自分がそのその他家に放銃する」分岐の確率が入れ替わり、original は自分の期待点を過小評価しやすい。
   - Go port は `targetID == actorID` をツモとして扱っており、麻雀ルール上はこちらが正しい。これは Go port の移植ミスではなく、original バグ修正由来の差分と見るのが妥当。

4. **Monte Carlo 乱数のライフサイクルが original と Go port で異なる。**
   - original は `getHoraEstimation` のたびに `seedRandom("")` を作り直す。
   - Go port は `start_game` ごとに PCG を 1 つ作り、以後の意思決定で使い回す。
   - このため Go port の評価値は「同じ局面に直接到達した場合」と「東1局から順に到達した場合」で変わり得る。original は少なくとも和了推定単位では毎回同じ乱数列を使う。
   - これは移植ミス候補として高優先度。compare で CoffeeScript 版との一致率を見るなら、Go 側も original 互換の乱数ライフサイクルにする必要がある。

5. **赤5を捨てる候補の将来和了価値で、Go port は original のバグらしい挙動を直している。**
   - original は `calculateFan(goal, tehais, reachMode)` に、候補打牌前の手牌 `tehais` を渡す。
   - そのため赤5を捨てた候補でも、将来和了の赤ドラ評価にその赤5が残る可能性がある。
   - Go port は `afterDiscardHand` から評価用手牌を作り、捨てた赤5を将来価値から除外している。
   - これは Go port 側が正しい可能性が高いが、「original と同一出力」を目標にする compare では差分要因になる。

## 根拠

### 0. 立直評価 `reachMode=default` の移植ミス

original:

- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:118` 付近の `decideDahai` で、立直可能時は `reachMode="now"` と `reachMode="never"` を別々に評価する。
- 立直できない通常局面では `getMetricsInternal(..., if @player().reachDeclared then "now" else "default")` を使う。
- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:652` 付近の `calculateFan(goal, tehais, reachMode)` は `reachMode != "never"` なら立直役を加える。
- 一方、テンパイ形への絞り込みは `reachMode == "now"` の経路でだけ行われる。

Go port 修正前:

- `internal/domain/ai/self_turn_candidates.go` の通常打牌候補では、立直済みでない限り `scoreAsRiichi=false` だった。
- `internal/domain/ai/win_goals.go` では `scoreAsRiichi && goal.Shanten > 0` で非テンパイ goal を除外していた。
- つまり Go port の `scoreAsRiichi` は、original の「立直 1 翻を評価に足す」と「今すぐ立直するのでテンパイ形だけ見る」を混同していた。

影響:

- original の `default` は「将来、門前でテンパイしたら立直できる」という期待値を含めるヒューリスティックである。
- Go port 修正前はこの 1 翻が消えるため、役なしだが受け入れの広い門前手の将来和了価値が 0 になり得る。
- その結果、受け入れ枚数やシャンテン改善よりも、タンヤオ/役牌/染めなどの明示的な手役が過剰に強く見える。
- ユーザー報告の「乱数差では説明できない」「単純な期待値通りの何切るでも過剰に手役を見る」症状と一致する。

対応:

- `internal/domain/ai/candidate.go` に `pruneToTenpai` を追加し、立直加点とテンパイ絞り込みを分離した。
- `internal/domain/ai/self_turn_candidates.go` で、通常局面の打牌候補を `scoreAsRiichi=true, pruneToTenpai=false` にした。
- 即立直/立直済みは `scoreAsRiichi=true, pruneToTenpai=true`、立直可能時の通常打牌は `scoreAsRiichi=false, pruneToTenpai=false` とした。
- `internal/domain/ai/call_candidates.go` では、副露を pass する `none` 候補だけ default 評価として `scoreAsRiichi=true` にした。original は実際に chi/pon/daiminkan した後の候補にも `reachMode="default"` を渡すが、`addYaku` が `goal.furos.length > 0` の場合に立直役を 0 翻にするため、実効的には立直加点されない。Go port ではその実効挙動を `scoreAsRiichi=false` として直接表現する。
- original は `goal.furos.length > 0` の場合に立直役の食い下がり翻 `kuiFan=0` で抑止する。Go port も `len(context.melds) == 0` のときだけ立直加点することで、実効挙動を合わせた。

original 側のバグ候補:

- original の `default` で将来立直を評価に含めること自体は、AI 評価のヒューリスティックとして一貫しており、今回の症状の原因ではなく Go port 側の移植ミスと判断する。
- ただし original は副露数で立直役を抑止しているため、暗槓だけを持つ門前扱い可能な形でも立直加点が消える可能性がある。これは麻雀ルールとしては疑わしいが、互換性のため Go port も同じ挙動に寄せている。

### 1. 最近接和了形 goal の対子列挙差分

original:

- `reference/repositories/mjai-manue-original/coffee/shanten_analysis.coffee` は刻子/順子について、現在手牌と 1 枚も重ならないブロックを `current1 < current + 3` で除外する。
- 同じく対子についても、`delta = max((vector1[i] + 2) - vector0[i], 0)` が `2` 以上なら候補から除外する。
- つまり original は、完全に手元にない牌を 2 枚引いて対子にする goal を列挙しない。

Go port:

- `internal/domain/game/round/service/shanten.go` は刻子/順子については、完全に手元にないブロックを除外していた。
- 対子については、現在手牌に 0 枚の牌でも `pairDistance == 2` の goal を許している。

検証結果:

- この差分は `myHoraProb` を一方向に高くする可能性があると考えたが、compare の結果はほとんど変化しなかった。
- そのため、残った不一致の主因ではないと判断する。

方針:

- 最近接和了形集合の定義としては、現在手牌から必要牌を補って到達できる和了形を含めるのが自然であり、手元にない牌を 2 枚引く対子も除外しないほうがよい。
- Go port の挙動を維持し、original との差分は original の過小列挙を補ったバグ修正扱いとする。

### 2. `zzz_before_fix_bug.txt` の評価値差分集計

`zzz_before_fix_bug.txt` の追加観察:

- 同一 action が original / Go port の両方に出ている評価行 813 件を集計した。
- `myHoraProb` 差分は平均 `+0.0015`、正 318 件 / 負 349 件 / 同値 146 件で、常に Go port が高いわけではない。
- 一方、`expPt` 差分は平均 `+521`、中央値 `+563`、正 812 件 / 負 1 件だった。
- mismatch ブロック単位でも、共通候補を持つ 74 ブロックすべてで Go port の平均 `expPt` が高かった。

このため、残差分の主因候補は `myHoraProb` 単体ではなく、将来局結果の点棒変化分布、特に候補共通で混ざる「他家和了」の扱いである。

### 2.1 他家和了時の target 確率割り当て差分

original:

- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:432` 付近の `getRandomHoraScoreChangesDist(actor)` が、他家 `actor` のランダム和了時の点棒変化分布を作る。
- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:444` 付近の `getHoraFactorsDist(actor)` で、各 `target` に確率を割り当てる。
- その確率は `prob = (if target == @player() then tsumoHoraProb else (1 - tsumoHoraProb) / 3)` になっている。
- これは自分の和了評価 `getScoreChangesDistOnHora` では正しいが、他家 `actor` の和了評価では不自然。ツモ確率を割り当てるべき target は `@player()` ではなく `actor` である。

影響:

- `actor != @player()` の場合、original は「他家のツモ」と「自分がそのその他家へ放銃」の確率を入れ替える。
- 自分が支払う額は、他家ツモなら通常 `1/4` または `1/2`、自分の放銃なら `1` である。
- そのため original は他家和了時の自分の失点を過大に見積もり、`expPt` と `avgRank` を悪く出しやすい。
- `zzz_before_fix_bug.txt` で Go port の `expPt` がほぼ常に高いこと、しかし `myHoraProb` は一方向に偏っていないことと整合する。

Go port:

- `internal/domain/ai/win_score.go:51` の `winScoreFactorDist(actorID, dealerID, selfDrawProb)` は、`targetID == actorID` をツモとして扱う。
- `targetID != actorID` はロン対象として一様に割り当てる。
- 麻雀ルール上はこちらが自然であり、Go port 側を original に合わせて壊すべきではない。

方針:

- この差分は original 互換性だけを見る compare では mismatch の原因になる。
- ただし AI の評価としては Go port 側が正しいため、バグ修正扱いの既知差分として分類するのがよい。
- original 完全互換モードを別途作る場合だけ、`getHoraFactorsDist` の target 割り当てバグを再現する互換オプションを検討する。

### 3. 乱数の差分

original:

- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:496` で `getHoraEstimation` に入る。
- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:532` で `random = seedRandom("")`。
- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:538` で `Util.shuffle(invisiblePids, random, numTsumos)`。

Go port:

- `internal/domain/ai/manue_agent.go:38` の `Reset` は `start_game` 単位。
- `internal/domain/ai/manue_agent.go:44` で `rand.New(rand.NewPCG(a.seed, 0))`。
- `internal/domain/ai/evaluator.go:94` で同じ evaluator の RNG を `winEstimatesFromState` に渡す。
- `internal/domain/ai/win_estimator.go:47` で試行ごとにシャッフルする。

影響:

- original は self turn の reach/never 比較、鳴き/none 比較で、各 `getMetricsInternal` が同じ seed から始まる。
- Go port は候補群をまとめて評価するため候補間の試行共有はあるが、ゲーム内の過去意思決定で RNG が進む。
- 1000 試行の Monte Carlo 評価なので、近い候補では乱数列差だけで `averageRank` / `expectedPoints` の順位が入れ替わる。

### 4. 赤5評価の差分

original:

- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:174` で `getHoraEstimation(candDahais, analysis, tehais, furos, reachMode)`。
- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:508` で `calculateFan(goal, tehais, reachMode)`。
- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:630` 以降の `calculateFan` は `tehais.concat(furoPais)` から赤ドラ数を数える。

Go port:

- `internal/domain/ai/self_turn_candidates.go:110` で候補に `afterDiscardHand` を持たせる。
- `internal/domain/ai/win_goals.go:74` に original の赤5評価差分に関するコメントがある。
- `internal/domain/ai/win_goals.go:75` で after-discard hand を使う方針。
- `internal/domain/ai/win_goals.go:76` で `scoringHandForGoal(candidate.afterDiscardHand, goal.Blocks)`。

影響:

- 赤5を打つ候補を original が過大評価し、Go port が低く評価する局面がある。
- 逆に、赤5を残す候補と赤5を打つ候補の順位が original と Go port で入れ替わる可能性がある。
- これは shallow copy というより、original の評価入力が候補打牌前の手牌のままになっている設計/実装バグ候補。

### 5. compare 側の水増し可能性

compare は action の意味を正規化して比較しているため、`log` や `consumed` の順序差は主因ではなさそう。

- `test/original_vs_port/compare/action.go:19` で正規化。
- `test/original_vs_port/compare/action.go:109` で比較。
- `test/original_vs_port/compare/action.go:113` で `pai` と `tsumogiri` も比較対象。

注意点:

- 同じ `pai` でも `tsumogiri` だけ違えば mismatch になる。
- ユーザーの「打牌が異なる」が、牌そのものの差なのか `tsumogiri` 差なのかは、現在の summary だけでは分からない。
- 20% の内訳確認には、mismatch を `type差` / `pai差` / `tsumogiriのみ差` / `red 5絡み` / `鳴き差` に分類する後処理が必要。

### 6. original の shallow copy / 破壊的変更候補

確認した範囲では、主要な打牌判断経路で shallow copy が即座に大量差分を生む証拠は薄い。

- `decideFuro` の `tehais = @player().tehais.concat([])` は浅いコピーだが、要素の `Pai` は実質 immutable で、配列から削除するだけ。
- `goal` へ `requiredBitVectors` / `furos` / `yakus` などを追加しているが、`getMetricsInternal` ごとに新しい `ShantenAnalysis` を作るため、呼び出し間での汚染は起きにくい。
- `Game.rankedPlayers` は `@_players.sort` で内部配列を破壊的に並べ替えるが、Manue AI 本体の意思決定経路では `rankedPlayers()` を使っていない。

ただし、original には別の数値的な危うさがある。

- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee:590` で `totalPointsVector[pid] / totalHoraVector[pid]` を計算する。
- 和了試行が 0 回の候補では `NaN` になり得る。
- Go port は `totalWins == 0` の場合に平均点・期待点を 0 として扱う。
- この差は通常はログ表示や確率 0 の分布に閉じるはずだが、低確率候補が多い局面では念のため確認対象。

## 優先して切り分けるべきこと

1. compare 出力の mismatch を分類する。
   - 同じ `dahai.pai` で `tsumogiri` だけ違う件数。
   - `5m` / `5mr` / `5p` / `5pr` / `5s` / `5sr` が関係する件数。
   - `reach` 直後の打牌差。
   - `chi` / `pon` 後の打牌差。

2. 他家和了 target 確率差分を compare で分類できるようにする。
   - Go port の正しい分布を維持する場合、original バグ修正由来の既知差分として扱う。
   - 完全互換検証をしたい場合だけ、一時的な互換モードで `targetID == selfID` をツモ扱いにする A/B 実験を行う。

3. Go port の実験ブランチで RNG を「評価ごとに reset」して compare する。
   - まず PCG のまま `evaluateCandidates` または `winEstimatesFromState` 単位で reset。
   - これで mismatch が大きく減るなら、主因は RNG ライフサイクル。
   - まだ多い場合は `seed-random` 互換 RNG か、original のシャッフル順そのものを移植する必要がある。

4. 赤5評価を original 互換に戻した実験版で compare する。
   - `afterDiscardHand` ではなく候補打牌前の手牌で赤ドラ評価する互換モードを一時的に試す。
   - 赤5絡み mismatch が減るなら、差分は original バグ修正由来。

5. 方針を決める。
   - 「original 完全互換」を重視するなら、他家和了 target 確率、RNG、赤5評価バグも互換に寄せる。
   - 「Go port として妥当な AI」を重視するなら、他家和了 target 確率と赤5評価の修正は維持し、compare では許容差分として分類・除外できるようにする。

## 推奨する次アクション

- compare に `--classify-mismatches` 相当の集計を追加する。
- 他家和了 target 確率の互換モードを A/B 実験し、`expPt` の一方向オフセットが消えるか確認する。
- 乱数 reset の A/B 実験を小さい差分で行う。
- 赤5評価の互換モードを実験し、赤5絡み mismatch の減少量を見る。
- その結果をもとに `.agents/design.md` の「同一入力 -> 同一出力」の扱いを、bug-compatible にするのか、known divergence を許容するのか明記する。
