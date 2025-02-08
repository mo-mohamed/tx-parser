/*
Package parser provides tools for parsing and monitoring the Ethereum blockchain.
It allows subscribing to addresses, fetching blocks, and retrieving transactions
through the Ethereum JSON-RPC interface.
*/

package parser

import (
	"context"
	"fmt"
	"time"

	"github.com/mo-mohamed/txparser/blockchain"
	store "github.com/mo-mohamed/txparser/storage"
)

// TxParser implements the Parser interface.
type TxParser struct {
	// store is the storage holding transactions, subscribers and blocks data
	store store.IStore

	// blockChain is a blockchain client for communicating with the blockchain network
	blockChain blockchain.IBlockchain
}

// NewTxParser initializes a new TxParser.
func NewTxParser(store store.IStore, blockchain blockchain.IBlockchain) *TxParser {
	parser := &TxParser{
		store:      store,
		blockChain: blockchain,
	}
	// Set the current block to the latest block on the network to start polling from this point
	parser.store.SetCurrentBlock(parser.blockChain.LatestNetworkBlock())
	return parser
}

// GetCurrentBlock fetches the latest block number from the Ethereum network.
func (p *TxParser) GetCurrentBlock() int {
	return p.store.CurrentBlock()
}

// Subscribe adds an address to the list of subscribers.
func (p *TxParser) Subscribe(address string) bool {
	return p.store.Subscribe(address)
}

// GetTransactions returns a list of transactions for a subscribed address.
func (p *TxParser) GetTransactions(address string) []store.Transaction {
	return p.store.Transactions(address)
}

// StartPolling starts fetching new blocks, alternatively "eth_subscribe" can be used for new blocks.
func (p *TxParser) StartPolling(ctx context.Context) {
	fmt.Println("Starting Polling Blocks")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Polling blocks stopped.")
			return
		default:
			latestBlockOnNetwork := p.blockChain.LatestNetworkBlock()

			for block := p.store.CurrentBlock() + 1; block <= latestBlockOnNetwork; block++ {
				p.processBlock(block)
			}
			p.store.SetCurrentBlock(latestBlockOnNetwork)

			// On Etherium network, there is a new block added every 12 seconds
			time.Sleep(5 * time.Second)
		}
	}
}

// processBlock helper fetches and extractstransactions from a block.
func (p *TxParser) processBlock(blockNumber int) {
	fmt.Println("Processing Block Number:", blockNumber)

	transactions := p.blockChain.ParseBlock(blockNumber)
	p.store.SaveTransactions(transactions)

	fmt.Println("Processing Block Completed:", blockNumber)
}
