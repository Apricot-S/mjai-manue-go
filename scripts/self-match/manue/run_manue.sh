#!/usr/bin/env bash

set -euo pipefail

HOST=${HOST:-server}
PORT=${PORT:-11600}
ROOM=${ROOM:-default}
NUM_GAMES=${NUM_GAMES:-1}
NAME=${NAME:-ManueGo}
STARTUP_DELAY_SECONDS=${STARTUP_DELAY_SECONDS:-5}

echo "Waiting ${STARTUP_DELAY_SECONDS} seconds for mjai server to start..."
sleep "$STARTUP_DELAY_SECONDS"

for i in $(seq 1 "$NUM_GAMES"); do
    echo "$NAME: $i/$NUM_GAMES games"
    mjai-manue --name "$NAME" "mjsonp://$HOST:$PORT/$ROOM" > /dev/null
    echo "$NAME: finished game $i"
done

echo "$NAME: all $NUM_GAMES games done"
