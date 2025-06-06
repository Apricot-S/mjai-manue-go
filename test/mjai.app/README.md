# Self-match test using [mjai.app](https://github.com/smly/mjai.app)

This directory contains the `test.sh` script for testing the functionality of the `mjai-app.zip` file generated by [/scripts/mjai.app/build.sh](/scripts/mjai.app/build.sh) before actually submitting it to RiichiLab. To execute this script, with the top-level directory of working tree of this repository as the current directory, run the following command:

```sh
test/mjai.app/test.sh PATH/TO/mjai-app.zip
```

This script performs repeated self-matches using four replicas of the model bundled in `mjai-app.zip` to check for any errors. Once this script starts running, it will not stop unless forcibly halted, whether by repeatedly pressing the `^C` key sequence or sending a `KILL` signal. Logs for each self-match are output to a directory named `logs.YYYY-MM-DD-hh-mm-ss`, where `YYYY-MM-DD-hh-mm-ss` is the start time of each self-match. Log directories that are determined to be clearly error-free are automatically deleted.
