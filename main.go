// Copyright 2023 Blink Labs, LLC.
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
	"github.com/blinklabs-io/gouroboros/ledger"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Magic      uint32
	SocketPath string `split_words:"true"`
}

func main() {
	// Set config defaults
	var cfg = Config{
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
	// Get mempool sizes
	capacity, size, numberOfTxs, err := o.LocalTxMonitor().Client.GetSizes()
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf(
		"Mempool size (bytes): %-10d Mempool capacity (bytes): %-10d Transactions: %-10d",
		size,
		capacity,
		numberOfTxs,
	))
	// Get all transactions
	fmt.Println("Transactions:")
	fmt.Println(fmt.Sprintf(" %-20s %s", "Size", "TxHash"))
	for {
		// Get NextTx from Node
		txRawBytes, err := o.LocalTxMonitor().Client.NextTx()
		if err != nil {
			panic(err)
		}
		// Break loop if empty
		if txRawBytes == nil {
			break
		}
		// Get Tx size
		size := len(txRawBytes)
		// Get Tx Hash (of Tx Body)
		tx, err := ledger.NewTransactionFromCbor(ledger.TxTypeBabbage, txRawBytes)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf(" %-20d %s", size, tx.Hash()))
	}
}
