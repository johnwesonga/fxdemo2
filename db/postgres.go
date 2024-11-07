package db

import (
	"context"
	"example/fxdemo2/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostGresService struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

func NewPostGresService(db *pgxpool.Pool, logger *zap.Logger) *PostGresService {
	return &PostGresService{logger: logger, db: db}
}
func (s *PostGresService) GetPlayer(id int) (*models.Player, error) {
	var player models.Player
	err := s.db.QueryRow(context.Background(), "SELECT id, name FROM players WHERE id=$1", id).Scan(&player.ID, &player.Name)
	if err != nil {
		return nil, err
	}
	s.logger.Info("Fetched player")
	return &player, nil
}

func (s *PostGresService) GetAllPlayers() ([]models.Player, error) {
	query := "SELECT id, name, score, created_at, updated_at FROM players"
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Player])
}
