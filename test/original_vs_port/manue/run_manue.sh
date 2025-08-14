#!/usr/bin/env bash

HOST=${HOST:-server}
PORT=${PORT:-11600}
ROOM=${ROOM:-default}
NUM_GAMES=${NUM_GAMES:-1}

echo "Waiting 15 seconds for mjai server to start..."
sleep 15

cd /app/mjai-manue/coffee

for i in $(seq 1 "$NUM_GAMES"); do
    echo "$i/$NUM_GAMES games"

    coffee main.coffee mjsonp://$HOST:$PORT/$ROOM

    echo "finished game $i"
done

echo "all $NUM_GAMES games done"
