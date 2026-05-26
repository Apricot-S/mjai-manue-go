# Manue AI Original Specification

この文書は CoffeeScript 版 `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee` の動作仕様である。Go への移植設計は含めず、オリジナル実装がどの入力からどの評価値を作り、どの順序で最終 action を選ぶかを記録する。

## 1. 全体の呼び出し構造

`ManueAI` は `AI` を継承し、`respondToAction(action)` を意思決定入口にする。局面状態、プレイヤー、可能行動は親クラスと mjai 入力由来の `action.possibleActions` から参照する。

```text
respondToAction
  ├─ categorizeActions
  ├─ hora があれば即 createAction(hora)
  ├─ 自分の tsumo/chi/pon/reach
  │   ├─ reach accepted 中の tsumo はツモ切り
  │   └─ decideDahai
  │       ├─ getMetrics
  │       │   ├─ getMetricsInternal(reachMode="now")
  │       │   ├─ selectTenpaiMetrics
  │       │   ├─ getMetricsInternal(reachMode="never"/"default")
  │       │   └─ mergeMetrics
  │       ├─ printMetrics
  │       ├─ printTenpaiProbs
  │       └─ chooseBestMetric
  └─ 他家の dahai/kakan
      ├─ hora があれば即 createAction(hora)
      └─ decideFuro
          ├─ none 候補の getMetricsInternal
          ├─ 各副露候補の手牌/副露を仮適用
          ├─ isKuikae で直後打牌候補を除外
          ├─ getMetricsInternal
          ├─ printMetrics
          ├─ printTenpaiProbs
          └─ chooseBestMetric
```

`respondToAction` が扱う自分のイベントは `tsumo`、`chi`、`pon`、`reach` である。他家イベントでは `dahai` と `kakan` にだけ反応する。それ以外、または選ぶべき行動がない場合は `type: "none"` を返す。

## 2. 行動候補

### 2.1 共通分類

`categorizeActions(actions)` は mjai の `possibleActions` を次の 3 種類へ分ける。

- `hora`: 最後に見つかった `type == "hora"`。
- `reach`: 最後に見つかった `type == "reach"`。
- `furos`: `hora` / `reach` 以外の action 全部。

和了は最優先で、打牌評価や副露評価を行わず即選択する。

### 2.2 打牌候補

`getMetrics(forbiddenDahais, reachDeclared, canReach)` は現在の手牌 `@player().tehais` から `candDahais` を作る。

- 既に候補に同じ牌がある場合は追加しない。比較は `Pai.equal`。
- `forbiddenDahais` に同じ牌がある場合は追加しない。
- 赤牌は `Pai.equal` の結果に従う。赤を同一視する追加処理はない。

`canReach` がある場合は、同じ打牌集合に対して 2 系統を作る。

- prefix `0`: `reachMode="now"` で評価し、`shanten <= 0` の候補だけ残す。選ばれた場合は reach action を返す。
- prefix `-1`: `reachMode="never"` で評価する。選ばれた場合は通常打牌を返す。

`canReach` がない場合は `reachDeclared` に応じて `reachMode="now"` または `"default"` で評価する。`reachDeclared` の場合は `shanten <= 0` の候補だけ残す。

`decideDahai` は最善 key を `"<actionIdx>.<pai>"` として受け取り、`actionIdx == 0` なら reach を返す。通常打牌を返す場合、元イベントが `tsumo` または `reach` で、選択牌が手牌末尾の牌と `Pai.equal` なら `tsumogiri: true` にする。

### 2.3 副露候補

`decideFuro(furoActions)` は見送り候補 `none` と各副露候補を同じ比較器で評価する。

- `none`: 現在の手牌/副露、候補 `[null]`、`reachMode="default"` で評価する。
- 各副露候補: `action.consumed` に一致する牌を手牌から 1 枚ずつ削除し、`Furo(type, taken, consumed, target)` を副露へ追加する。
- 副露後の打牌候補は副露後手牌の重複除去済み牌で、`isKuikae(action, pai)` が true の牌は除外する。
- 副露候補の key は `"<furoActionIndex>.<dahai>"`。最終的に key が `none` なら `type: "none"`、それ以外なら対応する副露 action を返す。鳴いた直後の打牌は返さない。

