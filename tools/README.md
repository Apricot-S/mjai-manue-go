# Tools to generate stats and decision trees from game records

## dump_light_game_stats

This tool analyzes game logs in mjai format, including gzip-compressed files, and generates per-round statistics on score differentials.

### What It Does

- Processes each game log and examines every round
- For each player, calculates the score difference between their starting and ending points
- Groups results by the player's seating position relative to the starting dealer (chicha)
- Aggregates and outputs the frequency of each differential

### Output

The tool writes the resulting statistics as JSON to standard output.
You can redirect or pipe it as needed for downstream processing.

> [!IMPORTANT]
> The output from this tool **cannot** be used directly as `configs/light_game_stats.json`.
> You must run it through `postprocess_light_game_stats` to convert it into the proper format.

### Usage

With the top-level directory of working tree of this repository as the current directory, run the following command:

```sh
go run ./tools/dump_light_game_stats <log_glob_patterns...> > <light_game_stats.json>
```

- Replace `<log_glob_patterns...>` with one or more file path patterns matching your target logs, such as `"logs/*/*.mjson"` and `"logs/*/*.mjson.gz"`. You can specify multiple patterns, separated by spaces.
- Output is written to `<light_game_stats.json>` in JSON format.

#### Sample Output (formatted)

```json
{
  "scoreStats": {
    "E1,0": {
      "0": 12,
      "1000": 5,
      ...
      "-8000": 3
      ...
    },
    "E1,1": {
      ...
    },
    ...
  }
}
```

## postprocess_light_game_stats
