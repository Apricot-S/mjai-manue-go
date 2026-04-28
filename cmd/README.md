# Command Line Applications

This directory contains the command line applications for mjai-compatible agents.

## Applications

| Application      | Description                                      | Default Name |
| ---------------- | ------------------------------------------------ | ------------ |
| `mjai-manue`     | AI-powered agent                                 | "Manue030"   |
| `mjai-tsumogiri` | Simple agent that always discards the drawn tile | "tsumogiri"  |

## Installation

```sh
# Install mjai-manue
go install github.com/Apricot-S/mjai-manue-go/cmd/mjai-manue@latest

# Install mjai-tsumogiri
go install github.com/Apricot-S/mjai-manue-go/cmd/mjai-tsumogiri@latest
```

## Usage

```sh
# Basic format
<APP_NAME> [--name <PLAYER_NAME>] [<URL>]

# Pipe mode (standard I/O)
mjai-manue
mjai-tsumogiri --name "SimpleBot"

# TCP/IP client mode
mjai-manue --name "ManueGo" mjsonp://example.com:11600/default
mjai-tsumogiri mjsonp://example.com:11600/room
```

## Modes

### Pipe mode

When `<URL>` is omitted, the application reads JSON Lines from stdin and writes protocol messages to stdout.

Output is sparse. The application writes one line only when it chooses an action. If it writes `{"type":"none"}`, that is an explicit pass for an available action opportunity, not a generic acknowledgement for every input message.

Receiving `end_game` does not terminate the process in pipe mode. The application keeps reading until EOF, so the same process can play multiple games by receiving another `start_game` after `end_game`.

### TCP/IP client mode

When `<URL>` is provided, it must be an `mjsonp://host:port/room` URL.

TCP mode is synchronous with the mjai server. The application sends one response for each input message that expects a response. If the application has no action to take, it sends `{"type":"none"}`.

When TCP mode receives `end_game`, it sends no response, disconnects, and exits.

## I/O rules

- stdout is reserved for protocol output.
- Logs and errors are written to stderr.
- Empty input lines and invalid JSON are treated as runtime errors.
- Output messages are flushed message by message.

## Command-specific notes

- [`mjai-manue`](mjai-manue/) documents `mjai-manue`-specific options and build-time configuration replacement.
- [`mjai-tsumogiri`](mjai-tsumogiri/) documents the simple tsumogiri agent.
