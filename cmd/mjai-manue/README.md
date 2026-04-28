# mjai-manue

`mjai-manue` is the main AI-powered agent. It uses the common command-line modes documented in [`../README.md`](../README.md).

## Usage

Install:

```sh
go install github.com/Apricot-S/mjai-manue-go/cmd/mjai-manue@latest
```

Run:

```sh
# stdio mode
mjai-manue [--name <PLAYER_NAME>] [--seed <INT>]

# mjsonp TCP client mode
mjai-manue [--name <PLAYER_NAME>] [--seed <INT>] mjsonp://example.com:11600/default
```

The default player name is `"Manue030"`.

## Random seed

`--seed <INT>` changes the random seed. Use it when reproducible decisions with a non-default seed are required, such as golden tests or comparisons with a fixed input stream.

When `--seed` is omitted, the default seed is `0`.

The random sequence is deterministic, but it does not match the original CoffeeScript implementation.

## Configuration files

`mjai-manue` embeds configuration files at build time. It does not replace configuration paths at runtime.

To customize the AI's strategic behavior, replace these files before building `mjai-manue`:

- `configs/danger_tree.all.json`
- `configs/game_stats.json`
- `configs/light_game_stats.json`

With the repository root as the current directory, build `mjai-manue` after replacing the files:

```sh
go build ./cmd/mjai-manue
```

See [`../../tools/`](../../tools/) for instructions on generating these configuration files.
