package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rajsinghtech/tailtunnel"
	"github.com/rajsinghtech/tailtunnel/internal/api"
	"github.com/rajsinghtech/tailtunnel/internal/tailscale"
)

func main() {
	log.Println("Starting TailTunnel...")

	ts, err := tailscale.NewTailscaleClient()
	if err != nil {
		log.Fatalf("Failed to create Tailscale client: %v", err)
	}
	defer ts.Close()

	log.Println("Tailscale client initialized")

	// Monitor for auth URLs and automatically open browser
	go func() {
		for authURL := range ts.AuthURL() {
			log.Printf("\n\n=================================")
			log.Printf("Tailscale Login Required!")
			log.Printf("Opening browser to: %s", authURL)
			log.Printf("=================================\n\n")
		}
	}()

	handler := api.NewHandler(ts)
	router := api.NewRouter(handler, tailtunnel.FrontendFS)

	// Serve on tailnet with HTTPS (if available)
	go func() {
		log.Println("Starting server on Tailscale network...")
		if err := ts.ListenHTTPS(router); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	log.Println("Server stopped")
}
