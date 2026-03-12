package main

import (
	"fmt"
	"net/http"

	"tictactoe/internal/api/handlers"
	"tictactoe/internal/repository"
	"tictactoe/internal/service"
)

func main() {
	// Dependency Injection
	repo := repository.NewMemoryRepository()
	svc := service.NewGameService(repo)
	handler := handlers.NewGameHandler(svc)

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateGame(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetGameState(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/move", handler.MakeMove)

	// CORS Middleware
	corsMux := corsMiddleware(mux)

	port := ":8080"
	fmt.Printf("Ultimate Tic Tac Toe Backend starting on %s...\n", port)
	if err := http.ListenAndServe(port, corsMux); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
