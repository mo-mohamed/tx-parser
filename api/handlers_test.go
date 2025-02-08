package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mo-mohamed/txparser/mock"
	store "github.com/mo-mohamed/txparser/storage"
)

func TestCurrentBlockHandler(t *testing.T) {
	mockParser := mock.NewMockParser()
	mockParser.CurrentBlock = 100

	req := httptest.NewRequest("GET", "/currentBlock", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		block := mockParser.GetCurrentBlock()
		json.NewEncoder(w).Encode(map[string]int{"currentBlock": block})
	})

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
	mockParser := mock.NewMockParser()

	req := httptest.NewRequest("GET", "/subscribe?address=0x123456789abcdef", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		if mockParser.Subscribe(address) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Subscribed successfully"))
		} else {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Address already subscribed"))
		}
	})

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if body := w.Body.String(); body != "Subscribed successfully" {
		t.Errorf("Handler returned wrong body: got %v want %v", body, "Subscribed successfully")
	}
}

func TestTransactionsHandler(t *testing.T) {
	mockParser := mock.NewMockParser()
	mockParser.Transactions["0x123456789abcdef"] = []store.Transaction{
		{Hash: "0xabc", From: "0x123456789abcdef", To: "0x987654321", Value: "1000", BlockNumber: "100"},
	}

	req := httptest.NewRequest("GET", "/transactions?address=0x123456789abcdef", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		transactions := mockParser.GetTransactions(address)
		json.NewEncoder(w).Encode(transactions)
	})

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
