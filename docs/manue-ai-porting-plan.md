# Manue AI Porting Plan

この文書は CoffeeScript 版 Manue AI を Go へ再実装する作業順だけを扱う。原仕様は `docs/manue-ai-original-spec.md`、Go 移植用の責務分割と接続仕様は `docs/manue-ai-porting-spec.md` を一次資料にする。

## 1. 前提

- `reference/repositories/mjai-manue-original/coffee/manue_ai.coffee` を AI ロジックの一次資料にする。
- 過去 Go 実装と現行 `internal/domain/ai` の Manue 固有コードは補助資料に留める。
- `agent.go` と `tsumogiri_agent.go` は維持する。Manue 固有コードだけを再実装対象にする。
- docs cleanup は AI 再実装の前提作業であり、実装しながら仕様を後追いで書く進め方には戻さない。

## 2. 実装順

1. **仕様固定**
   - `docs/manue-ai-original-spec.md` の characterization ケース候補を、実装前にテスト候補へ落とし込む。
   - 原仕様と Go 移植仕様の差分が必要な箇所は、実装前に `docs/manue-ai-porting-spec.md` へ理由付きで追記する。

2. **Manue 固有コードの整理**
   - 現行 `internal/domain/ai` の Manue 固有実装を、仕様に合う純粋関数と破棄対象へ分ける。
   - `agent.go` と `tsumogiri_agent.go` は残す。
   - 破棄対象を流用するために設計を曲げない。

3. **Agent dispatcher と候補生成**
   - `ManueAgent.Decide` を薄い dispatcher にする。
   - 和了優先、通常手番、副露反応の入口を分ける。
   - legal actions から打牌、reach+discard、pass、副露後打牌候補を生成する。

4. **候補比較と trace**
   - 平均順位、期待点、赤牌回避の比較規則を先に固定する。
   - trace table と `decidedKey` 相当の整形を I/O なしの formatter として実装する。

5. **score/rank 分布**
   - 自分和了、流局、他家和了、直後放銃の score delta 分布を実装する。
   - 期待点と平均順位を candidate evaluation へ接続する。

6. **和了推定**
   - 見えている牌から wall を作る。
   - `service.AnalyzeShanten` の goal と throwable vector を使い、候補別 win metrics を作る。
   - 点数計算は現行 `domain/game/round/service` を使い、原仕様との差分は characterization で確認する。

7. **危険度推定**
   - danger tree interface と scene builder を実装する。
   - 安全牌 shortcut、聴牌確率、feature evaluator を接続する。
   - `safeProb`、`dealInProb`、`immediateScoreDeltaDist` を実値化する。

8. **副露判断の完成**
   - `none` と副露候補を同じ評価経路で比較する。
   - 副露後の最善打牌評価を接続し、選択 action は副露 action のみ返す。

9. **CLI 接続と検証**
   - `cmd/mjai-manue` から stats と danger tree を load し、deps validation 済みの `ManueAgent` を生成する。
   - action golden と characterization test を追加する。
   - PowerShell で `$env:GOEXPERIMENT='jsonv2'; go test ./...` を実行する。

## 3. テスト方針

- 純粋関数はテーブル駆動テストにする。
- Agent 判断は mjai JSON Lines から `round.State` を構築して action を比較する。
- runtime golden は action のみ比較する。
- trace/log は action golden とは別 fixture にする。
- original-vs-port 比較は CI には入れず、差分調査の補助として使う。
