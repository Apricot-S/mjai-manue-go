# mjai-manue-go

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Apricot-S/mjai-manue-go)

Go port of [mjai-manue](https://github.com/gimite/mjai-manue)

For more information, see [the original README (Japanese)](https://github.com/gimite/mjai-manue/blob/master/README.md) or [its translation](docs/README-en.md).

## Differences from Original

### Protocol Support Extensions

- In addition to [Gimite's original Mjai protocol](https://gimite.net/pukiwiki/index.php?Mjai%20%E9%BA%BB%E9%9B%80AI%E5%AF%BE%E6%88%A6%E3%82%B5%E3%83%BC%E3%83%90), also supports [a minor modified version of the Mjai protocol used by RiichiLab](https://mjai.app/docs/mjai-protocol).

### Architecture Improvements

- The configuration files is embedded at build time rather than loaded at runtime.
- Fixed the miscalculation of the shanten number when the hand contains four identical tiles.
- Improved error handling to more reliably reject invalid or anomalous input.
- Refactored the code to improve readability and maintainability.

### Target Version

> [!NOTE]
> The original project includes an older version written in Ruby and a newer version written in CoffeeScript. This project ports only the new version.

## Installation

Requires [Go 1.24 or later](https://go.dev/dl/).

```sh
go install github.com/Apricot-S/mjai-manue-go/cmd/mjai-manue@latest
```

## Usage

### For TCP/IP (e.g., [mjai](https://github.com/gimite/mjai))

```sh
mjai-manue mjsonp://example.com:11600/default
```

### For Standard I/O (e.g., [mjai.app](https://github.com/smly/mjai.app))

```sh
mjai-manue
```

> [!NOTE]
> In practice, `mjai.app` runs `bot.py` in the submission `.zip` file.
> You need to call the above command from within `bot.py` and pipe the standard input and output.

> [!TIP]
> See [scripts/mjai.app/](scripts/mjai.app/) for how to generate a submission file for `mjai.app`.

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

## Licenses

This project is licensed under the [New BSD License](LICENSE) (3-Clause BSD License).

This project also contains configuration files copied from the original project:

- `configs/danger_tree.all.json`
- `configs/game_stats.json`
- `configs/light_game_stats.json`

These files are copyright Hiroshi Ichikawa and distributed under the New BSD License.
