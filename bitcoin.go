package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

type BitcoinMessage struct {
	Magic    []byte // Magic value indicating the network
	Command  string // Command name (e.g., "version", "verack")
	Length   uint32 // Length of the payload
	Checksum []byte // Checksum of the payload
	Payload  []byte // Payload of the message
}

type MessageHandler func(*BitcoinMessage) error

var ErrNoMoreMessagesSupported = errors.New("no more messages supported")

type BitcoinPeer struct {
	network     []byte
	nodeAddress string
	conn        net.Conn
	handlers    map[string]MessageHandler
}

var networks = map[string]string{
	"mainnet": "\xf9\xbe\xb4\xd9",
	"testnet": "\x0b\x11\x09\x07",
	"regtest": "\xfa\xbf\xb5\xda",
	"signet":  "\x0a\x03\xcf\x40",
}

func NewBitcoinPeer(network, address string) *BitcoinPeer {
	network = networks[network] // in future validate network

	p := &BitcoinPeer{
		network:     []byte(network),
		nodeAddress: address,
		handlers: map[string]MessageHandler{
			// "version": handleVersion,
			"verack": handleVerack,
			// Add more handlers for other message types
		},
	}
	p.handlers["version"] = p.handleVersion
	return p
}

func (p *BitcoinPeer) Start() error {
	var err error
	p.conn, err = net.Dial("tcp", p.nodeAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	defer p.conn.Close()

	// Send initial version message
	_, err = p.conn.Write(p.newVersionMsg())
	if err != nil {
		return fmt.Errorf("failed to send version message: %w", err)
	}

	// Start reading messages from the connection
	if err := p.readMessages(); err != nil {
		return fmt.Errorf("failed to read messages: %w", err)
	}

	return nil
}

func (p *BitcoinPeer) readMessages() error {
	for {
		// Read the message header first
		header := make([]byte, 24)
		_, err := io.ReadFull(p.conn, header)
		if err != nil {
			return fmt.Errorf("error reading header: %w", err)
		}

		message, err := parseBitcoinMessage(header)
		if err != nil {
			return fmt.Errorf("error parsing message: %w", err)
		}

		// Read the payload based on the length specified in the header
		payload := make([]byte, message.Length)
		_, err = io.ReadFull(p.conn, payload)
		if err != nil {
			return fmt.Errorf("error reading payload: %w", err)
		}
		message.Payload = payload

		// Handle the message based on its type
		if handler, ok := p.handlers[string(message.Command)]; ok {
			if err := handler(message); err != nil {
				return fmt.Errorf("error handling message: %w", err)
			}
		} else {
			fmt.Printf("No handler for message type: %s\n", message.Command)
		}
	}
}

func (p *BitcoinPeer) handleVersion(msg *BitcoinMessage) error {
	fmt.Println("Handling version message")
	return p.sendVerack()
}

func handleVerack(msg *BitcoinMessage) error {
	fmt.Println("Verack received")
	// Additional logic to be implemented after receiving verack
	return ErrNoMoreMessagesSupported // let's finish for now
}

func parseBitcoinMessage(header []byte) (*BitcoinMessage, error) {
	if len(header) != 24 {
		return nil, fmt.Errorf("invalid header length")
	}
	msg := &BitcoinMessage{
		Magic:    header[:4],
		Command:  string(bytes.Trim(header[4:16], "\x00")),
		Length:   binary.LittleEndian.Uint32(header[16:20]),
		Checksum: header[20:24],
	}
	return msg, nil
}

// newVersionMsg creates a new version message.
func (p *BitcoinPeer) newVersionMsg() []byte {
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

	return p.newMessage("version", buf.Bytes())
}

// sendVerack sends a verack message to the peer to acknowledge the reception of version message.
func (p *BitcoinPeer) sendVerack() error {
	verackMsg := p.newMessage("verack", nil) // verack has no payload
	_, err := p.conn.Write(verackMsg)
	return err
}

// newMessage constructs a Bitcoin protocol message with checksum.
func (p *BitcoinPeer) newMessage(command string, payload []byte) []byte {
	message := make([]byte, 24+len(payload)) // 24 byte header + payload

	// Start with the network magic number (mainnet in this case)
	copy(message[:4], p.network)

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
