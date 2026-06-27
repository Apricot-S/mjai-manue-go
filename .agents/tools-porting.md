# Tools Porting Notes

この文書は、未移植の `tools/` を現行 Go 実装へ移植するための棚卸しメモである。恒久方針は `.agents/design.md` を一次資料とし、この文書は実装順・参照元・受け入れ条件を整理する。

## 1. 対象

`tools/` は `configs/` に埋め込む JSON と、その生成過程で使う中間データを作る開発者向けコマンド群である。通常の `mjai-manue` 実行時依存にはしない。

| Tool | 生成物 | 現行 consumer |
| --- | --- | --- |
| `dump_game_stats` | `game_stats.json` | `configs.LoadGameStats()` / `ai.ManueStats` |
| `print_game_stats` | なし | `game_stats.json` の目視確認 |
| `dump_light_game_stats` | 中間 `scoreStats` JSON | `postprocess_light_game_stats` |
| `postprocess_light_game_stats` | `light_game_stats.json` | `configs.LoadGameStats()` / `ai.RankStats` |
| `estimate_danger dump_tree_json` | `danger_tree.all.json` | `configs.LoadDangerTree()` / `ai.NewDangerEstimator` |

## 2. 参照元

- CoffeeScript 版:
  - `reference/repositories/mjai-manue-original/coffee/dump_game_stats.coffee`
  - `reference/repositories/mjai-manue-original/coffee/dump_light_game_stats.coffee`
  - `reference/repositories/mjai-manue-original/coffee/postprocess_light_game_stats.coffee`
  - `reference/repositories/mjai-manue-original/coffee/print_game_stats.coffee`
  - `reference/repositories/mjai-manue-original/tools/estimate_danger.rb`
  - `reference/repositories/mjai-manue-original/coffee/danger_estimator.coffee` は tool 本体ではなく、Scene 実装の補助参照に限定する。
- 過去 Go 実装:
  - `reference/repositories/mjai-manue-go-main/tools/shared/`
  - `reference/repositories/mjai-manue-go-main/tools/dump_game_stats/`
  - `reference/repositories/mjai-manue-go-main/tools/dump_light_game_stats/`
  - `reference/repositories/mjai-manue-go-main/tools/postprocess_light_game_stats/`
  - `reference/repositories/mjai-manue-go-main/tools/print_game_stats/`
  - `reference/repositories/mjai-manue-go-main/tools/estimate_danger/`

過去 Go 実装は集計ロジックの補助資料として使う。ただし package 構造が現行と異なるため、`internal/game` / `internal/protocol` への依存はそのまま持ち込まない。

## 3. 現行 schema

### 3.1 `game_stats.json`

現行 loader は `configs.GameStats` へ unmarshal する。

- `numHoras`
- `numTsumoHoras`
- `numTurnsDistribution`
- `ryukyokuRatio`
- `averageHoraPoints`
- `koHoraPointsFreqs`
- `oyaHoraPointsFreqs`
- `yamitenStats`
- `ryukyokuTenpaiStat`

AI 側では `WinScoreStats` / `RoundEndStats` / `DrawTenpaiStats` / `TenpaiEstimatorStats` / `DealInStats` として読む。

### 3.2 `light_game_stats.json`

現行 loader は `configs.LightGameStats` へ unmarshal し、`GameStats.LightGameStats` に合成する。

- `winProbsMap`

AI 側では `RankStats.RelativeWinProbs` として読む。`dump_light_game_stats` の `scoreStats` は中間形式であり、`configs` には直接置かない。

### 3.3 `danger_tree.all.json`

現行 loader は `configs.DecisionNode` へ unmarshal する。

- `average_prob`
- `conf_interval`
- `num_samples`
- `feature_name`
- `negative`
- `positive`

AI 側では `ai.DangerTreeNode` として読み、`feature_name == null` を leaf として扱う。

## 4. 移植順

1. `tools/internal/archive`（実装済み）
   - glob、gzip、JSON Lines 読み取り、mjai inbound decode、必要に応じた `round.State` 更新を提供する。
   - 空行・不正 JSON の扱いは runtime と同じく error を基本にする。過去 Go 実装は空行を skip していたため、互換性が必要かは実装時に明示する。
   - Archive は `[]paths` を保持せず、`PlayPaths(paths, handlers)` で受け取る。進捗表示は `OnFileDone` callback を使って CLI 側で管理し、JSON 生成 stdout を汚さない。

