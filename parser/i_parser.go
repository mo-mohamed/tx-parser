package parser

import store "github.com/mo-mohamed/txparser/storage"

// Parser defines an interface for interacting with the blockchain parser.
type Parser interface {
	// GetCurrentBlock retrieves the most recently parsed block number from the blockchain.
	GetCurrentBlock() int

	// Subscribe adds an Ethereum address to the list of monitored addresses for transactions.
	Subscribe(address string) bool

	// GetTransactions retrieves the list of transactions involving a specified address.
	GetTransactions(address string) []store.Transaction
}
