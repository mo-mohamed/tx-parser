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
	"net/http"
	"sync"
	"time"
)

// Transaction represents a blockchain transaction.
type Transaction struct {
	// Hash is the unique identifier for this transaction.
	Hash string `json:"hash"`
	// From is the Ethereum address that initiated the transaction.
	From string `json:"from"`
	// To is the Ethereum address of the account that is the recipient of the transaction.
	To string `json:"to"`
	// Value is the transaction amount.
	Value string `json:"value"`
	// BlockNumber is the number of the transaction.
	BlockNumber string `json:"blockNumber"`
}

// TxParser implements the Parser interface.
type TxParser struct {
	// currentBlock stores the recent block that has been fetched.
	currentBlock int
	// subscribedAddr is a map where keys are Ethereum addresses, and values indicate
	// whether the address is subscribed for transaction monitoring.
	subscribedAddr map[string]bool
	/*
		transactions is a map that holds lists of transactions, indexed by Ethereum address.
		Each key corresponds to an address, and the associated value is a slice of Transaction structs.
	*/
	transactions map[string][]Transaction
	// mu is a mutex that ensures safe concurrent access to the TxParser's state.
	mu sync.Mutex
	// jsonRPCEndpoint is the Ethereum url endpoint used by parser to fetch data
	jsonRPCEndpoint string
}

// NewTxParser initializes a new TxParser.
func NewTxParser(jsonRPCEndpoint string) *TxParser {
	parser := &TxParser{
		subscribedAddr:  make(map[string]bool),
		transactions:    make(map[string][]Transaction),
		jsonRPCEndpoint: jsonRPCEndpoint,
	}
	// Set the current block to the latest block on the network to start polling from this point
	parser.currentBlock = parser.getLatestBlock()
	return parser
}

// GetCurrentBlock fetches the latest block number from the Ethereum network.
func (p *TxParser) GetCurrentBlock() int {
	return p.currentBlock
}

// Subscribe adds an address to the list of subscribers.
func (p *TxParser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.subscribedAddr[address]; exists {
		return false
	}
	p.subscribedAddr[address] = true
	return true
}

// GetTransactions returns a list of transactions for a subscribed address.
func (p *TxParser) GetTransactions(address string) []Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.transactions[address]
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
			latestBlockOnNetwork := p.getLatestBlock()

			// fmt.Printf("CURRENT BLOCK: %d & LATEST ETH BLOCK IS: %d\n", p.GetCurrentBlock(), latestBlockOnNetwork)

			for block := p.GetCurrentBlock() + 1; block <= latestBlockOnNetwork; block++ {
				p.parseNewBlock(block)
			}
			p.currentBlock = latestBlockOnNetwork

			// On Etherium networik, there a new block being added every 12 seconds
			time.Sleep(8 * time.Second)
		}
	}
}

// parseNewBlock helper fetches and extractstransactions from a block.
func (p *TxParser) parseNewBlock(blockNumber int) {
	fmt.Println("Processing Block Number:", blockNumber)
	response, _ := p.jsonRPCRequest("eth_getBlockByNumber", []interface{}{fmt.Sprintf("0x%x", blockNumber), true})

	var blockData struct {
		Result struct {
			Transactions []Transaction `json:"transactions"`
		} `json:"result"`
	}
	json.Unmarshal(response, &blockData)
	p.store(blockData.Result.Transactions)

	// fmt.Printf("CURRENT TRANSACTIONS: %+v\n", p.transactions)
}

// jsonRPCRequest sends a JSON-RPC request to the Ethereum node.
func (p *TxParser) jsonRPCRequest(method string, params []interface{}) ([]byte, error) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	})

	resp, err := http.Post(p.jsonRPCEndpoint, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// getLatestBlock returns thge latest block on the network
func (p *TxParser) getLatestBlock() int {
	p.mu.Lock()
	defer p.mu.Unlock()

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

// store stores transaction in the transactions store
func (p *TxParser) store(transactions []Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, tx := range transactions {
		if p.subscribedAddr[tx.From] || p.subscribedAddr[tx.To] {
			p.transactions[tx.From] = append(p.transactions[tx.From], tx)
			p.transactions[tx.To] = append(p.transactions[tx.To], tx)
		}
	}
}
