# Performance Improvement Candidates

この文書は、実装を効率化できる候補をまとめた調査メモである。
設計の一次資料は引き続き `docs/design.md` とし、この文書は性能改善タスクの棚卸しと優先順位付けに使う。

## 調査時点

- 調査日: 2026-05-06
- 主な対象:
  - `internal/domain/game/round/service`
  - `internal/domain/game/round`
  - `internal/domain/game/round/player/hand`
  - `internal/domain/game/tile`

## ベンチマーク

既存の向聴計算ベンチマークを実行した。

```powershell
go test ./internal/domain/game/round/service -bench BenchmarkShantenAnalysis -benchmem
```

結果の抜粋:

| Benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| `BenchmarkShantenAnalysis_Normal` | 408598 | 1465427 | 1144 |
| `BenchmarkShantenAnalysis_HalfFlush` | 76236 | 355088 | 279 |
| `BenchmarkShantenAnalysis_FullFlush` | 35877 | 182439 | 145 |
| `BenchmarkShantenAnalysis_NonSimple` | 113710 | 423050 | 252 |
| `BenchmarkShantenAnalysis_14_Normal` | 1070379 | 4156273 | 3735 |
| `BenchmarkShantenAnalysis_14_HalfFlush` | 193739 | 691248 | 761 |
| `BenchmarkShantenAnalysis_14_FullFlush` | 82433 | 339399 | 351 |
| `BenchmarkShantenAnalysis_14_NonSimple` | 220120 | 726994 | 605 |

通常14枚の向聴解析が特に重く、allocation も多い。
最初に見るべき箇所は向聴計算の DFS と、その呼び出し側である。

## 優先度 高

### 1. `AnalyzeShanten` の DFS 中の allocation 削減

Status: 改善済み（一部）。DFS 中の block slice copy をやめ、Goal 保存時だけ block 列をコピーするように変更済み。

対象:

- `internal/domain/game/round/service/shanten.go`
- `AnalyzeShanten`
- `analyzeShantenInternal`
- `makeNewBlocks`

現状:

- DFS の各枝で `makeNewBlocks` が `blocks` slice をコピーしている。
- Goal を返すために探索途中の block 列を保持しているが、向聴数だけが必要な呼び出しでも同じ経路を通る。
- 14枚通常手で `3735 allocs/op` が出ており、ここが主要因の可能性が高い。

改善案:

- 再帰中は固定長の `[5]block.Block` または再利用 slice を使い、append/backtrack で探索する。
- `Goal` を確定して保存する時だけ block 列をコピーする。
- `AnalyzeShanten` を「向聴数だけ欲しい API」と「Goal も欲しい API」に分ける。
  - 例: `ShantenNumber(h, opts...) int`
  - 例: `AnalyzeShanten(h, opts...) (int, []Goal)`
- `IsTenpaiGeneral` / `IsTenpaiAll` / リーチ可否判定のような用途では Goal を生成しない経路を使う。

期待効果:

- allocation と GC 圧力の削減。
- `IsTenpaiAll` を多数回呼ぶ合法手列挙や将来の AI 評価で効きやすい。

注意点:

- `Goal.Blocks` の順序がテストや移植元互換に影響している可能性があるため、既存テストに加えて Goal 内容を確認するテストを維持する。
- DFS の枝刈り条件は挙動互換性に直結するため、最初の変更は allocation 削減に限定する。

### 2. 和了形判定の lookup table / memo 化

対象:

- `internal/domain/game/round/service/win.go`
- `isWinningFormGeneral`
- `isSingleColorWinningFormWithoutPair`
- `isSingleColorWinningFormWithPair`
- `isHonorsWinningForm`

現状:

- `IsWinningForm` は呼び出しごとに `TileCounts34` から各色と字牌を分解して判定している。
- `waitsFor` では最大34回、`canConcealedKanAfterRiichi` では2つの手牌に対して呼ばれる。

改善案:

- 1色9枚の count vector を base-5 などで整数 key に encode し、`withPair` / `withoutPair` の判定結果を table 化する。
- 牌数は各 tile count が 0..4 なので、1色あたり最大 `5^9 = 1,953,125` 通りである。
- 字牌7枚も base-5 key で table 化できる。
- 起動時生成ではなく package init 前の静的生成が望ましいが、まずは `sync.Once` による lazy table でもよい。

