package service

import (
	"context"

	"github.com/google/uuid"
	"tictactoe/internal/domain/dto"
	"tictactoe/internal/domain/models"
	"tictactoe/internal/ports"
)

type gameService struct {
	repo ports.GameRepository
}

// NewGameService creates a new instance of the game service.
func NewGameService(repo ports.GameRepository) ports.GameService {
	return &gameService{repo: repo}
}

// CreateGame initializes a new Ultimate Tic Tac Toe session.
func (s *gameService) CreateGame(ctx context.Context) (string, error) {
	newGame := &models.Game{
		ID:            uuid.New().String(),
		CurrentPlayer: models.PlayerX,
		NextBoardIdx:  -1,
		IsGameOver:    false,
	}

	// Initialize empty boards
	for i := 0; i < 9; i++ {
		newGame.SubBoards[i] = models.Board{
			Cells:  [9]models.CellState{},
			Winner: models.Empty,
		}
	}

	if err := s.repo.Save(ctx, newGame); err != nil {
		return "", err
	}

	return newGame.ID, nil
}

func (s *gameService) GetGameState(ctx context.Context, id string) (*models.Game, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *gameService) MakeMove(ctx context.Context, req dto.MoveRequest) (*models.Game, error) {
	game, err := s.repo.FindByID(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// Bouncers: Validate move preconditions
	if game.IsGameOver {
		return nil, models.ErrGameOver
	}

	if game.NextBoardIdx != -1 && game.NextBoardIdx != req.BoardIdx {
		return nil, models.ErrWrongBoard
	}

	if req.BoardIdx < 0 || req.BoardIdx > 8 || req.CellIdx < 0 || req.CellIdx > 8 {
		return nil, models.ErrInvalidMove
	}

	targetBoard := &game.SubBoards[req.BoardIdx]
	if targetBoard.Winner != models.Empty {
		return nil, models.ErrBoardAlreadyWon
	}

	if targetBoard.Cells[req.CellIdx] != models.Empty {
		return nil, models.ErrCellAlreadyTaken
	}

	// Execute Move
	targetBoard.Cells[req.CellIdx] = game.CurrentPlayer

	// Check Sub-Board Winner
	if winner := s.calculateWinner(targetBoard.Cells[:]); winner != models.Empty {
		targetBoard.Winner = winner
	} else if s.isBoardFull(targetBoard.Cells[:]) {
		targetBoard.Winner = models.Tie
	}

	// Check Global Winner
	var globalCells [9]models.CellState
	for i := 0; i < 9; i++ {
		globalCells[i] = game.SubBoards[i].Winner
	}
	if globalWinner := s.calculateWinner(globalCells[:]); globalWinner != models.Empty {
		game.Winner = globalWinner
		game.IsGameOver = true
	} else if s.isBoardFull(globalCells[:]) {
		game.Winner = models.Tie
		game.IsGameOver = true
	}

	// Update Next Move Constraints
	game.NextBoardIdx = req.CellIdx
	if game.SubBoards[game.NextBoardIdx].Winner != models.Empty {
		game.NextBoardIdx = -1 // Target board is already won, player can go anywhere
	}

	// Switch Player
	if game.CurrentPlayer == models.PlayerX {
		game.CurrentPlayer = models.PlayerO
	} else {
		game.CurrentPlayer = models.PlayerX
	}

	if err := s.repo.Save(ctx, game); err != nil {
		return nil, err
	}

	return game, nil
}

func (s *gameService) calculateWinner(cells []models.CellState) models.CellState {
	winPatterns := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6},             // Diagonals
	}

	for _, p := range winPatterns {
		if cells[p[0]] != models.Empty && cells[p[0]] != models.Tie &&
			cells[p[0]] == cells[p[1]] && cells[p[0]] == cells[p[2]] {
			return cells[p[0]]
		}
	}
	return models.Empty
}

func (s *gameService) isBoardFull(cells []models.CellState) bool {
	for _, c := range cells {
		if c == models.Empty {
			return false
		}
	}
	return true
}

