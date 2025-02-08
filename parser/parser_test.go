package parser_test

import (
	"testing"

	"github.com/mo-mohamed/txparser/mock"
	"github.com/mo-mohamed/txparser/parser"
	store "github.com/mo-mohamed/txparser/storage"
)

func TestGetCurrentBlock(t *testing.T) {
	blockNumberInTest := 123456
	store := store.NewMemoryStore()
	store.SetCurrentBlock(blockNumberInTest)
	blockchain := &mock.BlockchainMock{
		LatestNetworkBlockFunc: func() int { return blockNumberInTest },
	}
	parser := parser.NewTxParser(store, blockchain)

	block := parser.GetCurrentBlock()
	if block != blockNumberInTest {
		t.Errorf("Expected block number to be 123456, got %d", block)
	}
}

func TestSubscribe(t *testing.T) {
	store := store.NewMemoryStore()
	blockchain := &mock.BlockchainMock{
		LatestNetworkBlockFunc: func() int { return 10 },
	}
	parser := parser.NewTxParser(store, blockchain)
	address := "0x123456789abcdef"

	if !parser.Subscribe(address) {
		t.Errorf("Expected subscription to succeed")
	}

	if parser.Subscribe(address) {
		t.Errorf("Expected subscription to fail for already subscribed address")
	}
}

func TestGetTransactions(t *testing.T) {
	store := store.NewMemoryStore()
	blockchain := &mock.BlockchainMock{
		LatestNetworkBlockFunc: func() int { return 10 },
	}
	parser := parser.NewTxParser(store, blockchain)
	address := "0x123456789abcdef"

	transactions := parser.GetTransactions(address)
	if len(transactions) != 0 {
		t.Errorf("Expected no transactions for a new subscription, got %d", len(transactions))
	}
}