期待効果:

- `IsWinningForm` の定数時間化。
- 待ち計算、和了可否、1翻有無判定の下支えになる。

注意点:

- 通常形のみの table と、七対子・国士無双は分ける。
- 赤牌は `TileCounts34` 化された後の判定なので、table key には含めない。

### 3. `waitsFor` の bitset 化と手牌生成削減

Status: 改善中。待ち集合の `uint64` bitset 化は実装済み。`VisibleHand` 生成削減は未実施。

対象:

- `internal/domain/game/round/legal_actions_waits.go`
- `canConcealedKanAfterRiichi`
- `waitsFor`

現状:

- `waitsFor` は `map[int]struct{}` を作る。
- 34種類の牌について `VisibleHand.Draw` を試し、`service.IsWinningForm` を呼ぶ。
- `canConcealedKanAfterRiichi` は `waitsFor(handBeforeKan)` と `waitsFor(handAfterKan)` を `maps.Equal` で比較する。
- `handAfterKan` 生成時に `TileCounts34.ToTiles()` から `VisibleHand` を作り直している。

改善案:

- 待ち集合を `uint64` の bitset にする。
  - 34種類なので `1 << id` で表現できる。
  - 比較は整数比較だけで済む。
- `VisibleHand` を毎回生成せず、`TileCounts34` を直接 +1/-1 して `isWinningForm...` を呼ぶ helper を追加する。
- 必要に応じて `TileCounts34` を key にした memo を導入する。

期待効果:

- map allocation の削減。
- 暗槓後待ち不変チェックと将来の待ち列挙で効く。

注意点:

- `VisibleHand.Draw` が担っている「5枚目を引けない」検証を、直接 count を触る helper 側でも維持する。
- 赤牌を含む手牌でも待ち判定は 34 種類で行うため、`RemoveRed` 済みの count を使う。

## 優先度 中

### 4. 牌コード変換と么九牌判定の table 化

Status: 改善済み。`NewTileFromCode` / `MustTileFromCode` は code -> ID map、`Tile.IsYaochu` は bool lookup table を使うように変更済み。

対象:

- `internal/domain/game/tile/tile.go`
- `NewTileFromCode`
- `MustTileFromCode`
- `Tile.IsYaochu`

現状:

- `NewTileFromCode` / `MustTileFromCode` は `slices.Index(tileCodes[:], code)` で線形探索する。
- `Tile.IsYaochu` は `slices.Contains(YaochuhaiIDs[:], t.ID())` で線形探索する。

改善案:

- `map[string]int` または `switch` で tile code から ID へ変換する。
- `IsYaochu` は `[tile.NumTileType38]bool` の lookup table にする。

期待効果:

- protocol decode や役判定で頻出する小さな線形探索を削減できる。
- 実装リスクが低く、テーブル駆動テストで確認しやすい。

注意点:

- 不正 code の error message は既存テストに合わせる。
- unknown tile `?` の扱いを変えない。

### 5. `MustTileFromID` の値生成削減

対象:

- `internal/domain/game/tile/tile.go`
- `newTileFromValidID`
- `NewTileFromID`
- `MustTileFromID`
- `Tile.AddRed`
- `Tile.RemoveRed`
- `Tile.Next`
- `Tile.NextForDora`

現状:

- `MustTileFromID` は呼び出しごとに `newTileFromValidID` で `Tile` を作る。
- `AddRed` / `RemoveRed` / `NextForDora` などの小さな変換でも `MustTileFromID` を経由する。

改善案:

- package-level の `[tile.NumTileType38]Tile` を用意し、ID から既存値を返す。
- pointer を返す API は既存のままにする場合、配列要素の address を返す。
- 値返し API を追加できるなら、内部 hot path は値返しを使う。

期待効果:

- 小さな allocation / escape の削減。
- table-driven な牌変換と相性がよい。

注意点:

- 呼び出し側が返された `*Tile` を変更しない前提で安全にする。
  - `Tile` の field は unexported なので現状は外部 package から変更できない。
