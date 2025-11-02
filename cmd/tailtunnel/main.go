package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	handler := api.NewHandler(ts)
	router := api.NewRouter(handler, tailtunnel.FrontendFS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
