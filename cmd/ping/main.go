package main

import (
	"fmt"
	"time"

	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/blinklabs-io/gouroboros/protocol/chainsync"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Magic      uint32 `default:"764824073"`
	SocketPath string `split_words:"true" default:"/ipc/node.socket"`
	Address    string `default:"backbone.cardano.iog.io:3001"`
	Network    string `default:"mainnet"`
}

type NodeConnection interface {
	Dial(network, address string) error
	ChainSync() *chainsync.ChainSync
	Close() error
}

type PingResult struct {
	ConnectionTime time.Duration
	ProtocolTime   time.Duration
	Error          error
}

func GetConnectionDetails(cfg *Config) (string, string) {
	if cfg.SocketPath != "" {
		return "unix", cfg.SocketPath
	}
	return "tcp", cfg.Address
}

func NewConnection(cfg *Config) (NodeConnection, error) {
	return ouroboros.NewConnection(
		ouroboros.WithNetworkMagic(cfg.Magic),
		ouroboros.WithNodeToNode(true),
	)
}

var getTip = func(sync *chainsync.ChainSync) (*chainsync.Tip, error) {
	return sync.Client.GetCurrentTip()
}

func PingNode(conn NodeConnection, cfg *Config) PingResult {
	network, address := GetConnectionDetails(cfg)

	start := time.Now()
	if err := conn.Dial(network, address); err != nil {
		return PingResult{Error: fmt.Errorf("connection failed: %w", err)}
	}
	connTime := time.Since(start)

	start = time.Now()

	if _, err := getTip(conn.ChainSync()); err != nil {
		return PingResult{
			ConnectionTime: connTime,
			Error:          fmt.Errorf("protocol error: %w", err),
		}
	}
	protoTime := time.Since(start)

	return PingResult{
		ConnectionTime: connTime,
		ProtocolTime:   protoTime,
	}
}

func main() {
	var cfg Config
	if err := envconfig.Process("cardano_node", &cfg); err != nil {
		fmt.Printf("Config error: %v\n", err)
		return
	}

	conn, err := NewConnection(&cfg)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	defer conn.Close()

	result := PingNode(conn, &cfg)
	if result.Error != nil {
		fmt.Printf("Ping failed: %v\n", result.Error)
		return
	}

	fmt.Printf("Node-to-Node Ping Results:\n")
	fmt.Printf("Connection established in: %s\n", result.ConnectionTime)
	fmt.Printf("Protocol response time:   %s\n", result.ProtocolTime)
}
