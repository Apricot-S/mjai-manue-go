#!/usr/bin/env bash

HOST=${HOST:-server}
PORT=${PORT:-11600}
ROOM=${ROOM:-default}
NUM_GAMES=${NUM_GAMES:-1}

GAME_TYPE=${GAME_TYPE:-tonnan}

mjai server \
    --host=$HOST \
    --port=$PORT \
    --game_type=$GAME_TYPE \
    --room=$ROOM \
    --games=$NUM_GAMES \
    --log_dir="/log_dir"

echo "all $NUM_GAMES games done"
