package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"poker-server/internal/network"
)

func main() {
	log.Println("Starting Poker Server (Phase 7 - Lobby System)...")

	// 1. Initialize Server (which owns Hub & LobbyManager)
	// Circular dependency handled inside NewServer if possible, or manual wiring.
	// NewServer creates Hub internally? Or takes it?

	// Let's use manual wiring for clarity.
	// Hub needs a TableManager (Server).
	// Server needs Hub.

	// Step A: Create Server struct (empty Hub)
	srvLogic := network.NewServer(nil)

	// Step B: Create Hub, passing Server as handler
	hub := network.NewHub(srvLogic)

	// Step C: Inject Hub back into Server
	srvLogic.Hub = hub

	// 2. Start Hub
	go hub.Run()

	// 3. HTTP Routes
	// Serve Frontend Static Files
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	// WebSocket Endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		network.ServeWs(hub, w, r)
	})

	// 4. Server Config
	addr := ":8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: nil, // Use DefaultServeMux
	}

	// 5. Start Server in Goroutine
	go func() {
		log.Printf("Listening on %s", addr)
		log.Printf("WS Endpoint: ws://localhost:8080/ws")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