`isKuikae` は `furoAction.consumed + dahai` の 3 枚をソートし、3 枚同種牌、または連続した同色数牌なら喰い替えとする。

## 3. 候補評価

`getMetricsInternal(tehais, furos, candDahais, reachMode)` が候補ごとの中心評価である。入力手牌から `ShantenAnalysis(..., {allowedExtraPais: 1})` を作り、候補ごとに以下を合成する。

```text
getMetricsInternal
  ├─ analysis = ShantenAnalysis(tehais)
  ├─ safeProbs = getSafeProbs(candDahais, analysis)
  ├─ immediateScoreChangesDists = getImmediateScoreChangesDists(candDahais)
  ├─ metrics = getHoraEstimation(candDahais, analysis, tehais, furos, reachMode)
  ├─ 流局時平均点/分布
  ├─ 他家和了時 score changes 分布
  └─ 候補ごとに最終 scoreChangesDist / expectedPoints / averageRank を計算
```

候補 key は打牌なら `pai.toString()`、見送りなら `none`。各 metric は少なくとも次を持つ。

- `red`: 候補牌が赤牌か。
- `safeProb`: 直後に全相手へ放銃しない確率。
- `hojuProb`: `1 - safeProb`。
- `horaProb`: 自分が和了する確率。
- `averageHoraPoints`: 自分が和了した条件付き平均点。
- `horaPointsDist`: 自分の和了点分布。
- `expectedHoraPoints`: Monte Carlo 全試行あたりの自分の和了点期待値。
- `ryukyokuProb`: 最終的には `(1 - horaProb) * getRyukyokuProbOnMyNoHora()`。
- `othersHoraProb`: `(1 - horaProb) * (1 - getRyukyokuProbOnMyNoHora())`。
- `ryukyokuAveragePoints`: 現在候補の `shanten <= 0` なら聴牌扱い、そうでなければノーテン扱いの流局平均点。
- `immediateScoreChangesDist`: 直後放銃または無変化の点差分布。
- `scoreChangesDistOnRyukyoku`: 流局時の点差分布。
- `scoreChangesDistOnHora`: 自分和了時の点差分布。
- `futureScoreChangesDist`: 自分和了、流局、他家和了を混合した将来分布。
- `scoreChangesDist`: `immediateScoreChangesDist` の無変化 branch を `futureScoreChangesDist` で置換した最終分布。
- `expectedPoints`: `scoreChangesDist.expected()[@player().id]`。
- `averageRank`: `getAverageRank(scoreChangesDist)`。
- `shanten`: `getHoraEstimation` が求めた候補別向聴。

`safeExpectedPoints`、`unsafeExpectedPoints`、`ryukyokuExpectedPoints` も計算されるが、最終比較には使われない。

## 4. 危険度と直後放銃

`getSafeProbs(candDahais, analysis)` は候補ごとに初期値 1 を置き、各相手に対して安全確率を掛け合わせる。

- 候補が `null` の場合、`none` は常に 1。
- `scene.anpai(pai)` が true なら、その相手への `safeProb` は 1。
- それ以外は `DangerEstimator.estimateProb(scene, pai).prob` を放銃条件付き確率とし、相手の `getTenpaiProb(player)` を掛けて `1 - tenpaiProb * prob` を安全確率とする。
- feature はログ用に `"name value"` 形式へ変換される。

`getImmediateScoreChangesDists(candDahais)` は直後放銃による score change を候補ごとに分布化する。

- 初期分布は無変化 `[0,0,0,0]` が確率 1。
- 各相手について、親/子別の `*_HoraPointsFreqs` から和了点分布を作る。
- 放銃時の単位 score change は和了者 `+1`、自分 `-1`、他者 `0`。これに和了点分布を掛ける。
- 候補が安全牌なら放銃確率 0。そうでなければ `tenpaiProb * dangerProb`。
- ダブロン/トリロンの組み合わせ爆発を避けるため、無変化 branch を相手ごとに置換し、最初のロンだけを考慮する。

