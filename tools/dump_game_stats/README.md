# dump_game_stats

This tool analyzes game logs in Mjai format, including gzip-compressed files, and generates overall game statistics such as number of wins, Tsumo rate, draw ratio, average winning points, and Tenpai timing distributions.

## What It Does

- Parses each game log and replays all actions in order
- Tracks overall stats like number of rounds, Ryukyokus, and win-related totals
- Computes distribution of round lengths
- Measures winning point distributions by dealer and non-dealer status
- Aggregates counts of Yamiten cases (i.e. situations where a player in Tenpai does not declare Riichi and quietly remains in Tenpai) grouped by turn number and number of melds, limited to the player has not declared Riichi
- Checks for each player whether they were in Tenpai at the time of Ryukyoku, and records the turn they first entered Tenpai

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
