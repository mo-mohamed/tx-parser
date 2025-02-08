/*
Package api provides HTTP handlers and routing for interacting with the blockchain parser.
This package exposes endpoints that allow clients to subscribe to addresses, retrieve
current block information, and fetch transactions related to subscribed addresses.

Available Endpoints:

- /currentBlock: Retrieves the latest block number parsed by the system.
                 Method: GET
                 Response: { "currentBlock": <block number> }

- /subscribe: Subscribes an address for monitoring inbound or outbound transactions.
              Method: GET
              Query Parameters:
              - address: The Ethereum address to subscribe.
              Response: "Subscribed successfully" or "Address already subscribed"

- /transactions: Fetches transactions associated with a subscribed address.
                 Method: GET
                 Query Parameters:
                 - address: The Ethereum address to fetch transactions for.
                 Response: JSON array of transactions.
*/

package api

import (
	"encoding/json"
	"net/http"

	"github.com/mo-mohamed/txparser/parser"
)

// CurrentBlockHandler handles the /currentBlock endpoint.
func CurrentBlockHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		block := p.GetCurrentBlock()
		json.NewEncoder(w).Encode(map[string]int{"currentBlock": block})
	}
}

// SubscribeHandler handles the /subscribe endpoint.
func SubscribeHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		if p.Subscribe(address) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Subscribed successfully"))
		} else {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Address already subscribed"))
		}
	}
}

// TransactionsHandler handles the /transactions endpoint.
func TransactionsHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		transactions := p.GetTransactions(address)
		json.NewEncoder(w).Encode(transactions)
	}
}

// Router sets up the HTTP routes and returns an http.Handler.
func Router(p parser.Parser) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/current-block", CurrentBlockHandler(p))
	mux.HandleFunc("/subscribe", SubscribeHandler(p))
	mux.HandleFunc("/transactions", TransactionsHandler(p))
	return mux
}
