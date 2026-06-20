# Tools to generate stats and decision trees from game logs

Tools that analyze game logs and generate structured data for use in mjai-manue's decision-making.

These tools output JSON suitable for use as configuration files under [../configs/](../configs/).

## Deal-in risk decision tree

| Tool                                | Output                 | Purpose                                                                |
| ----------------------------------- | ---------------------- | ---------------------------------------------------------------------- |
| [estimate_danger](estimate_danger/) | `danger_tree.all.json` | Generates a decision tree to estimate deal-in risk based on game state |

## Game-level statistics

| Tool                                  | Output            | Purpose                                             |
| ------------------------------------- | ----------------- | --------------------------------------------------- |
| [dump_game_stats](dump_game_stats/)   | `game_stats.json` | Aggregates per-game metrics from logs               |
| [print_game_stats](print_game_stats/) | —                 | Displays game stats JSON in a human-readable format |

## Round-level statistics

| Tool                                                          | Output                  | Purpose                                            |
| ------------------------------------------------------------- | ----------------------- | -------------------------------------------------- |
| [dump_light_game_stats](dump_light_game_stats/)               | (intermediate JSON)     | Extracts round-level score differentials from logs |
| [postprocess_light_game_stats](postprocess_light_game_stats/) | `light_game_stats.json` | Converts score differentials into win rates        |

See each tool's `README.md` for details.
