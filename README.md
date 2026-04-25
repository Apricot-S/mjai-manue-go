# mjai-manue-go

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Apricot-S/mjai-manue-go)

Go port of [mjai-manue](https://github.com/gimite/mjai-manue) — a Mahjong AI for the [Mjai Mahjong AI match server](https://gimite.net/pukiwiki/index.php?Mjai%20%E9%BA%BB%E9%9B%80AI%E5%AF%BE%E6%88%A6%E3%82%B5%E3%83%BC%E3%83%90)

[Sample game record of a self-match](https://apricot-s.github.io/mjai-manue-go/)

## Differences from Original

### Pipe Mode Support

- Adds standard input/output support for JSON Lines streams.
- In pipe mode, the bot reads input line by line and emits an action only when the current state requires a decision.
- `{"type":"none"}` emitted in pipe mode means an explicit pass, such as skipping a call or win. Inputs that do not require a decision produce no output.

### Architecture Improvements

- Embed configuration files at build time instead of loading them at runtime.
- Fixed an incorrect shanten number calculation when a hand contains four identical tiles.
- Log more detailed information about the game state.
- Improved error handling to more reliably reject invalid or anomalous input.
- Refactored the code for better readability and maintainability.

### Target Version

The original project includes an older version written in Ruby and a newer version written in CoffeeScript. This project ports only the new version.

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

### TCP/IP (mjsonp)

```sh
mjai-manue mjsonp://example.com:11600/default
```

Use this mode with an [mjai server](https://github.com/gimite/mjai) that expects one response for each input message. When the bot has no action to take, the TCP adapter sends `{"type":"none"}` as the protocol response.

### Pipe (JSON Lines)

```sh
mjai-manue
```

When no URL is provided, `mjai-manue` reads JSON Lines from stdin and writes protocol messages to stdout. Output is sparse: the bot writes a line only when it chooses an action. A written `{"type":"none"}` is an explicit pass, not a generic acknowledgement.

For more information, see [cmd/](cmd/).

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
