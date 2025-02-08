package parser_test

import (
	"testing"

	"github.com/mo-mohamed/txparser/mock"
	store "github.com/mo-mohamed/txparser/storage"
)

func TestGetCurrentBlock(t *testing.T) {
	store := store.NewMemoryStore()
	store.SetCurrentBlock(123456)
	mockParser := mock.NewMockParser(store)

	block := mockParser.GetCurrentBlock()
	if block != 123456 {
		t.Errorf("Expected block number to be 123456, got %d", block)
	}
}

func TestSubscribe(t *testing.T) {
	store := store.NewMemoryStore()
	mockParser := mock.NewMockParser(store)
	address := "0x123456789abcdef"

	if !mockParser.Subscribe(address) {
		t.Errorf("Expected subscription to succeed")
	}

	if mockParser.Subscribe(address) {
		t.Errorf("Expected subscription to fail for already subscribed address")
	}
}

func TestGetTransactions(t *testing.T) {
	store := store.NewMemoryStore()
	mockParser := mock.NewMockParser(store)
	address := "0x123456789abcdef"

	transactions := mockParser.GetTransactions(address)
	if len(transactions) != 0 {
		t.Errorf("Expected no transactions for a new subscription, got %d", len(transactions))
	}
}
