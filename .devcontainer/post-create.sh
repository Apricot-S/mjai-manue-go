#!/usr/bin/env bash

set -euxo pipefail

# Set up colorful debug output
PS4='+${BASH_SOURCE[0]}:$LINENO: '
if [[ -t 1 ]] && type -t tput >/dev/null; then
  if (( "$(tput colors)" == 256 )); then
    PS4='$(tput setaf 10)'$PS4'$(tput sgr0)'
  else
    PS4='$(tput setaf 2)'$PS4'$(tput sgr0)'
  fi
fi

sudo apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends gnuplot
sudo apt-get clean
sudo rm -rf /var/lib/apt/lists/*

export UV_LINK_MODE=copy
uv venv --clear
uv pip install mjai
