package store

import (
	"sync"
)

type MemoryStore struct {
	// currentBlock stores the recent block that has been fetched.
	currentBlock int
	/*
		subscribedAddr is a map where keys are Ethereum addresses, and values indicate
		whether the address is subscribed for transaction monitoring.
	*/
	subscribedAddr map[string]bool
	/*
		transactions is a map that holds lists of transactions, indexed by Ethereum address.
		Each key corresponds to an address, and the associated value is a slice of Transaction structs.
	*/
	transactions map[string][]Transaction
	// mu is a mutex that ensures safe concurrent access to the TxParser's state.
	mu sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		subscribedAddr: make(map[string]bool),
		transactions:   make(map[string][]Transaction),
	}
}

// Transactions fetches transactions records for a given address
func (m *MemoryStore) Transactions(address string) []Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.transactions[address]
}

// CurrentBlock retrieves the latest processed block
func (m *MemoryStore) CurrentBlock() int {
	return m.currentBlock
}

// SaveTransactions stores transaction in the transactions store
func (m *MemoryStore) SaveTransactions(transactions []Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, tx := range transactions {
		if m.subscribedAddr[tx.From] || m.subscribedAddr[tx.To] {
			m.transactions[tx.From] = append(m.transactions[tx.From], tx)
			m.transactions[tx.To] = append(m.transactions[tx.To], tx)
		}
	}
}

// Subscribe adds an address to the list of subscribers.
func (m *MemoryStore) Subscribe(address string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.subscribedAddr[address]; exists {
		return false
	}
	m.subscribedAddr[address] = true
	return true
}

// SetCurrentBlock stores the latest processed block
func (m *MemoryStore) SetCurrentBlock(blockNumber int) {
	m.currentBlock = blockNumber
}