`getTenpaiProb(player)` は、リーチ状態が `none` でなければ 1。そうでなければ残り巡目と副露数から `yamitenStats["numRemainTurns,numFuros"]` を引き、存在すれば `tenpai / total`、なければ 1 を返す。

## 5. 和了推定

`getHoraEstimation(candDahais, analysis, tehais, furos, reachMode)` は Monte Carlo で候補別の自分和了確率と和了点を推定する。

1. `analysis.goals()` を走査する。
2. `reachMode == "now"` かつ `goal.shanten > 0` の goal は除外する。
3. `analysis.shanten() > 3` かつ `goal.shanten > analysis.shanten()` の goal は速度のため除外する。
4. `goal.requiredVector` を bit vector 化し、候補副露を goal に持たせ、`calculateFan` で fan/fu/points/yakus を計算する。
5. `goal.points > 0` の goal だけ採用する。
6. 自分から見えている牌を `game.visiblePais(@player())` で集め、全牌集合から除いて未見牌 ID 列を作る。
7. `getNumExpectedRemainingTurns()` でツモ試行枚数を決め、1000 回固定で試行する。
8. 乱数は `seedRandom("")` 固定。呼び出しごとに同じ乱数列を使い、リーチあり/なしなどの比較条件を揃える。
9. 各試行で未見牌を shuffle し、先頭 `numTsumos` 枚をツモ集合として bit vector 化する。
10. goal の required bit vectors がツモ bit vectors の subset なら達成扱いにする。
11. 達成 goal について、`pid == Pai.NUM_IDS` または `goal.throwableVector[pid] > 0` の候補を和了可能にする。
12. 同一候補で複数 goal が達成した場合は points が最大の goal を採用する。

候補別 `shantenVector` は打牌後手牌を再解析せず、`analysis.goals()` の `throwableVector` から、その牌を捨てられる goal の最小 `goal.shanten` を採用する。`Pai.NUM_IDS`、つまり `none` は `analysis.shanten()`。

返却 metric は候補ごとに次を持つ。

- `horaProb = totalHoraVector[pid] / 1000`
- `averageHoraPoints = totalPointsVector[pid] / totalHoraVector[pid]`
- `horaPointsDist = points frequency / totalHoraVector[pid]`
- `expectedHoraPoints = totalPointsVector[pid] / 1000`
- `shanten = shantenVector[pid]`

`totalHoraVector[pid] == 0` の場合の除算は特別扱いされていない。

## 6. 役・点数推定

`calculateFan(goal, tehais, reachMode)` は goal の面子と副露を合わせた `mentsus` から簡易的に役と点数を計算する。副露種別は `FURO_TYPE_TO_MENTSU_TYPE` で `chi -> shuntsu`、`pon -> kotsu`、各 kan -> `kantsu` に変換する。

対象役は以下である。

- `reach`: `reachMode != "never"` のとき 1 飜、鳴き時 0 飜。
- `tyc`: 断么九。
- `cty`: 混全帯么九。門前 2、鳴き 1。
- `pf`: 平和。全 shuntsu、雀頭が役牌でない場合。両面条件は TODO。
- `ykh`: 役牌刻子/槓子の合計飜。
- `ipk`: 一盃口。門前 1、鳴き 0。
- `ssj`: 三色同順。門前 2、鳴き 1。
- `ikt`: 一気通貫。門前 2、鳴き 1。
- `tth`: 対々和。
- `cis`: 清一色。門前 6、鳴き 5。
- `his`: 混一色。門前 3、鳴き 2。
- `dr`: goal 内の牌と同種のドラ表示由来ドラ数。1 飜以上の役がある場合だけ加算。
- `adr`: 手牌/副露内の赤牌のうち、goal 構成牌に同種牌がある数。1 飜以上の役がある場合だけ加算。

