// Copyright 2024 Blink Labs Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/kelseyhightower/envconfig"
)

// We parse environment variables using envconfig into this struct
type Config struct {
	Magic      uint32
	SocketPath string `split_words:"true"`
}

// This code will be executed when run
func main() {
	// Set config defaults
	cfg := Config{
		Magic:      764824073,
		SocketPath: "/ipc/node.socket",
	}
	// Parse environment variables
	if err := envconfig.Process("cardano_node", &cfg); err != nil {
		panic(err)
	}
	// Create error channel
	errorChan := make(chan error)
	// start error handler
	go func() {
		for {
			err := <-errorChan
			panic(err)
		}
	}()
	// Configure Ouroboros
	o, err := ouroboros.NewConnection(
		ouroboros.WithNetworkMagic(uint32(cfg.Magic)),
		ouroboros.WithErrorChan(errorChan),
		ouroboros.WithNodeToNode(false),
	)
	if err != nil {
		panic(err)
	}
	// Connect to Node socket
	if err = o.Dial("unix", cfg.SocketPath); err != nil {
		panic(err)
	}
	// Get current tip from Node via NtC ChainSync Ouroboros mini-protocol
	tip, err := o.ChainSync().Client.GetCurrentTip()
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"Chain Tip:\nSlot: %-10d Block Hash: %x\n",
		tip.Point.Slot,
		tip.Point.Hash,
	)
	fmt.Println()
}
