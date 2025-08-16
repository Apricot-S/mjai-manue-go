#!/usr/bin/env bash

HOST=${HOST:-server}
PORT=${PORT:-11600}
ROOM=${ROOM:-default}
NUM_GAMES=${NUM_GAMES:-1}

echo "Waiting 5 seconds for mjai server to start..."
sleep 5

for i in $(seq 1 "$NUM_GAMES"); do
    echo "$i/$NUM_GAMES games"

    ./system.exe mjai_client $PORT ./setup_mjai.json

    echo "finished game $i"
done

echo "all $NUM_GAMES games done"
