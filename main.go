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
	"encoding/json"
	"fmt"

	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/blinklabs-io/gouroboros/ledger"
	"github.com/kelseyhightower/envconfig"
)

// We parse environment variables using envconfig into this struct
type Config struct {
	Magic      uint32
	SocketPath string `split_words:"true"`
}

// Create an Asset type using generics since the ledger code does not expose it
type Asset[T ledger.MultiAssetTypeOutput | ledger.MultiAssetTypeMint] struct {
	Name        string `json:"name"`
	NameHex     string `json:"nameHex"`
	PolicyId    string `json:"policyId"`
	Fingerprint string `json:"fingerprint"`
	Amount      T      `json:"amount"`
}

// This code will be executed when run
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
	// Get mempool sizes from Node via LocalTxMonitor Ouroboros mini-protocol
	capacity, size, numberOfTxs, err := o.LocalTxMonitor().Client.GetSizes()
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"Mempool size (bytes): %-10d Mempool capacity (bytes): %-10d Transactions: %-10d\n",
		size,
		capacity,
		numberOfTxs,
	)
	fmt.Println()
	// Get all transactions
	fmt.Println("Transactions:\n")
	// The Ouroboros LocalTxMonitor mini-protocol allows fetching all of the
	// contents of the Node mempool. However, you have to loop and fetch
	// each Tx until the mempool is empty.
	for {
		// Get raw Tx bytes from Node via LocalTxMonitor
		txRawBytes, err := o.LocalTxMonitor().Client.NextTx()
		if err != nil {
			panic(err)
		}
		// Break loop if empty
		if txRawBytes == nil {
			break
		}
		// Get Tx size of raw Tx bytes
		size := len(txRawBytes)
		// Determine transaction type (era) from raw Tx bytes
		txType, err := ledger.DetermineTransactionType(txRawBytes)
		if err != nil {
			panic(err)
		}
		// Get ledger.Transaction from raw Tx bytes
		tx, err := ledger.NewTransactionFromCbor(txType, txRawBytes)
		if err != nil {
			panic(err)
		}
		fmt.Println(" ---")
		// Print Tx size and Tx Hash (of Tx Body)
		fmt.Printf(
			" %-20s %d\n",
			"Size:",
			size,
		)
		fmt.Printf(
			" %-20s %s\n",
			"TxHash:",
			tx.Hash(),
		)
		// Print number of inputs
		fmt.Printf(
			" %-20s %d\n",
			"Inputs:",
			len(tx.Inputs()),
		)
		// Loop through transaction inputs and print ID#Index
		for i, input := range tx.Inputs() {
			fmt.Printf(
				" %-20s %s\n",
				fmt.Sprintf("Input[%d]:", i),
				fmt.Sprintf("%s#%d", input.Id().String(), input.Index()),
			)
		}
		// Print number of outputs
		fmt.Printf(
			" %-20s %d\n",
			"Outputs:",
			len(tx.Outputs()),
		)
		// Loop through transaction outputs
		for o, output := range tx.Outputs() {
			fmt.Printf(
				" %-20s %s\n",
				fmt.Sprintf("Output[%d]:", o),
				fmt.Sprintf("Address: %s", output.Address().String()),
			)
			fmt.Printf(
				" %-20s %s\n",
				fmt.Sprintf("Output[%d]:", o),
				fmt.Sprintf("Amount: %d", output.Amount()),
			)
			if output.Assets() == nil {
				continue
			}
			// Marshal to JSON bytes from ledger.MultiAsset
			j, _ := output.Assets().MarshalJSON()
			var assets []Asset[ledger.MultiAssetTypeOutput]
			// Unmarshal JSON bytes to list of Assets
			err := json.Unmarshal(j, &assets)
			if err != nil {
				panic(err)
			}
			// Loop through each asset and display
			for a, asset := range assets {
				fmt.Printf(
					" %-20s %s\n",
					fmt.Sprintf("Output[%d]:", o),
					fmt.Sprintf(
						"Asset[%d]: Policy: %s, Name: %s, Amount: %d",
						a,
						asset.PolicyId,
						asset.Name,
						asset.Amount,
					),
				)
			}

		}
		fmt.Println()
	}
}