- pointer identity に依存していないことを確認する。

### 6. 役判定の nested map を配列または bit mask に置換

対象:

- `internal/domain/game/round/service/yaku.go`
- `sanshokuDoujun`
- `ikkiTsuukan`
- `sanshokuDoukou`

現状:

- 呼び出しごとに `map[rune]map[int]bool` を作る。
- `Has1Han` では Goal ごとに複数の役判定が呼ばれる。

改善案:

- 色を 0..2 に encode して `[3][10]bool` を使う。
- さらに軽くするなら、数字ごとの bit mask にする。
  - 三色同順: `sequenceByNumber[n]` に色 bit を立て、`0b111` を確認する。
  - 一気通貫: `sequenceByColor[color]` に開始数字 bit を立て、1/4/7 を確認する。
  - 三色同刻: triplet/quad の数字別色 bit を見る。

期待効果:

- map allocation の削減。
- `Has1Han` の Goal 数が多い局面で効きやすい。

注意点:

- `isOpen` による翻数差だけは既存ロジック通りに残す。
- 字牌を三色同刻の map に入れている現状は実質使われていないが、挙動確認のため変更時はテストを追加する。

## 優先度 低から中

### 7. 合法打牌列挙の `Distinct` 回避

対象:

- `internal/domain/game/round/legal_actions_self_draw.go`
- `legalActionsOnSelfDraw`
- `legalPromotedKanActions`
- `tile.Tiles.Distinct`

現状:

- 手牌を slice 化し、sort/compact して distinct tile を列挙している。
- 手牌は内部的に count 配列を持っているため、distinct 列挙だけなら sort は不要である。

改善案:

- `VisibleHand` に distinct tile iterator/helper を追加する。
  - 例: `DistinctTiles37() []tile.Tile`
  - 例: `ForEachDistinct(func(tile.Tile) bool)`
- 合法打牌列挙では count が 1 以上の ID を昇順に走査する。

期待効果:

- 合法手列挙の allocation と sort cost を削減できる。

注意点:

- 赤5と通常5を別 tile として出すべき局面と、同一 symbol として扱う局面を分ける。
- `SwapCallTiles` 判定は `HasSameSymbol` を使っているため、赤牌を含む場合の挙動を維持する。

### 8. リーチ可否判定の向聴計算共有

対象:

- `internal/domain/game/round/legal_actions_self_draw.go`
- `canRiichi`
- `canDiscardAsRiichiDeclarationTile`
- `service.IsTenpaiAll`

現状:

- リーチ宣言後の打牌候補確認で、候補ごとに `IsTenpaiAll` を呼ぶ。
- `IsTenpaiAll` は一般形、七対子、国士無双の向聴/聴牌判定を行う。

改善案:

- 1つの手牌に対する「各打牌後に聴牌か」をまとめて計算する helper を追加する。
- `AnalyzeShanten` の「向聴数だけ API」と合わせて、Goal 生成なしで判定する。
- 待ち集合が必要な場合は `waitsFor` の bitset helper と共有する。

期待効果:

- 合法手列挙と将来の AI 候補評価で重複計算を避けられる。

注意点:

- 現状の `mjai-tsumogiri` では呼び出し頻度は限定的だが、`mjai-manue` 本体実装後は重要度が上がる。

## 実装順序の提案

1. `AnalyzeShanten` の allocation 削減。
2. `waitsFor` の bitset 化。
3. `tile` の code / yaochu lookup table 化。
4. 和了形判定の lookup table / memo 化。
5. 役判定の map 削減。
6. 合法手列挙の distinct helper 追加。

## 検証方針

- 既存テスト:

```powershell
$env:GOEXPERIMENT='jsonv2'; go test ./...
```

- 向聴計算の性能比較:

```powershell
go test ./internal/domain/game/round/service -bench BenchmarkShantenAnalysis -benchmem
```

- 変更前後で見る指標:
  - `ns/op`
  - `B/op`
  - `allocs/op`

性能改善は挙動互換性を崩しやすいため、1回の差分では1テーマに絞る。
特に向聴計算と和了形判定は、既存の単体テストに加えてランダム手牌の旧実装比較テストを一時的に置くと安全に進めやすい。
