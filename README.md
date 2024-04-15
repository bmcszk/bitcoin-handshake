# Bitcoin Protocol Message Handler

This project implements a basic Bitcoin node that connects to other nodes in the network, sends a version message to complete a protocol handshake. It's written in Go and is designed to be a minimal example of interacting within the Bitcoin network.

## Features

- Connects to a specified Bitcoin node.
- Sends a version message upon connection.
- Continuously reads and handles incoming messages based on their type.
- Modular message handling architecture.

## Important Notes

- This project is meant to be a minimal example of using the Bitcoin protocol. It is not intended to be a complete implementation of the Bitcoin protocol.
- It stops working after receiving a `verack` message.

## Prerequisites

- Docker (docker compose)
- Make

## Testing

- Run `make test` to run tests.
- Output should be similar to the following:
```
test-1          | === RUN   TestBitcoinPeerStart
test-1          | Sending message: version
test-1          | Handling message: version
test-1          | Sending message: verack
test-1          | No handler for message type: wtxidrelay
test-1          | No handler for message type: sendaddrv2
test-1          | Handling message: verack
test-1          | --- PASS: TestBitcoinPeerStart (0.00s)
test-1          | PASS
test-1          | ok    github.com/bmcszk/bitcoin-handshake     0.003s
test-1 exited with code 0
```

## License
This project is licensed under the MIT License.