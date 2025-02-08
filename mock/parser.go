package mock

import (
	store "github.com/mo-mohamed/txparser/storage"
)

// MockParser is a mock implementation of the Parser interface.
type MockParser struct {
	store store.IStore
}

func NewMockParser(store store.IStore) *MockParser {
	return &MockParser{store: store}
}

func (m *MockParser) GetCurrentBlock() int {
	return m.store.CurrentBlock()
}

func (m *MockParser) Subscribe(address string) bool {
	return m.store.Subscribe(address)
}

func (m *MockParser) GetTransactions(address string) []store.Transaction {
	return m.store.Transactions(address)
}
