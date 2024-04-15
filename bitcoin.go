package bitcoin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type MessageHandler func(*bitcoinMessage) error

var ErrNoMoreMessagesSupported = errors.New("no more messages supported")

var networks = map[string]string{
	"mainnet": "\xf9\xbe\xb4\xd9",
	"testnet": "\x0b\x11\x09\x07",
	"regtest": "\xfa\xbf\xb5\xda",
	"signet":  "\x0a\x03\xcf\x40",
}

type BitcoinPeer struct {
	magic       []byte
	nodeAddress string
	conn        net.Conn
	userAgent   string
	handlers    map[string]MessageHandler
}

func NewBitcoinPeer(networkName, nodeAddress string) *BitcoinPeer {
	network := networks[networkName] // in future validate network

	p := &BitcoinPeer{
		magic:       []byte(network),
		nodeAddress: nodeAddress,
		userAgent:   "/my-client:0.1/",
	}
	p.handlers = map[string]MessageHandler{
		"version": p.handleVersion,
		"verack":  p.handleVerack,
		// Add more handlers for other message types
	}
	return p
}

func (p *BitcoinPeer) Start() error {
	var err error
	p.conn, err = net.Dial("tcp", p.nodeAddress)
	if err != nil {
		return fmt.Errorf("connecting to node: %w", err)
	}
	defer p.conn.Close()

	// Send initial version message
	err = p.sendMessage(p.newVersionMsg())
	if err != nil {
		return fmt.Errorf("sending version message: %w", err)
	}

	// Start reading messages from the connection
	if err := p.readMessages(); err != nil {
		return fmt.Errorf("reading messages: %w", err)
	}

	return nil
}

func (p *BitcoinPeer) sendMessage(message *bitcoinMessage) error {
	fmt.Println("Sending message: ", message.command)
	_, err := p.conn.Write(message.toBytes())
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}
	return nil
}

func (p *BitcoinPeer) readMessages() error {
	for {
		message, err := readBitcoinMessage(p.conn)
		if err != nil {
			return fmt.Errorf("reading message: %w", err)
		}

		// Handle the message based on its type
		if handler, ok := p.handlers[string(message.command)]; ok {
			if err := handler(message); err != nil {
				return fmt.Errorf("handling message: %w", err)
			}
		} else {
			fmt.Printf("No handler for message type: %s\n", message.command)
		}
	}
}

func (p *BitcoinPeer) handleVersion(msg *bitcoinMessage) error {
	fmt.Println("Handling version message")
	verackMsg := newBitcoinMessage(p.magic, "verack", nil)
	return p.sendMessage(verackMsg)
}

func (p *BitcoinPeer) handleVerack(msg *bitcoinMessage) error {
	fmt.Println("Handling verack message")
	// Additional logic to be implemented after receiving verack
	return ErrNoMoreMessagesSupported // let's finish for now
}

// newVersionMsg creates a new version message.
func (p *BitcoinPeer) newVersionMsg() *bitcoinMessage {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int32(70016))      // Protocol version
	binary.Write(buf, binary.LittleEndian, uint64(1))         // Services
	binary.Write(buf, binary.LittleEndian, time.Now().Unix()) // Timestamp

	// Add dummy net addresses for the sender and receiver
	buf.Write(make([]byte, 26)) // from addr: 10 zero bytes + 16 byte IP + 2 byte port
	buf.Write(make([]byte, 26)) // to addr

	binary.Write(buf, binary.LittleEndian, rand.Uint64()) // Nonce
	buf.WriteByte(byte(len(p.userAgent)))                 // UserAgent length
	buf.WriteString(p.userAgent)                          // UserAgent
	binary.Write(buf, binary.LittleEndian, int32(0))      // Start height
	buf.WriteByte(0)                                      // Relay

	return newBitcoinMessage(p.magic, "version", buf.Bytes())
}
