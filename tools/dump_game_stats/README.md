# dump_game_stats

This tool analyzes game logs in mjai format, including gzip-compressed files, and generates overall game statistics such as number of wins, tsumo rate, draw ratio, average winning points, and tenpai timing distributions.

## What It Does

(TODO)

## Output

The tool writes the resulting statistics as JSON to standard output.
You can redirect or pipe it as needed for downstream processing.

The output is directly usable as `configs/game_stats.json`.

## Usage

With the top-level directory of working tree of this repository as the current directory, run the following command:

```sh
go run ./tools/dump_game_stats <log_glob_patterns...> > PATH/TO/game_stats.json
```

- Replace `<log_glob_patterns...>` with one or more file path patterns matching your target logs, such as `"logs/*/*.mjson"` and `"logs/*/*.mjson.gz"`. You can specify multiple patterns, separated by spaces.

### Sample Output (formatted)

```json
{
  "numHoras": 17793,
  "numTsumoHoras": 7195,
  "numTurnsDistribution": [
    0.008629732049203776,
    ...
  ],
  "ryukyokuRatio": 0.15700390960236482,
  "averageHoraPoints": 5533.648063845332,
  "koHoraPointsFreqs": {
    "1000": 1087,
    "1100": 445,
    ...
    "32000": 1,
    "total": 12834
  },
  "oyaHoraPointsFreqs": {
    "1500": 555,
    "2000": 154,
    ...
    "36000": 1,
    "total": 4959
  },
  "yamitenStats": {
    "17,0": {
      "total": 41903,
      "tenpai": 2
    },
    ...
    "12,4": {
      "total": 4,
      "tenpai": 4
    }
  },
  "ryukyokuTenpaiStat": {
    "total": 13172,
    "tenpai": 5468,
    "noten": 7704,
    "tenpaiTurnDistribution": {
      "0": 0,
      ...
      "17.5": 46,
    }
  }
}
```
