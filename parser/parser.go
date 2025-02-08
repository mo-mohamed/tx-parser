/*
Package parser provides tools for parsing and monitoring the Ethereum blockchain.
It allows subscribing to addresses, fetching blocks, and retrieving transactions
through the Ethereum JSON-RPC interface.
*/

package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	store "github.com/mo-mohamed/txparser/storage"
)

// TxParser implements the Parser interface.
type TxParser struct {
	// store is the storage holding transactions, subscribers and blocks data
	store store.IStore

	// jsonRPCEndpoint is the Ethereum url endpoint used by parser to fetch data
	jsonRPCEndpoint string
}

// NewTxParser initializes a new TxParser.
func NewTxParser(jsonRPCEndpoint string, store store.IStore) *TxParser {
	parser := &TxParser{
		jsonRPCEndpoint: jsonRPCEndpoint,
		store:           store,
	}
	// Set the current block to the latest block on the network to start polling from this point
	parser.store.SetCurrentBlock(parser.getLatestNetworkBlock())
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
			latestBlockOnNetwork := p.getLatestNetworkBlock()

			for block := p.store.CurrentBlock() + 1; block <= latestBlockOnNetwork; block++ {
				p.parseNewBlock(block)
			}
			p.store.SetCurrentBlock(latestBlockOnNetwork)

			// On Etherium networik, there a new block being added every 12 seconds
			time.Sleep(5 * time.Second)
		}
	}
}

// parseNewBlock helper fetches and extractstransactions from a block.
func (p *TxParser) parseNewBlock(blockNumber int) {
	fmt.Println("Processing Block Number:", blockNumber)
	response, _ := p.jsonRPCRequest("eth_getBlockByNumber", []interface{}{fmt.Sprintf("0x%x", blockNumber), true})

	var blockData struct {
		Result struct {
			Transactions []store.Transaction `json:"transactions"`
		} `json:"result"`
	}
	json.Unmarshal(response, &blockData)
	p.store.SaveTransactions(blockData.Result.Transactions)

	fmt.Println("Processing Block Completed:", blockNumber)
	// fmt.Printf("CURRENT TRANSACTIONS: %+v\n", p.store.Transactions())
}

// jsonRPCRequest sends a JSON-RPC request to the Ethereum node.
func (p *TxParser) jsonRPCRequest(method string, params []interface{}) ([]byte, error) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      p.randomID(),
	})

	resp, err := http.Post(p.jsonRPCEndpoint, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// getLatestNetworkBlock returns thge latest block on the network
func (p *TxParser) getLatestNetworkBlock() int {
	response, err := p.jsonRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		fmt.Println("Error fetching block number:", err)
		return 0
	}

	var result struct {
		Result string `json:"result"`
	}
	json.Unmarshal(response, &result)
	var latestBlock int
	fmt.Sscanf(result.Result, "0x%x", &latestBlock)
	return latestBlock
}

// randomID generates random Identifier
func (p *TxParser) randomID() string {
	return strconv.Itoa(int(rand.Uint32()))
}
