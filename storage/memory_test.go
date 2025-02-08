package store_test

import (
	"testing"

	store "github.com/mo-mohamed/txparser/storage"
)

func TestTransactions(t *testing.T) {
	memoryStore := store.NewMemoryStore()
	transactions := []store.Transaction{
		{Hash: "0xabc", From: "0x123", To: "0x456", Value: "1000", BlockNumber: "1"},
	}

	memoryStore.Subscribe("0x123")
	memoryStore.SaveTransactions(transactions)

	transactions = memoryStore.Transactions("0x123")

	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}
	if transactions[0].Hash != "0xabc" {
		t.Errorf("Expected transaction hash to be '0xabc', got '%s'", transactions[0].Hash)
	}
}

func TestSaveTransactions(t *testing.T) {
	memoryStore := store.NewMemoryStore()
	memoryStore.Subscribe("0x123")
	memoryStore.Subscribe("0x456")
	transactions := []store.Transaction{
		{Hash: "0xabc", From: "0x123", To: "0x456", Value: "1000", BlockNumber: "1"},
	}

	memoryStore.SaveTransactions(transactions)

	savedTransactions := memoryStore.Transactions("0x123")
	if len(savedTransactions) != 1 {
		t.Errorf("Expected 1 transaction for address '0x123', got %d", len(savedTransactions))
	}
	if savedTransactions[0].Hash != "0xabc" {
		t.Errorf("Expected transaction hash to be '0xabc', got '%s'", savedTransactions[0].Hash)
	}
}

func TestSubscribe(t *testing.T) {
	memoryStore := store.NewMemoryStore()

	if !memoryStore.Subscribe("0x123") {
		t.Error("Expected subscription to succeed for new address")
	}

	if memoryStore.Subscribe("0x123") {
		t.Error("Expected subscription to fail for already subscribed address")
	}
}

func TestSetCurrentBlock(t *testing.T) {
	memoryStore := store.NewMemoryStore()

	memoryStore.SetCurrentBlock(200)

	if memoryStore.CurrentBlock() != 200 {
		t.Errorf("Expected current block to be 200, got %d", memoryStore.CurrentBlock())
	}
}
