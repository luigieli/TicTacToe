package repository

import (
	"context"
	"sync"

	"tictactoe/internal/domain/models"
	"tictactoe/internal/ports"
)

type memoryRepo struct {
	mu    sync.RWMutex
	games map[string]*models.Game
}

// NewMemoryRepository creates a new thread-safe in-memory game repository.
func NewMemoryRepository() ports.GameRepository {
	return &memoryRepo{
		games: make(map[string]*models.Game),
	}
}

func (r *memoryRepo) Save(ctx context.Context, game *models.Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Deep copy or pointer? For in-memory, pointer is fine as long as we're careful.
	// But usually, we store a copy to avoid external mutations.
	r.games[game.ID] = game
	return nil
}

func (r *memoryRepo) FindByID(ctx context.Context, id string) (*models.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	game, ok := r.games[id]
	if !ok {
		return nil, models.ErrGameNotFound
	}
	return game, nil
}
