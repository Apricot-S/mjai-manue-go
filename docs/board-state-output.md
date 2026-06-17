# Board State Output Plan

## 目的

Ruby版 `mjai` の `Game.render_board()` 相当として、状態更新後の局面を stderr へ出力できるようにする。

## 方針

- `round.State` は I/O を持たず、盤面の整形済み文字列を返す pure method を提供する。
- `application.Bot` は状態遷移の直後に reporter を呼び出す。
- stderr への実書き込みは adapter / runtime 側で `io.Writer` として注入する。
- mjai の受信/送信 JSON trace は盤面出力と別責務として `mjairuntime` で扱う。
- `end_game` は Bot に渡さない。`end_kyoku` 相当の `EndRound` は type 以外の局面情報を持たないため、盤面出力のために round へ Apply しない。

## レイヤ分担

- `internal/domain/game/round`
  - `RenderBoard() string`
  - `BoardRenderer` interface
  - ForTest 構築 state を使って整形ロジックを単体テストする。
- `internal/application`
  - `Reporter` を `Bot` へ明示引数で注入する。
  - `StartRound` 後、通常 event の `Apply` 成功後に reporter を呼ぶ。
  - reporter がない場合は従来動作を維持する。
- `internal/adapter/mjai/runtime`
  - `io.Writer` へ盤面文字列を書く reporter 実装を持つ。
  - 受信/送信 JSON を `"<-\t..."` / `"->\t..."` 形式で trace 出力する。
- `cmd`
  - stdout は protocol output 専用のまま、stderr を runtime の trace / board output writer として渡す。

## TODO

- [x] 盤面出力の pure method と最小フォーマットを追加する。
- [x] Bot へ reporter 注入を追加する。
- [x] runtime で stderr writer に接続する。
- [x] 受信/送信 JSON trace を runtime に追加する。
- [ ] Ruby版 mjai の `Game.render_board()` と表示項目・順序をさらに照合する。
- [ ] `round.State.Apply` の未対応 event を埋め、event 列ベースの盤面 golden test を増やす。

## 参考 mjai の `Game.render_board()` の実装抜粋

```ruby
module Mjai
    class Game
        def render_board()
          result = ""
          if @bakaze && @kyoku_num && @honba
            result << ("%s-%d kyoku %d honba  " % [@bakaze, @kyoku_num, @honba])
          end
          result << ("pipai: %d  " % self.num_pipais) if self.num_pipais
          result << ("dora_marker: %s  " % @dora_markers.join(" ")) if @dora_markers
          result << "\n"
          @players.each_with_index() do |player, i|
            if player.tehais
              result << ("%s%s%d%s tehai: %s %s\n" %
                   [player == @actor ? "*" : " ",
                    player == @oya ? "{" : "[",
                    i,
                    player == @oya ? "}" : "]",
                    Pai.dump_pais(player.tehais),
                    player.furos.join(" ")])
              if player.reach_ho_index
                ho_str =
                    Pai.dump_pais(player.ho[0...player.reach_ho_index]) + "=" +
                    Pai.dump_pais(player.ho[player.reach_ho_index..-1])
              else
                ho_str = Pai.dump_pais(player.ho)
              end
              result << ("     ho:    %s\n" % ho_str)
            end
          end
          result << ("-" * 80) << "\n"
          return result
        end
```

```ruby
module Mjai
    class Pai
        def self.dump_pais(pais)
          return pais.map(){ |pai| "%-3s" % pai }.join("")
        end
```
