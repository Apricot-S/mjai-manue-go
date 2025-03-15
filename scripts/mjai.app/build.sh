#!/usr/bin/env bash

set -euxo pipefail

docker build --pull -f scripts/mjai.app/Dockerfile -t mjai-manue.mjai-app .

container_id=$(docker run -d --rm mjai-manue.mjai-app sleep infinity)

trap 'docker stop $container_id' EXIT

docker cp ${container_id}:/build/mjai-app.zip .
