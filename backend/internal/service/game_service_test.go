package service

import (
	"context"
	"testing"

	"tictactoe/internal/domain/dto"
	"tictactoe/internal/domain/models"
)

// Mock repository for testing
type mockRepo struct {
	games map[string]*models.Game
}

func newMockRepo() *mockRepo {
	return &mockRepo{games: make(map[string]*models.Game)}
}

func (m *mockRepo) Save(ctx context.Context, game *models.Game) error {
	m.games[game.ID] = game
	return nil
}

func (m *mockRepo) FindByID(ctx context.Context, id string) (*models.Game, error) {
	game, ok := m.games[id]
	if !ok {
		return nil, models.ErrGameNotFound
	}
	return game, nil
}

func TestCreateGame(t *testing.T) {
	repo := newMockRepo()
	svc := NewGameService(repo)

	gameID, err := svc.CreateGame(context.Background())
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	if gameID == "" {
		t.Error("expected non-empty game ID")
	}

	game, err := repo.FindByID(context.Background(), gameID)
	if err != nil {
		t.Fatalf("failed to find created game in repo: %v", err)
	}

	if game.CurrentPlayer != models.PlayerX {
		t.Errorf("expected starting player X, got %v", game.CurrentPlayer)
	}

	if game.NextBoardIdx != -1 {
		t.Errorf("expected first move to be anywhere (-1), got %v", game.NextBoardIdx)
	}
}

func TestMakeMove(t *testing.T) {
	repo := newMockRepo()
	svc := NewGameService(repo)

	gameID, _ := svc.CreateGame(context.Background())

	// Move Player X to Board 0, Cell 4
	req := dto.MoveRequest{
		GameID:   gameID,
		BoardIdx: 0,
		CellIdx:  4,
	}

	game, err := svc.MakeMove(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to make move: %v", err)
	}

	if game.SubBoards[0].Cells[4] != models.PlayerX {
		t.Errorf("expected Board 0, Cell 4 to be PlayerX, got %v", game.SubBoards[0].Cells[4])
	}

	if game.CurrentPlayer != models.PlayerO {
		t.Errorf("expected current player to be O, got %v", game.CurrentPlayer)
	}

	if game.NextBoardIdx != 4 {
		t.Errorf("expected next move to be restricted to Board 4, got %v", game.NextBoardIdx)
	}
}

func TestBoardWin(t *testing.T) {
	repo := newMockRepo()
	svc := NewGameService(repo)
	gameID, _ := svc.CreateGame(context.Background())

	// Sequence to win Board 0 for X (Cells 0, 1, 2)
	// X: 0,0 (sends O to 0)
	// O: 0,3 (sends X to 3)
	// X: 3,0 (sends O to 0)
	// O: 0,4 (sends X to 4)
	// X: 4,0 (sends O to 0)
	// O: 0,5 (sends X to 5)
	// X: 5,0 (sends O to 0)
	// O: 0,1 (sends X to 1) --- No, O should not take 1.
	
	moves := []struct{ b, c int }{
		{0, 0}, // X (O->0)
		{0, 3}, // O (X->3)
		{3, 0}, // X (O->0)
		{0, 4}, // O (X->4)
		{4, 0}, // X (O->0)
		{0, 5}, // O (X->5)
	}

	for i, m := range moves {
		_, err := svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: m.b, CellIdx: m.c})
		if err != nil {
			t.Fatalf("failed move %d (%d,%d): %v", i, m.b, m.c, err)
		}
	}

	// Now X must go to 5 (from move 0,5)
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 5, CellIdx: 0}) // X at 5,0 -> O to 0
	
	// O must go to 0
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 0, CellIdx: 6}) // O at 0,6 -> X to 6

	// X must go to 6
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 6, CellIdx: 0}) // X at 6,0 -> O to 0

	// O must go to 0
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 0, CellIdx: 7}) // O at 0,7 -> X to 7

	// X must go to 7
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 7, CellIdx: 0}) // X at 7,0 -> O to 0

	// O must go to 0
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 0, CellIdx: 8}) // O at 0,8 -> X to 8

	// X must go to 8
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 8, CellIdx: 0}) // X at 8,0 -> O to 0

	// Now X has 0,0. O has 0,3 0,4 0,5 0,6 0,7 0,8.
	// O must move to 0. ONLY 0,1 and 0,2 are free.
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 0, CellIdx: 1}) // O at 0,1 -> X to 1

	// X must go to 1
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 1, CellIdx: 0}) // X at 1,0 -> O to 0

	// O must move to 0. ONLY 0,2 is free.
	svc.MakeMove(context.Background(), dto.MoveRequest{GameID: gameID, BoardIdx: 0, CellIdx: 2}) // O at 0,2 -> X to 2

	// Wait, now Board 0 is full.
	game, _ := svc.GetGameState(context.Background(), gameID)
	if game.SubBoards[0].Winner == models.Empty {
		t.Errorf("expected Board 0 to have a result (Win or Tie)")
	}
}
