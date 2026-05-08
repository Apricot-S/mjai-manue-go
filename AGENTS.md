# mjai-manue-go-rewrite: Agent ガイド

このリポジトリでは、設計の一次資料を `docs/design.md` に集約しています。実装・修正は原則として設計書に従い、設計に影響する変更は同じ差分で設計書も更新してください。

## 参照すべき一次資料

- `docs/README.md`: docs 配下の文書の役割と分割方針
- `docs/design.md`: 全体方針（ゴール/非ゴール、NFR、DDD/レイヤ、ユースケース、現状実装、実装計画、テスト戦略）
- `docs/protocols.md`: mjai / stdio / RiichiLab 各プロトコルの **メッセージ仕様の集約先**
  - `docs/design.md` 側に仕様詳細を重複記載せず、必要箇所から参照する方針
- `docs/board-state-output.md`: 盤面状態出力の実装状況と移植元参照（恒久仕様は `docs/design.md` に反映）
- `docs/terminology-en.md`: 移植元コードの英語用語を確認する補助資料

## 実装ルール（設計書の要点）

- 依存方向は内側へ（外側→内側）。`domain` は I/O や外部仕様に依存しない。
  - `internal/domain/`: ビジネスルール（状態/判定/意思決定の核）
  - `internal/application/`: ユースケース（受信→状態適用→意思決定→送信のオーケストレーション）
  - `internal/adapter/`: 具体 I/O（mjai TCP/stdio runtime、JSON codec 等）。外部プロトコル差分は ACL で吸収する。
  - `cmd/`: CLI エントリ。フラグ解析、Agent 選択、runtime 起動、終了コード変換のみ。
- 現在実装済みの CLI は `cmd/mjai-tsumogiri`。`cmd/mjai-manue` は README のみ先行しており、本体実装は今後の計画対象。
- `README.md` / `cmd/README.md` / `cmd/mjai-manue/README.md` は、あえて完成形を先行して記載している利用者向け文書。実装状況と逐次同期する作業は行わず、現状判断は `docs/design.md` と実ファイルを根拠にする。
- stdout は **プロトコル出力専用**。ログ/エラーは stderr（`docs/design.md` の I/O 安全性）。
- 入力が空行・不正 JSON の場合は **エラー終了**（継続しない）。
- 送信はメッセージ単位で必ず flush する（透過性）。
- `--seed` を持つコマンドでは乱数を決定的にする（決定性）。現行 `mjai-tsumogiri` は乱数を使わず、`--seed` も持たない。
- `ai.Request` の `Round` は AI 分野でいう observation（obs）として扱う。legal actions は外側から別フィールドで渡すのではなく、obs（現状は `round.ActionStateViewer`）に含める設計を維持する。
- `round.State` / `EventApplier` / `LegalActions` は現状を最終形として扱う。責務分割目的での追加 struct/service 化は、間接参照が増えて読みにくくなるため原則行わない。
- `service/yaku.go` の TODO はオリジナル実装由来のコメントを保持しているもの。未完了作業とは扱わず、必要ならその旨をコメントに追記する。

## リファクタリング作業順

レビュー負荷を抑えるため、以下は1件ずつ独立した差分で進める。

1. viewer / player getter の防御コピー
2. `hand.VisibleHand.Call` / `hand.InvisibleHand.Call` の未知 meld 防止
3. mjai runtime の JSON Lines loop 共通化
4. `VisiblePlayer` / `InvisiblePlayer` の小規模重複整理（読みにくくなるなら実施しない）
5. `service/yaku.go` の TODO 注釈整理（ロジック再設計や panic 排除はしない）
6. 内部 API と非自明処理へのコメント追加。不変条件・コーナーケース・移植仕様上の注意点に絞り、自明な getter や単純処理には増やしすぎない。

## ドキュメント更新の運用

- 仕様/設計に影響する変更（パッケージ構造、責務分担、主要インタフェース、I/O 方針、実装計画、テスト方針など）を入れる場合は、該当箇所の `docs/design.md` も同時に更新する。
- プロトコル（メッセージ種別・フィールド・解釈）を変更/追加する場合は、まず `docs/protocols.md` を更新し、`docs/design.md` から参照する（重複記載しない）。
- docs 配下に新しい設計メモや計画文書を追加する場合は、`docs/README.md` に役割を追記する。
- `docs/design.md` の見出し番号（`## 4.x` など）は参照の手掛かりになるため、可能な限り維持する（大きく組み替える場合は最小限の破壊に留める）。

## テスト方針（抜粋）

- `domain` の純粋ロジックはテーブル駆動で単体テスト（TDD）。
- プロトコル入出力はゴールデンテストで「action のみ」を比較（詳細は `docs/design.md` のテスト章を参照）。
- `encoding/json/v2` を使うテストを実行する際は、実験機能のため `GOEXPERIMENT=jsonv2` を有効化する（例: PowerShell なら `$env:GOEXPERIMENT='jsonv2'; go test ./...`）。
- `t.Fatal/t.Fatalf` は「この時点でテスト継続が不可能」なケース（前提条件の破綻、初期化エラー、`nil` により以降がパニックする等）に限定する。
- 値の不一致など「継続して他の差分も報告できる」ケースは `t.Error/t.Errorf` を優先し、常に `Fatal` を使う書き方は避ける。