fu は平和または副露ありなら 30、それ以外は 40 とする。`getPoints(fu, fan, oya)` は満貫以上を base points で丸め、最終的に親なら `base * 6`、子なら `base * 4` を 100 点切り上げする。これはロン点相当の単一値として扱われる。

## 7. 流局・他家和了・順位期待値

`getRyukyokuProb()` は現在巡目以降の `numTurnsDistribution` の残り質量を分母にし、`ryukyokuRatio / den` を返す。`getRyukyokuProbOnMyNoHora()` はその `3/4` 乗。

`getNotenRyukyokuTenpaiProb()` は、現在ターン `game.turn() + 1/4` から最終ターンまでの `ryukyokuTenpaiStat.tenpaiTurnDistribution` と `noten` を使い、現在ノーテンのプレイヤーが流局時聴牌になる確率を推定する。

`getRyukyokuAveragePoints(selfTenpai)` と `getScoreChangesDistOnRyukyoku(selfTenpai)` は、各プレイヤーの現在聴牌確率とノーテンから聴牌になる確率を掛け合わせ、流局時の 3000 点授受を期待値または分布として計算する。全員聴牌/全員ノーテンは 0。

`getRandomHoraScoreChangesDist(actor)` は、指定 actor が和了した場合の点数分布と `getHoraFactorsDist(actor)` を掛ける。他家和了は `getMetricsInternal` で相手 3 人へ `othersHoraProb / 3` ずつ割り当てられる。

`getAverageRank(scoreChangesDist)` は、自分と各相手の次局以降の勝率を `getWinProb` で計算し、それぞれを独立 Bernoulli として `winsDist` に加算する。自分の平均順位は `4 - 自分を上回る相手数` の期待値。

`getWinProb` は次局、起家からの自分/相手位置、点差から `winProbsMap` を引く。統計が欠けている場合やオーラスの場合は、起家から近い方が同点有利になる規則で 0/1 を返す。

## 8. 最終選択とログ

`chooseBestMetric(metrics, preferBlack)` は全 metric を走査し、`compareMetric` が小さい候補を最善とする。比較順は固定である。

1. `averageRank` が小さい候補。
2. `expectedPoints` が大きい候補。
3. `preferBlack` が true の場合だけ、赤牌を切らない候補。
4. それ以外は同等。走査順で先に現れた候補が残る。

通常打牌では `preferBlack=true`、副露判断では `preferBlack=false`。

`printMetrics` は候補を `compareMetric(..., true)` で並べ、以下の列を `@log` へ表形式で出す。

- `action`: key。
- `avgRank`: `averageRank`。
- `expPt`: `expectedPoints`。
- `hojuProb`: `hojuProb`。
- `myHoraProb`: `horaProb`。
- `ryukyokuProb`: `ryukyokuProb`。
- `otherHoraProb`: `othersHoraProb`。
- `avgHoraPt`: `averageHoraPoints`。
- `ryukyokuAvgPt`: `ryukyokuAveragePoints`。
- `shanten`: `shanten`。

`decideDahai` と `decideFuro` は `console.log("decidedKey", key)` を出す。`printTenpaiProbs` は相手ごとの聴牌確率を `tenpaiProbs:` として `@log` へ出す。

## 9. Characterization ケース候補

Go 移植前に、次の単位で CoffeeScript 版の出力を固定できる局面を用意する。

- 和了可能時は評価せず `hora` を返す。
- リーチ accepted 中の自摸はツモ切りを返す。
- `canReach` ありでは `reachMode="now"` の聴牌候補と `reachMode="never"` の通常候補を比較する。
- `reachDeclared` 中は聴牌候補だけが残る。
- 通常打牌の同一牌重複除去と `cannotDahai` 除外。
- 副露判断で `none` と副露候補を比較し、選択 action は副露のみを返す。
- 喰い替え判定で副露後打牌候補が除外される。
- 安全牌は直後放銃確率 0 になる。
- リーチ者の聴牌確率は 1 になる。
- 候補比較は平均順位、期待点、赤牌回避の順で tie-break する。
