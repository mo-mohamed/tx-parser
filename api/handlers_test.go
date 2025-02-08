package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mo-mohamed/txparser/api"
	"github.com/mo-mohamed/txparser/mock"
	"github.com/mo-mohamed/txparser/parser"
	store "github.com/mo-mohamed/txparser/storage"
)

func TestCurrentBlockHandler(t *testing.T) {
	store := store.NewMemoryStore()
	store.SetCurrentBlock(100)
	blockchain := &mock.BlockchainMock{
		LatestNetworkBlockFunc: func() int { return 100 },
	}
	parser := parser.NewTxParser(store, blockchain)

	req := httptest.NewRequest("GET", "/currentBlock", nil)
	w := httptest.NewRecorder()

	handler := api.CurrentBlockHandler(parser)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]int
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if response["currentBlock"] != 100 {
		t.Errorf("Handler returned wrong block number: got %v want %v", response["currentBlock"], 100)
	}
}
func TestSubscribeHandler(t *testing.T) {
	store := store.NewMemoryStore()
	blockchain := &mock.BlockchainMock{
		LatestNetworkBlockFunc: func() int { return 10 },
	}
	parser := parser.NewTxParser(store, blockchain)

	req := httptest.NewRequest("GET", "/subscribe?address=0x123456789abcdef", nil)
	w := httptest.NewRecorder()

	handler := api.SubscribeHandler(parser)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if body := w.Body.String(); body != "Subscribed successfully" {
		t.Errorf("Handler returned wrong body: got %v want %v", body, "Subscribed successfully")
	}
}

func TestTransactionsHandler(t *testing.T) {
	mockedTrans := []store.Transaction{
		{Hash: "0xabc", From: "0x123456789abcdef", To: "0x987654321", Value: "1000", BlockNumber: "100"},
	}
	storage := store.NewMemoryStore()
	blockchain := &mock.BlockchainMock{
		LatestNetworkBlockFunc: func() int { return 10 },
	}
	parser := parser.NewTxParser(storage, blockchain)

	parser.Subscribe("0x123456789abcdef")
	storage.SaveTransactions(mockedTrans)

	req := httptest.NewRequest("GET", "/transactions?address=0x123456789abcdef", nil)
	w := httptest.NewRecorder()

	handler := api.TransactionsHandler(parser)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var transactions []store.Transaction
	err := json.NewDecoder(w.Body).Decode(&transactions)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if len(transactions) != 1 || transactions[0].Hash != "0xabc" {
		t.Errorf("Handler returned wrong transactions: got %v", transactions)
	}
}
