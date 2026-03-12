package ports

import (
	"context"
	"tictactoe/internal/domain/models"
)

// GameRepository defines the persistence operations for games.
type GameRepository interface {
	Save(ctx context.Context, game *models.Game) error
	FindByID(ctx context.Context, id string) (*models.Game, error)
}
