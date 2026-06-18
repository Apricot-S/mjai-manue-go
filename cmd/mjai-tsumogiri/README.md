# mjai-tsumogiri

`mjai-tsumogiri` is a simple agent that always discards the drawn tile. It uses the common command-line modes documented in [`../README.md`](../README.md).

## Usage

Install:

```sh
go install github.com/Apricot-S/mjai-manue-go/cmd/mjai-tsumogiri@latest
```

Run:

```sh
# stdio mode
mjai-tsumogiri [--name <PLAYER_NAME>] [--id <ID>]

# mjsonp TCP client mode
mjai-tsumogiri [--name <PLAYER_NAME>] [--id <ID>] mjsonp://example.com:11600/default
```

The default player name is `"tsumogiri"`.
