package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mo-mohamed/txparser/api"
	"github.com/mo-mohamed/txparser/parser"
	store "github.com/mo-mohamed/txparser/storage"
)

func main() {
	store := store.NewMemoryStore()
	p := parser.NewTxParser("https://ethereum-rpc.publicnode.com", store)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		fmt.Println("Shutting down...")
		cancel()
	}()

	// Start tyhe background polling worker
	go p.StartPolling(ctx)

	server := &http.Server{
		Addr:    ":8080",
		Handler: api.Router(p),
	}

	go func() {
		fmt.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()
	fmt.Println("Context canceled, shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	} else {
		fmt.Println("HTTP server stopped.")
	}

	fmt.Println("Server stopped.")
}
