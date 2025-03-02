# mjai-manue-go

Go port of [mjai-manue](https://github.com/gimite/mjai-manue)

***Work in progress***

Differences from the original:

- Supports both [Gimite's original Mjai protocol](https://gimite.net/pukiwiki/index.php?Mjai%20%E9%BA%BB%E9%9B%80AI%E5%AF%BE%E6%88%A6%E3%82%B5%E3%83%BC%E3%83%90) and [a minor modified version of the Mjai protocol used by RiichiLab](https://mjai.app/docs/mjai-protocol).
- Fixed the miscalculation of the shanten number when the hand contains four identical tiles.
- Improved error handling to more reliably reject invalid or anomalous input.
- Refactored the code to improve readability and maintainability.

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

### For [mjai](https://github.com/gimite/mjai) (Line-by-line, TCP/IP)

```sh
mjai-manue http://example.com:11600/default
```

### For [mjai.app](https://github.com/smly/mjai.app) (Batch, Standard I/O)

```sh
mjai-manue --batch --stdio
```

In practice, `mjai.app` runs `bot.py` in the submission zip file.
You need to call the above command from within `bot.py` and pipe the standard input and output.

## License

Licensed under the [New BSD License](LICENSE) (3-Clause BSD License).

This project includes the following files from the original mjai-manue project,
which are licensed under the New BSD License. Full credit for these files goes
to the original author, [Hiroshi Ichikawa](https://github.com/gimite).

Files used:

- configs/danger_tree.all.json
- configs/game_stats.json
- configs/light_game_stats.json
