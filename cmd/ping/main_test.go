package main

import (
	"errors"
	//"log"
	"os"
	"testing"
	"time"

	//"time"

	"github.com/blinklabs-io/gouroboros/protocol/chainsync"
)

type MockNodeConnection struct {
	DialFunc      func(network, address string) error
	ChainSyncFunc func() *chainsync.ChainSync
	CloseFunc     func() error
}

func (m *MockNodeConnection) Dial(network, address string) error {
	return m.DialFunc(network, address)
}

func (m *MockNodeConnection) ChainSync() *chainsync.ChainSync {
	return m.ChainSyncFunc()
}

func (m *MockNodeConnection) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func restoreDefaultGetTip() {
	getTip = func(sync *chainsync.ChainSync) (*chainsync.Tip, error) {
		return sync.Client.GetCurrentTip()
	}
}

func TestPingNode_Success(t *testing.T) {
	getTip = func(sync *chainsync.ChainSync) (*chainsync.Tip, error) {
		return &chainsync.Tip{}, nil
	}
	defer restoreDefaultGetTip()

	mockConn := &MockNodeConnection{
		DialFunc: func(network, address string) error {
			return nil
		},
		ChainSyncFunc: func() *chainsync.ChainSync {
			return &chainsync.ChainSync{}
		},
	}

	cfg := &Config{SocketPath: "/fake/socket"}
	result := PingNode(mockConn, cfg)

	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}
	if result.ConnectionTime <= 0 {
		t.Errorf("Expected positive connection time")
	}
	if result.ProtocolTime <= 0 {
		t.Errorf("Expected positive protocol time")
	}
}

func TestPingNode_ProtocolError(t *testing.T) {
	getTip = func(sync *chainsync.ChainSync) (*chainsync.Tip, error) {
		return nil, errors.New("fake protocol failure")
	}
	defer restoreDefaultGetTip()

	mockConn := &MockNodeConnection{
		DialFunc: func(network, address string) error {
			return nil
		},
		ChainSyncFunc: func() *chainsync.ChainSync {
			return &chainsync.ChainSync{}
		},
	}

	cfg := &Config{SocketPath: "/fake/socket"}
	result := PingNode(mockConn, cfg)

	if result.Error == nil || result.Error.Error() != "protocol error: fake protocol failure" {
		t.Errorf("Expected protocol error, got: %v", result.Error)
	}
}

func TestPingNode_ConnectionError(t *testing.T) {
	// Override getTip to make sure no protocol call happens - only connection is needed
	getTip = func(sync *chainsync.ChainSync) (*chainsync.Tip, error) {
		return nil, nil
	}
	defer restoreDefaultGetTip()

	// Simulate a connection failure by returning an error in DialFunc
	mockConn := &MockNodeConnection{
		DialFunc: func(network, address string) error {
			return errors.New("connection failed") // Simulating connection error
		},
		ChainSyncFunc: func() *chainsync.ChainSync {
			return &chainsync.ChainSync{}
		},
	}

	cfg := &Config{SocketPath: "/fake/socket"}
	result := PingNode(mockConn, cfg)

	// Check that the correct error is returned
	if result.Error == nil || result.Error.Error() != "connection failed: connection failed" {
		t.Errorf("Expected connection error, got: %v", result.Error)
	}
}

func TestPingNode_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Try to get connection details from environment
	socketPath := os.Getenv("CARDANO_NODE_SOCKET_PATH")
	network := os.Getenv("CARDANO_NODE_NETWORK")
	address := os.Getenv("CARDANO_NODE_ADDRESS")

	// Skip if no connection method is configured
	if socketPath == "" && address == "" && network == "" {
		t.Skip("skipping integration test: set CARDANO_NODE_SOCKET_PATH, CARDANO_NODE_ADDRESS, or CARDANO_NODE_NETWORK to run")
	}

	// Configure based on available options
	cfg := &Config{
		SocketPath: socketPath,
		Address:    address,
		Network:    network,
	}

	// Create real connection
	conn, err := NewConnection(cfg)
	if err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}
	defer conn.Close()

	// Run with timeout
	resultChan := make(chan PingResult)
	go func() {
		resultChan <- PingNode(conn, cfg)
	}()

	select {
	case result := <-resultChan:
		if result.Error != nil {
			t.Fatalf("ping failed: %v", result.Error)
		}
		t.Logf("Integration test successful:")
		t.Logf("Connection time: %v", result.ConnectionTime)
		t.Logf("Protocol time: %v", result.ProtocolTime)
	case <-time.After(30 * time.Second):
		t.Fatal("test timed out after 30 seconds")
	}
}
