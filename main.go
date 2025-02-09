package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mo-mohamed/txparser/api"
	"github.com/mo-mohamed/txparser/blockchain"
	"github.com/mo-mohamed/txparser/parser"
	store "github.com/mo-mohamed/txparser/storage"
)

func main() {
	store := store.NewMemoryStore()
	blockchain := blockchain.NewBlockchain("https://ethereum-rpc.publicnode.com")
	p := parser.NewTxParser(store, blockchain)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	// Start the background polling worker
	go p.StartPolling(ctx)

	server := &http.Server{
		Addr:    ":8080",
		Handler: api.Router(p),
	}

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("Context canceled, shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v\n", err)
	} else {
		log.Println("HTTP server stopped.")
	}

	log.Println("Server stopped.")
}
