# estimate_danger

> 🚧 **Work in Progress**
>
> estimate_danger is currently in development and not usable yet.

This tool analyzes game logs in Mjai format and generates a decision tree to estimate deal-in risk based on game state.

> [!IMPORTANT]
> Unlike the original implementation, this tool only supports the Mjai format.
> `mjlog` format is not supported.

## extract

The `extract` command extracts feature vectors from Mjai format game logs and generates training data for decision tree learning.
It focuses specifically on **situations where exactly one player has declared Riichi**, analyzing the safety of each tile discarded by the other players.

### Usage

With the top-level directory of working tree of this repository as the current directory, run the following command:

```sh
go run ./tools/estimate_danger extract -o OUTPUT_FILEPATH [OPTIONS] <PATH/TO/INPUT_FILES>
```

Required Option

- `-o OUTPUT_FILEPATH`  
Path to the output file for the extracted feature data (in `.gob` format)

Optional Flags

- `-v`  
Verbose mode (prints feature vectors for each discard candidate to standard output)
- `-filter FILTER_SPEC`  
Filters extracted scenes by feature conditions and prints only matching candidates to standard output.  
FILTER_SPEC format: `feature1:1&feature2:0&hit:1` where conditions are joined by `&`, each condition is `key:value`, and values are `1` (true) or `0` (false).
Supports any feature name defined in `Scene` struct (e.g., suji, urasuji, visible>=3, dora) plus `hit` for actual deal-in results.
- `-start FILEPATH`  
Start processing from the specified file
- -n NUMBER  
Limit the number of files to process

### What It Does

- Identifies discard situations after a Riichi declaration by another player; excludes cases with multiple Riichi declarations
- Evaluates feature vectors for each discard candidate
- Determines whether each discard candidate would deal into the Riichi player's hand
- Stores feature vectors and hit information in `.gob` format

### Output

The output file includes the following information:

- Metadata: List of feature names
- Candidate Data: Feature vectors and hit information for each discard candidate in each situation

### Example of Usage

```sh
# Basic usage
go run ./tools/estimate_danger extract -o features.gob logs/*.mjson

# Process only the first 100 files in verbose mode
go run ./tools/estimate_danger extract -v -n 100 -o features.gob logs/*.mjson

# Start processing from a specific file
go run ./tools/estimate_danger extract -start logs/game_050.mjson -o features.gob logs/*.mjson
```

### Sample Output (formatted)

(TODO)
