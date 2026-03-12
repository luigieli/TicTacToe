package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"tictactoe/internal/api/handlers"
	"tictactoe/internal/repository"
	"tictactoe/internal/service"
)

func main() {
	// Dependency Injection
	repo := repository.NewMemoryRepository()
	svc := service.NewGameService(repo)
	handler := handlers.NewGameHandler(svc)

	// Router setup
	r := chi.NewRouter()

	// Standard Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/games", func(r chi.Router) {
		r.Post("/", handler.CreateGame)
		r.Get("/{id}", handler.GetGameState)
		r.Post("/{id}/move", handler.MakeMove)
	})

	port := ":8080"
	fmt.Printf("Ultimate Tic Tac Toe Backend starting on %s...\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
