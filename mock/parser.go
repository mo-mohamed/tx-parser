package mock

import (
	store "github.com/mo-mohamed/txparser/storage"
)

// MockParser is a mock implementation of the Parser interface.
type MockParser struct {
	CurrentBlock   int
	SubscribedAddr map[string]bool
	Transactions   map[string][]store.Transaction
}

func NewMockParser() *MockParser {
	return &MockParser{
		CurrentBlock:   100,
		SubscribedAddr: make(map[string]bool),
		Transactions:   make(map[string][]store.Transaction),
	}
}

func (m *MockParser) GetCurrentBlock() int {
	return m.CurrentBlock
}

func (m *MockParser) Subscribe(address string) bool {
	if _, exists := m.SubscribedAddr[address]; exists {
		return false
	}
	m.SubscribedAddr[address] = true
	return true
}

func (m *MockParser) GetTransactions(address string) []store.Transaction {
	return m.Transactions[address]
}
