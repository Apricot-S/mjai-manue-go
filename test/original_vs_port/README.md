# Original-vs-Port

## Usage

### gen_log

With `test/original_vs_port/gen_log/` as the current directory, run the following command:

```sh
NUM_GAMES=<NUMBER_OF_GAMES> docker compose up
```

### verify

With the top-level directory of working tree of this repository as the current directory, run the following command:

```sh
go run ./test/original_vs_port/verify <LOG_GLOB_PATTERNS>... > <OUTPUT_FILEPATH>
```
