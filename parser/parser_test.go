package parser_test

import (
	"context"
	"strconv"
	"testing"
	"time"

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

func TestStartPolling(t *testing.T) {
	storage := store.NewMemoryStore()
	mockBlockchain := &mock.BlockchainMock{
		LatestNetworkBlockFunc: func() int { return 100 },
		ParseBlockFunc: func(block int) ([]store.Transaction, error) {
			return []store.Transaction{
				{Hash: "0x" + strconv.Itoa(block), From: "0xabc", To: "0xdef", Value: "500", BlockNumber: strconv.Itoa(block)},
			}, nil
		},
	}
	parser := parser.NewTxParser(storage, mockBlockchain)
	mockBlockchain.LatestNetworkBlockFunc = func() int { return 105 }
	parser.Subscribe("0xabc")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go parser.StartPolling(ctx)

	<-ctx.Done()

	if parser.GetCurrentBlock() != 105 {
		t.Errorf("Expected current block to be updated to 105, got %d", storage.CurrentBlock())
	}

	expectedProcessedBlocks := []int{101, 102, 103, 104, 105}
	for _, block := range expectedProcessedBlocks {
		found := false
		for _, tx := range parser.GetTransactions("0xabc") {
			if tx.BlockNumber == strconv.Itoa(block) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected block %d to be processed, but it was not", block)
		}
	}
}
