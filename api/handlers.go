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

// Router sets up the HTTP routes and returns an http.Handler.
func Router(p parser.Parser) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/currentBlock", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		block := p.GetCurrentBlock()
		json.NewEncoder(w).Encode(map[string]int{"currentBlock": block})
	})

	mux.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
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
	})

	mux.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
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
	})

	return mux
}
