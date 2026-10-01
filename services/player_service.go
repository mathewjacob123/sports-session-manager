package services

import (
	"errors"
	"sports-session-manager/models"
	"sports-session-manager/repository"
)

type PlayerService struct {
	playerRepo *repository.PlayerRepository
}

func NewPlayerService(playerRepo *repository.PlayerRepository) *PlayerService {
	return &PlayerService{playerRepo: playerRepo}
}

func (s *PlayerService) CreatePlayer(req models.CreatePlayerRequest) (*models.Player, error) {
	existing, err := s.playerRepo.GetPlayerByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("a player with this phone number already exists")
	}
	return s.playerRepo.CreatePlayer(req)
}

func (s *PlayerService) GetPlayers(activeOnly bool) ([]models.Player, error) {
	return s.playerRepo.GetPlayers(activeOnly)
}

func (s *PlayerService) GetPlayerByID(id int) (*models.Player, error) {
	player, err := s.playerRepo.GetPlayerByID(id)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, errors.New("player not found")
	}
	return player, nil
}

func (s *PlayerService) UpdatePlayer(id int, req models.UpdatePlayerRequest) (*models.Player, error) {
	player, err := s.playerRepo.GetPlayerByID(id)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, errors.New("player not found")
	}
	return s.playerRepo.UpdatePlayer(id, req)
}

func (s *PlayerService) BootPlayer(id int) error {
	player, err := s.playerRepo.GetPlayerByID(id)
	if err != nil {
		return err
	}
	if player == nil {
		return errors.New("player not found")
	}
	if !player.IsActive {
		return errors.New("player is already booted")
	}
	return s.playerRepo.BootPlayer(id)
}

func (s *PlayerService) AddStrike(playerID, sessionID int, strikeType string) error {
	player, err := s.playerRepo.GetPlayerByID(playerID)
	if err != nil {
		return err
	}
	if player == nil {
		return errors.New("player not found")
	}
	if !player.IsActive {
		return errors.New("player is already booted")
	}
	return s.playerRepo.AddStrike(playerID, sessionID, strikeType)
}

func (s *PlayerService) GetPlayerStrikes(playerID int) ([]map[string]interface{}, error) {
	player, err := s.playerRepo.GetPlayerByID(playerID)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, errors.New("player not found")
	}
	return s.playerRepo.GetPlayerStrikes(playerID)
}