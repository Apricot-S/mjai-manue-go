# mjai-manue-go-rewrite: Agent ガイド

このリポジトリでは、設計の一次資料を `.agents/design.md` に集約しています。実装・修正は原則として設計書に従い、設計に影響する変更は同じ差分で設計書も更新してください。

## 参照すべき一次資料

- `.agents/README.md`: `.agents` 配下の文書の役割と分割方針
- `.agents/design.md`: 全体方針（ゴール/非ゴール、NFR、DDD/レイヤ、ユースケース、現状実装、実装計画、テスト戦略）
- `.agents/board-state-output.md`: 盤面状態出力の実装状況と移植元参照（恒久仕様は `.agents/design.md` に反映）
- `.agents/tools-porting.md`: `tools/` Go 移植の完了済み棚卸しメモ。生成物 schema、参照元、受け入れ条件の確認に使う。
- `.agents/terminology-en.md`: 移植元コードの英語用語を確認する補助資料
- `.agents/archive/`: 完了済みフェーズの詳細メモ。AI 本体移植の経緯確認が必要な場合だけ参照する。

## 参照実装

- `reference/repositories/mjai-manue-go-main/`: main merge 済みの過去 Go 実装。参照専用。必要な実装方針や差分確認に使う。
- `reference/repositories/mjai-manue-original/`: CoffeeScript 版オリジナル。参照専用。移植元挙動・命名・tools ロジック確認に使う。
- `reference/` 配下は編集しない。変更が必要な場合は現行コードへ反映し、参照元はそのままにする。

## 実装ルール（設計書の要点）

- 依存方向は内側へ（外側→内側）。`domain` は I/O や外部仕様に依存しない。
  - `internal/domain/`: ビジネスルール（状態/判定/意思決定の核）
  - `internal/application/`: ユースケース（受信→状態適用→意思決定→送信のオーケストレーション）
  - `internal/adapter/`: 具体 I/O（mjai TCP/stdio runtime、JSON codec 等）。外部プロトコル差分は ACL で吸収する。
  - `cmd/`: CLI エントリ。フラグ解析、Agent 選択、runtime 起動、終了コード変換のみ。
- 現在実装済みの CLI は `cmd/mjai-tsumogiri` と `cmd/mjai-manue`。AI 本体移植と `tools/` 配下の Go 移植は完了済み。
- 次の主作業は、`log` フィールドを JSON object に変換するラッパースクリプト、Original-vs-Port Testing Framework。
- `README.md` / `cmd/README.md` / `cmd/mjai-manue/README.md` / `tools/README.md` は利用者向け文書。現状判断は `.agents/design.md` と実ファイルを根拠にし、tools の CLI 仕様を変更する場合は必要に応じて `tools/*/README.md` と同期する。
- stdout は **プロトコル出力専用**。ログ/エラーは stderr（`.agents/design.md` の I/O 安全性）。
- 入力が空行・不正 JSON の場合は **エラー終了**（継続しない）。
- 送信はメッセージ単位で必ず flush する（透過性）。
- `--seed` を持つコマンドでは乱数を決定的にする（決定性）。現行 `mjai-tsumogiri` は乱数を使わず、`--seed` も持たない。
- `ai.Request` の `Round` は AI 分野でいう observation（obs）として扱う。legal actions は外側から別フィールドで渡すのではなく、obs（現状は `round.ActionStateViewer`）に含める設計を維持する。
- `round.State` / `EventApplier` / `LegalActions` は現状を最終形として扱う。責務分割目的での追加 struct/service 化は、間接参照が増えて読みにくくなるため原則行わない。
- `tools/` 配下を保守する場合は、参照元 CoffeeScript と過去 Go 実装を補助にしつつ、現行の `internal/domain` / `configs` の型と責務境界に寄せる。大量ログ処理の I/O は `tools` 側に閉じ込め、`domain` にファイル形式や集計 CLI の都合を持ち込まない。
- `log` フィールド JSON 化は mjai-manue 本体に組み込まず、出力を書き換える外部ラッパーとして追加する。本体の CLI オプションを増やさず、後処理の参考実装としても使える形を優先する。
- Original-vs-Port Testing Framework は通常 CI へ無理に組み込まず、差分調査・移植確認用の開発者向け検証基盤として設計する。

## ドキュメント更新の運用

- 仕様/設計に影響する変更（パッケージ構造、責務分担、主要インタフェース、I/O 方針、実装計画、テスト方針など）を入れる場合は、該当箇所の `.agents/design.md` も同時に更新する。
- プロトコル（メッセージ種別・フィールド・解釈）を変更/追加する場合は、adapter codec とその単体テストを更新し、設計判断に影響する場合は `.agents/design.md` も更新する。
- `.agents` 配下に新しい設計メモや計画文書を追加する場合は、`.agents/README.md` に役割を追記する。
- `.agents/design.md` の見出し番号（`## 4.x` など）は参照の手掛かりになるため、可能な限り維持する（大きく組み替える場合は最小限の破壊に留める）。

## テスト方針（抜粋）

- `domain` の純粋ロジックはテーブル駆動で単体テスト（TDD）。
- プロトコル入出力はゴールデンテストで「action のみ」を比較（詳細は `.agents/design.md` のテスト章を参照）。
- Go コードを変更した場合は、`go fix ./...` と `go vet ./...` を実行し、必要に応じて `GOEXPERIMENT=jsonv2` を有効化した `go test ./...` も実行する。
- `encoding/json/v2` を使うテストを実行する際は、実験機能のため `GOEXPERIMENT=jsonv2` を有効化する（例: PowerShell なら `$env:GOEXPERIMENT='jsonv2'; go test ./...`）。
- `t.Fatal/t.Fatalf` は「この時点でテスト継続が不可能」なケース（前提条件の破綻、初期化エラー、`nil` により以降がパニックする等）に限定する。
- 値の不一致など「継続して他の差分も報告できる」ケースは `t.Error/t.Errorf` を優先し、常に `Fatal` を使う書き方は避ける。
