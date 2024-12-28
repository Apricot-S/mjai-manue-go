# mjai-manue-go

Go port of [mjai-manue](https://github.com/gimite/mjai-manue)

Differences from the original:

- Supports both [Gimite's original Mjai protocol](https://gimite.net/pukiwiki/index.php?Mjai%20%E9%BA%BB%E9%9B%80AI%E5%AF%BE%E6%88%A6%E3%82%B5%E3%83%BC%E3%83%90) and a [minor modified version of the Mjai protocol](https://mjai.app/docs/mjai-protocol) used by [RiichiLab](https://mjai.app/).
- Fixed the calculation of the shanten number when the hand contains four identical tiles.
- Improved error handling to more reliably reject invalid or anomalous input.
- Refactored the code to improve readability and maintainability.

> [!NOTE]
> The original project includes an older version (written in Ruby), but this project ports only the new version (written in CoffeeScript).

## Installation

There are two options to install this application:

### Option 1: Download from releases

Download the executable file for your platform from the [releases page](https://github.com/Apricot-S/mjai-manue-go/releases/latest).

### Option 2: Build from source

```sh
go build github.com/Apricot-S/mjai-manue-go/cmd/manue
```

`mjai-manue` will be built in current directory.

## Usage

(TODO)

### For [mjai](https://github.com/gimite/mjai)

```sh
mjai-manue --tcp http://example.com:11600/default-room
```

### For [mjai.app](https://github.com/smly/mjai.app)

```sh
mjai-manue --stdio --batch
```

## License

Licensed under the [New BSD License](LICENSE) (3-Clause BSD License).

This project includes the following files from the original mjai-manue project,
which are licensed under the New BSD License. Full credit for these files goes
to the original author, [Hiroshi Ichikawa](https://github.com/gimite).

Files used:

- configs/danger_tree.all.json
- configs/game_stats.json
- configs/light_game_stats.json
