package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
)

type bitcoinMessage struct {
	magic    []byte // Magic value indicating the network
	command  string // Command name (e.g., "version", "verack")
	length   uint32 // Length of the payload
	checksum []byte // Checksum of the payload
	payload  []byte // Payload of the message
}

func newBitcoinMessage(magic []byte, command string, payload []byte) *bitcoinMessage {
	return &bitcoinMessage{
		magic:    magic,
		command:  command,
		length:   uint32(len(payload)),
		checksum: checksum(payload),
		payload:  payload,
	}
}

func readBitcoinMessage(reader io.Reader) (*bitcoinMessage, error) {
	// Read the message header first
	header := make([]byte, 24)
	_, err := io.ReadFull(reader, header)
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}
	message, err := parseHeader(header)
	if err != nil {
		return nil, fmt.Errorf("parsing header: %w", err)
	}

	// Read the payload based on the length specified in the header
	payload := make([]byte, message.length)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return nil, fmt.Errorf("reading payload: %w", err)
	}
	message.payload = payload

	return message, nil
}

func parseHeader(header []byte) (*bitcoinMessage, error) {
	if len(header) != 24 {
		return nil, fmt.Errorf("invalid header length")
	}
	msg := &bitcoinMessage{
		magic:    header[:4],
		command:  string(bytes.Trim(header[4:16], "\x00")),
		length:   binary.LittleEndian.Uint32(header[16:20]),
		checksum: header[20:24],
	}
	return msg, nil
}

func (m *bitcoinMessage) write(writer io.Writer) error {
	// Start with the network magic number (mainnet in this case)
	if _, err := writer.Write(m.magic); err != nil {
		return fmt.Errorf("writing magic: %w", err)
	}

	// Command
	if _, err := writer.Write([]byte(m.command + string(make([]byte, 12-len(m.command))))); err != nil {
		return fmt.Errorf("writing command: %w", err)
	}

	// Payload length
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, m.length)
	if _, err := writer.Write(length); err != nil {
		return fmt.Errorf("writing length: %w", err)
	}

	// Checksum: first four bytes of double SHA-256
	if _, err := writer.Write(m.checksum[:4]); err != nil { // Only the first 4 bytes of the second hash
		return fmt.Errorf("writing checksum: %w", err)
	}

	// Payload
	if _, err := writer.Write(m.payload); err != nil {
		return fmt.Errorf("writing payload: %w", err)
	}

	return nil
}

func checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:4]
}
