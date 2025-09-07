# Command Line Application

This directory contains the command line applications.

## Applications

| Application  | Description                         | Default Name |
| ------------ | ----------------------------------- | ------------ |
| `mjai-manue` | AI-powered agent                    | "Manue020"   |
| `tsumogiri`  | Simple rule-based agent for testing | "Tsumogiri"  |

## Installation

```sh
# Install mjai-manue
go install github.com/Apricot-S/mjai-manue-go/cmd/mjai-manue@latest

# Install tsumogiri
go install github.com/Apricot-S/mjai-manue-go/cmd/tsumogiri@latest
```

### To customize configuration files

With the top-level directory of working tree of this repository as the current directory, run the following command:

```sh
go build ./cmd/mjai-manue
```

## Usage

```sh
# Basic format
<APP_NAME> [--name <PLAYER_NAME>] [<URL>]

# Pipe mode (standard I/O)
mjai-manue 2> mjai-manue.log
tsumogiri --name "SimpleBot"

# TCP/IP client mode
mjai-manue --name "ManueGo" mjsonp://example.com:11600/default
tsumogiri mjsonp://example.com:11600/room
```
