package blockchain

import store "github.com/mo-mohamed/txparser/storage"

// IBlockchain defines the interface for interacting with a blockchain network.
type IBlockchain interface {
	// ParseBlock fetches and extracts transactions from the specified block number.
	ParseBlock(block int) []store.Transaction

	// LatestNetworkBlock retrieves the number of the latest block available on the blockchain network.
	LatestNetworkBlock() int
}
