package bitcoin

import (
	"errors"
	"os"
	"testing"
)

func TestBitcoinPeerStart(t *testing.T) {
	network := getEnvOrDefault("NETWORK", "regtest")
	address := getEnvOrDefault("ADDRESS", "localhost:18444")

	peer := NewBitcoinPeer(network, address)

	err := peer.Start()
	if !errors.Is(err, ErrNoMoreMessagesSupported) {
		t.Errorf("unexpected error: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
