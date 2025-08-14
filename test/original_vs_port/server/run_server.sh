#!/usr/bin/env bash

HOST=${HOST:-server}
PORT=${PORT:-11600}
ROOM=${ROOM:-default}
NUM_GAMES=${NUM_GAMES:-1}

GAME_TYPE=${GAME_TYPE:-tonnan}

for i in $(seq 1 "$NUM_GAMES"); do
    echo "$i/$NUM_GAMES games"

    mjai server \
        --host=0.0.0.0 \
        --port=$PORT \
        --game_type=$GAME_TYPE \
        --room=$ROOM \
        --games=$NUM_GAMES \
        --log_dir="/log_dir"

    echo "finished game $i"
done

echo "all $NUM_GAMES games done"