2. `postprocess_light_game_stats`（実装済み）
   - domain 依存がなく、`configs.LightGameStats` の schema 互換を最初に固定しやすい。
   - 小さい `scoreStats` fixture で `winProbsMap` を golden 化する。

3. `print_game_stats`（実装済み）
   - domain 依存は turn range だけ。必要なら定数を tool 側へ閉じ込める。
   - `configs.GameStats` の読み取り互換を確認する smoke test を置く。

4. `dump_light_game_stats`（実装済み）
   - `start_game` / `start_kyoku` / `hora` / `ryukyoku` / `end_game` の score 推移だけを使うため、`dump_game_stats` より先に実装する。
   - chicha からの相対席 `0..3` と `E1..S4` key を固定する。

5. `dump_game_stats`（実装済み）
   - `round.State` と tenpai 判定を使うため、shared の state 更新を先に固める。
   - `yamitenStats`、`ryukyokuTenpaiStat` は off-by-one が出やすいので小さい fixture で代表ケースを固定する。

6. `estimate_danger`
   - サブコマンドは `dump_tree_json` までの生成経路を優先する。
   - tool 本体の一次参照は `reference/repositories/mjai-manue-original/tools/estimate_danger.rb` とする。`coffee/danger_estimator.coffee` は Scene feature 判定の補助参照としてのみ使う。
   - Scene は現行 `internal/domain/ai` の danger scene と似ているが、オリジナル tool と微妙に判定が異なる箇所がある。`estimate_danger` 移植では現行 AI 側へ寄せず、オリジナル tool の判定差分をそのまま再現する。
   - Scene 差分の背景確認には <https://github.com/gimite/mjai-manue/issues/2> を参照する。
   - `interesting_graph` は gnuplot 依存があるため optional とし、通常テスト対象から外す。

### 4.1 `estimate_danger` PR 分割計画

`estimate_danger` は PR 単位を小さく保つため、1 PR につきサブコマンドを 1 つだけ有効化する。共通コードは、その PR のサブコマンド実行に必要な最小範囲だけ同梱する。後続サブコマンド用の実装を先に書いた場合は、通常 PR へ混ぜず、別 branch などに退避してから必要な PR へ順に取り込む。

| PR | サブコマンド | 目的 | 進捗 |
| --- | --- | --- | --- |
| 1 | `extract` | Mjai log から feature gob を生成する。Scene は `internal/domain/ai/danger_scene.go` コピーを起点に Ruby tool 差分だけ修正し、Ruby 版と CoffeeScript/runtime 版との差分が別物であることをコメントに残す。 | 実装済み |
| 2 | `tree` | `features.gob` から probability 集計と `configs.DecisionNode` 互換の決定木 gob を生成する。 | 実装済み |
| 3 | `dump_tree_json` | tree gob を `configs/danger_tree.all.json` 互換 JSON へ変換する。 | 実装済み |
| 4 | `dump_tree` | 保存済み tree gob を text tree として表示する。 | 実装済み |
| 5 | `single` | feature ごとの true/false 危険率、信頼区間、sample 数を表示する。 | 実装済み |
| 6 | `interesting` | Ruby の interesting criteria を移植し、必要に応じて probability map gob を保存する。 | 未着手 |
| 7 | `benchmark` | `interesting` と同じ criteria builder で probability map 作成までを実行する。 | 未着手 |
| 8 | `interesting_graph` | `interesting` の probability gob から points / plot / png / html を生成する。gnuplot 依存のため通常テスト対象外。 | 未着手 |

## 5. 受け入れ条件

- `go run ./tools/<name>` が README の usage と一致する。
- stdout は生成物または表示結果専用、進捗・診断は stderr。
- 生成 JSON は `encoding/json/v2` で現行 `configs` 型へ unmarshal できる。
- Go コードを追加した差分では `go fix ./...`、`go vet ./...`、必要に応じて `GOEXPERIMENT=jsonv2` 付き `go test ./...` を実行する。
