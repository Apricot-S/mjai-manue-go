# mjai-manue-go

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Apricot-S/mjai-manue-go)

Go port of [mjai-manue](https://github.com/gimite/mjai-manue) — a Mahjong AI for the [Mjai Mahjong AI match server](https://gimite.net/pukiwiki/index.php?Mjai%20%E9%BA%BB%E9%9B%80AI%E5%AF%BE%E6%88%A6%E3%82%B5%E3%83%BC%E3%83%90)

> [!NOTE]
> The original project includes an older version written in Ruby and a newer version written in CoffeeScript. This project ports only the new version.

[Sample game log of a self-match](https://apricot-s.github.io/mjai-manue-go/)

## Differences from Original

### stdio Mode Support

This project adds standard input/output support for JSON Lines streams, following the same style as [Akochan](https://github.com/critter-mj/akochan) and [Mortal](https://github.com/Equim-chan/Mortal). In stdio mode, the bot reads input line by line and emits an action only when the current state requires a decision.

### No `possible_actions` Dependency

In mjai protocol messages, `possible_actions` may be attached to events such as `tsumo` and `dahai` to tell the bot which responses are currently legal. Unlike the original project, this project does not require that field to be present. Instead, it updates the game state from the event stream and derives available decisions from that state.

This makes the bot usable with inputs that contain the game events but omit server-provided action candidates, including mjson game logs and environments that do not provide `possible_actions`, such as [RiichiEnv](https://github.com/smly/RiichiEnv).

### Embedded Configuration

Unlike the original project, this project embeds configuration files at build time. The installed binary can run on its own without depending on files in the repository checkout.

### Other Differences

- Ported the AI logic from the original implementation while reimplementing the rest in Go.
- Treats malformed or unexpected input more strictly than the original implementation.
- Fixed an incorrect shanten number calculation when a hand contains four identical tiles.
- Logs more detailed information about the game state.

## How It Works

(TODO)

The discard that minimizes this avgRank is selected.

Decisions such as "whether to call or not" and "whether to declare Riichi or not" are also made in a similar method.

## Prerequisites

This project (including all tools under [tools/](tools/)) requires:

- [Go 1.26 or later](https://go.dev/dl/)
- Environment variable `GOEXPERIMENT=jsonv2` enabled when building, installing or running with `go run`

## Installation

```sh
go install github.com/Apricot-S/mjai-manue-go/cmd/mjai-manue@latest
```

## Usage

```sh
mjai-manue [--name <PLAYER_NAME>] [--seed <INT>] [<URL>]
```

See [cmd/](cmd/) for more information.

> [!TIP]
> To customize the AI's strategic behavior, replace the following configuration files before building `mjai-manue`:
>
> - `configs/danger_tree.all.json`
> - `configs/game_stats.json`
> - `configs/light_game_stats.json`
>
> See [tools/](tools/) for instructions on how to generate these files.

## Credits

This project is a Go port of [mjai-manue](https://github.com/gimite/mjai-manue), created by [Hiroshi Ichikawa](https://github.com/gimite).

Some parts of the code are ported from [mjai](https://github.com/gimite/mjai), created by [Hiroshi Ichikawa](https://github.com/gimite).

## Licenses

This project is licensed under the [New BSD License](LICENSE) (3-Clause BSD License).

This project also contains configuration files copied from the original project:

- `configs/danger_tree.all.json`
- `configs/game_stats.json`
- `configs/light_game_stats.json`

These files are copyright Hiroshi Ichikawa and distributed under the New BSD License.
