package ports

import (
	"context"
	"tictactoe/internal/domain/dto"
	"tictactoe/internal/domain/models"
)

// GameService defines the business operations for Ultimate Tic Tac Toe.
type GameService interface {
	CreateGame(ctx context.Context) (string, error)
	GetGameState(ctx context.Context, id string) (*models.Game, error)
	MakeMove(ctx context.Context, req dto.MoveRequest) (*models.Game, error)
}
