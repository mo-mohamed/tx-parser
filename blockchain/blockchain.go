package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	store "github.com/mo-mohamed/txparser/storage"
)

type Blockchain struct {
	//jsonRPCEndpoint is the endpoint for blockchain network
	jsonRPCEndpoint string
}

type blockData struct {
	Result struct {
		Transactions []store.Transaction `json:"transactions"`
	} `json:"result"`
}

// NewBlockchain returns new instance of the blockchain client
func NewBlockchain(endpoint string) *Blockchain {
	return &Blockchain{
		jsonRPCEndpoint: endpoint,
	}
}

// ParseBlock returns the transactions within a block
func (b *Blockchain) ParseBlock(block int) ([]store.Transaction, error) {
	var blockData blockData
	response, err := b.jsonRPCRequest("eth_getBlockByNumber", []interface{}{fmt.Sprintf("0x%x", block), true})
	if err != nil {
		log.Println("Error fetching block number:", err)
		return nil, fmt.Errorf("error fetching block number: %s", err.Error())
	}
	json.Unmarshal(response, &blockData)
	return blockData.Result.Transactions, nil
}

// LatestNetworkBlock returns the latest block on the network
func (b *Blockchain) LatestNetworkBlock() int {
	response, err := b.jsonRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		log.Println("Error fetching block number:", err)
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

// jsonRPCRequest issues a RPC request to the Etherium blockchain network.
func (b *Blockchain) jsonRPCRequest(method string, params []interface{}) ([]byte, error) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      b.randomID(),
	})

	resp, err := http.Post(b.jsonRPCEndpoint, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// randomID generates random Identifier
func (b *Blockchain) randomID() string {
	return strconv.Itoa(int(rand.Uint32()))
}
