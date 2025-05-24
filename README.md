# mjai-manue-go

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Apricot-S/mjai-manue-go)

Go port of [mjai-manue](https://github.com/gimite/mjai-manue)

> 🚧 **Work in Progress**
>
> mjai-manue-go is currently in active development and not usable yet.

Differences from the original:

- Supports both [Gimite's original Mjai protocol](https://gimite.net/pukiwiki/index.php?Mjai%20%E9%BA%BB%E9%9B%80AI%E5%AF%BE%E6%88%A6%E3%82%B5%E3%83%BC%E3%83%90) and [a minor modified version of the Mjai protocol used by RiichiLab](https://mjai.app/docs/mjai-protocol).
- Fixed the miscalculation of the shanten number when the hand contains four identical tiles.
- Improved error handling to more reliably reject invalid or anomalous input.
- Refactored the code to improve readability and maintainability.
- Does not include tools to generate stats and decision trees from game records, as there is no motivation to modify the pre-generated files provided in the original project.

> [!NOTE]
> The original project includes an older version written in Ruby and a newer version written in CoffeeScript. This project ports only the new version.

For more information, see [the original README (Japanese)](https://github.com/gimite/mjai-manue/blob/master/README.md) or [its translation](docs/README-en.md).

## Installation

There are two options to install this application:

### Option 1: Download from releases

Download the executable file for your platform from the [releases page](https://github.com/Apricot-S/mjai-manue-go/releases/latest).

### Option 2: Build from source

```sh
go install github.com/Apricot-S/mjai-manue-go/cmd/mjai-manue
```

`mjai-manue` will be built in current directory.

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

## License

Licensed under the [New BSD License](LICENSE) (3-Clause BSD License).

## Copyright and Attribution

This project is a Go port of [mjai-manue](https://github.com/gimite/mjai-manue), originally created by [Hiroshi Ichikawa](https://github.com/gimite).

The following configuration files are copied from the original project and their copyright belongs to Hiroshi Ichikawa. These files are distributed under the New BSD License:

- `configs/danger_tree.all.json`
- `configs/game_stats.json`
- `configs/light_game_stats.json`
