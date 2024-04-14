package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	// Change the address to the address of the Bitcoin node you want to connect to.
	// nodeAddress = "node_ip:port"
	// nodeAddress = "172.21.176.1:8333"
	// nodeAddress = "localhost:18444"
	nodeAddress = "bitcoin-core:18444"
	// Bitcoin network magic number (mainnet)
	mainNet    = "\xf9\xbe\xb4\xd9"
	regTestNet = "\xfa\xbf\xb5\xda"
)

// main starts the handshake process.
func main() {
	if err := StartBitcoinConnection(); err != nil {
		panic(err)
	}
}

func StartBitcoinConnection() error {

	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}

	// Send initial version message
	_, err = conn.Write(NewVersionMsg())
	if err != nil {
		return fmt.Errorf("failed to send version message: %w", err)
	}

	// Start reading messages from the connection
	if err := readMessages(conn); err != nil {
		return fmt.Errorf("failed to read messages: %w", err)
	}

	return nil
}

// NewVersionMsg creates a new version message.
func NewVersionMsg() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int32(70016))      // Protocol version
	binary.Write(buf, binary.LittleEndian, uint64(1))         // Services
	binary.Write(buf, binary.LittleEndian, time.Now().Unix()) // Timestamp

	// Add dummy net addresses for the sender and receiver
	buf.Write(make([]byte, 26)) // from addr: 10 zero bytes + 16 byte IP + 2 byte port
	buf.Write(make([]byte, 26)) // to addr

	binary.Write(buf, binary.LittleEndian, rand.Uint64()) // Nonce
	buf.WriteByte(byte(len("/my-client:0.1/")))           // UserAgent length
	buf.WriteString("/my-client:0.1/")                    // UserAgent
	binary.Write(buf, binary.LittleEndian, int32(0))      // Start height
	buf.WriteByte(0)                                      // Relay

	return newMessage("version", buf.Bytes())
}

// sendVerack sends a verack message to the peer to acknowledge the reception of version message.
func sendVerack(conn net.Conn) error {
	verackMsg := newMessage("verack", nil) // verack has no payload
	_, err := conn.Write(verackMsg)
	return err
}

// newMessage constructs a Bitcoin protocol message with checksum.
func newMessage(command string, payload []byte) []byte {
	message := make([]byte, 24+len(payload)) // 24 byte header + payload

	// Start with the network magic number (mainnet in this case)
	copy(message[:4], []byte(regTestNet))

	// Command
	copy(message[4:16], []byte(command+string(make([]byte, 12-len(command)))))

	// Payload length
	binary.LittleEndian.PutUint32(message[16:20], uint32(len(payload)))

	// Checksum: first four bytes of double SHA-256
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	copy(message[20:24], secondHash[:4]) // Only the first 4 bytes of the second hash

	// Payload
	copy(message[24:], payload)

	return message
}
