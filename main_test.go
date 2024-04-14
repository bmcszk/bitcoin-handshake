package main

import (
	"errors"
	"testing"
)

func TestStartBitcoinConnection(t *testing.T) {
	err := StartBitcoinConnection()
	if !errors.Is(err, NoMoreMessagesSupported) {
		t.Errorf("unexpected error: %v", err)
	}
}
