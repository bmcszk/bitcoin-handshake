package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type BitcoinMessage struct {
	Magic    []byte // Magic value indicating the main network
	Command  string // Command name (e.g., "version", "verack")
	Length   uint32 // Length of the payload
	Checksum []byte // Checksum of the payload
	Payload  []byte // Payload of the message
}

type MessageHandler func(*BitcoinMessage, net.Conn) error

var NoMoreMessagesSupported = errors.New("no more messages supported")

var handlers = map[string]MessageHandler{
	"version": handleVersion,
	"verack":  handleVerack,
	// Add more handlers for other message types
}

func handleVersion(msg *BitcoinMessage, conn net.Conn) error {
	fmt.Println("Handling version message")
	return sendVerack(conn)
}

func handleVerack(msg *BitcoinMessage, conn net.Conn) error {
	fmt.Println("Verack received")
	// Additional logic to be implemented after receiving verack
	return NoMoreMessagesSupported // let's finish for now
}

func readMessages(conn net.Conn) error {
	defer conn.Close()

	for {
		// Read the message header first
		header := make([]byte, 24)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("error reading: %w", err)
			}
			break
		}

		message, err := parseBitcoinMessage(header)
		if err != nil {
			fmt.Println("Error parsing message:", err)
			continue
		}

		// Read the payload based on the length specified in the header
		payload := make([]byte, message.Length)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			return fmt.Errorf("error reading payload: %w", err)
		}
		message.Payload = payload

		// Handle the message based on its type
		if handler, ok := handlers[string(message.Command)]; ok {
			if err := handler(message, conn); err != nil {
				return fmt.Errorf("error handling message: %w", err)
			}
		} else {
			fmt.Printf("No handler for message type: %s\n", message.Command)
		}
	}

	return nil
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
