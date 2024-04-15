# Bitcoin Protocol Message Handler

This project implements a basic Bitcoin node that connects to other nodes in the network, sends a version message to complete a protocol handshake. It's written in Go and is designed to be a minimal example of interacting within the Bitcoin network.

## Features

- Connects to a specified Bitcoin node.
- Sends a version message upon connection.
- Continuously reads and handles incoming messages based on their type.
- Modular message handling architecture.

## Prerequisites

- Docker
- Make

## Testing

- Run `make test` to run tests.
- Output should be similar to the following:
```
app-1           | === RUN   TestBitcoinPeerStart
app-1           | Handling version message
app-1           | No handler for message type: wtxidrelay
app-1           | No handler for message type: sendaddrv2
app-1           | Verack received
app-1           | --- PASS: TestBitcoinPeerStart (0.00s)
app-1           | PASS
app-1           | ok    github.com/bmcszk/bitcoin-handshake     0.004s
app-1 exited with code 0
```

## License
This project is licensed under the MIT License.