package db

import (
	"example/fxdemo2/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MongoService struct {
	logger *zap.Logger
	client *mongo.Client
}

func NewMongoService(logger *zap.Logger, client *mongo.Client) *MongoService {
	return &MongoService{
		logger: logger,
		client: client,
	}
}
func (s *MongoService) GetPlayer(id int) (*models.Player, error) {
	// MongoDB logic to get player by id
	player := &models.Player{}
	s.logger.Info("Fetched player")
	return player, nil
}
func (s *MongoService) CreatePlayer(player models.Player) error {
	// MongoDB logic to create player
	s.logger.Info("Created player")
	return nil
}
func (s *MongoService) GetAllPlayers() ([]*models.Player, error) {
	var players []*models.Player
	// MongoDB logic to get all players
	s.logger.Info("Fetched all players")
	return players, nil
}
