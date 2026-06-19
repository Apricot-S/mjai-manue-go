# Self-match log generation

This directory contains the Docker Compose setup for generating README sample self-match logs with four Go `mjai-manue` clients.

## Run

From the repository root:

```powershell
scripts\self-match\run.ps1
```

or:

```sh
scripts/self-match/run.sh
```

The default output directory is `scripts/self-match/out`. You can override the runtime settings with environment variables:

| Variable    | Default   | Description                                         |
| ----------- | --------- | --------------------------------------------------- |
| `LOG_DIR`   | `./out`   | Directory mounted as the mjai server log directory. |
| `NUM_GAMES` | `1`       | Number of games to run.                             |
| `ROOM`      | `default` | mjai room name.                                     |
| `GAME_TYPE` | `tonnan`  | mjai game type passed to the server.                |
| `PORT`      | `11600`   | mjai server port.                                   |

## Bundled images

Tile images are bundled in `scripts/self-match/images`. The required file names are:

- `blank.png`
- `p_<tile>_1.gif` and `p_<tile>_3.gif` for normal tiles and tile backs
- `p_ms5r_1.png`, `p_ms5r_3.png`, `p_ps5r_1.png`, `p_ps5r_3.png`, `p_ss5r_1.png`, and `p_ss5r_3.png` for red fives

## Replace viewer images

After generating a log, rewrite the generated viewer to use local images:

```sh
python scripts/self-match/replace_assets.py scripts/self-match/out/<log>.html
```

The script copies `scripts/self-match/images` into the generated
`<log>.html.files/images` directory and rewrites
`<log>.html.files/js/archive_player.js`.

## Asset source

Most bundled tile images are based on free Mahjong assets by Mahjong kingdom., Inc:

<https://mj-king.net/sozai/>

`blank.png` and red five images are modified versions of the Mahjong kingdom assets
bundled with Mjai ([gimite/mjai](https://github.com/gimite/mjai)).

Mahjong kingdom., Inc states that:

- the assets may be used freely;
- contact is not required;
- links or source attribution are appreciated;
- modification and redistribution are allowed; and
- copyright is not waived.
