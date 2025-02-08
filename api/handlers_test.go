package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mo-mohamed/txparser/api"
	"github.com/mo-mohamed/txparser/mock"
	store "github.com/mo-mohamed/txparser/storage"
)

func TestCurrentBlockHandler(t *testing.T) {
	store := store.NewMemoryStore()
	store.SetCurrentBlock(100)
	mockParser := mock.NewMockParser(store)

	req := httptest.NewRequest("GET", "/currentBlock", nil)
	w := httptest.NewRecorder()

	handler := api.CurrentBlockHandler(mockParser)
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
	mockParser := mock.NewMockParser(store)

	req := httptest.NewRequest("GET", "/subscribe?address=0x123456789abcdef", nil)
	w := httptest.NewRecorder()

	handler := api.SubscribeHandler(mockParser)
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
	mockParser := mock.NewMockParser(storage)
	mockParser.Subscribe("0x123456789abcdef")
	storage.SaveTransactions(mockedTrans)

	req := httptest.NewRequest("GET", "/transactions?address=0x123456789abcdef", nil)
	w := httptest.NewRecorder()

	handler := api.TransactionsHandler(mockParser)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var transactions []store.Transaction
	err := json.NewDecoder(w.Body).Decode(&transactions)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	fmt.Println("LENGTH is:", len(transactions))

	if len(transactions) != 1 || transactions[0].Hash != "0xabc" {
		t.Errorf("Handler returned wrong transactions: got %v", transactions)
	}
}
