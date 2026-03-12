package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"tictactoe/internal/domain/dto"
	"tictactoe/internal/domain/models"
	"tictactoe/internal/ports"
)

type GameHandler struct {
	service ports.GameService
}

// NewGameHandler creates a new instance of GameHandler.
func NewGameHandler(service ports.GameService) *GameHandler {
	return &GameHandler{service: service}
}

// CreateGame handles the POST /game request.
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID, err := h.service.CreateGame(r.Context())
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, dto.CreateGameResponse{GameID: gameID})
}

// GetGameState handles the GET /game/{id} request.
func (h *GameHandler) GetGameState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For simplicity, we'll get the ID from the query param "id"
	// In a real app, we'd use a router like mux or chi for URL params.
	id := r.URL.Query().Get("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "game id is required")
		return
	}

	game, err := h.service.GetGameState(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrGameNotFound) {
			h.respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, game)
}

// MakeMove handles the POST /move request.
func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	game, err := h.service.MakeMove(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, models.ErrInvalidMove) || errors.Is(err, models.ErrWrongBoard) ||
			errors.Is(err, models.ErrBoardAlreadyWon) || errors.Is(err, models.ErrGameOver) ||
			errors.Is(err, models.ErrCellAlreadyTaken) {
			status = http.StatusBadRequest
		} else if errors.Is(err, models.ErrGameNotFound) {
			status = http.StatusNotFound
		}
		h.respondWithError(w, status, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, game)
}

func (h *GameHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *GameHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
