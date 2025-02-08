/*
Package store provides data access to transactions, subscribers and blocks data.
*/
package store

// IStore defines an interface for storing and managing blockchain data.
type IStore interface {
	// CurrentBlock returns the most recently processed block number.
	CurrentBlock() int

	// Transactions retrieves all transactions associated with the specified address.
	Transactions(address string) []Transaction

	// SaveTransactions stores a list of transactions in the store.
	SaveTransactions(transactions []Transaction)

	// SetCurrentBlock updates the current block number in the store.
	SetCurrentBlock(blockNumber int)

	// Subscribe adds an address to the list of monitored addresses for transactions.
	Subscribe(address string) bool
}
