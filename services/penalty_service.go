package services

import (
	"errors"
	"sports-session-manager/models"
	"sports-session-manager/repository"
)

type PenaltyService struct {
	penaltyRepo *repository.PenaltyRepository
	playerRepo  *repository.PlayerRepository
	sessionRepo *repository.SessionRepository
}

func NewPenaltyService(
	penaltyRepo *repository.PenaltyRepository,
	playerRepo *repository.PlayerRepository,
	sessionRepo *repository.SessionRepository,
) *PenaltyService {
	return &PenaltyService{
		penaltyRepo: penaltyRepo,
		playerRepo:  playerRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *PenaltyService) GetPenaltiesByPlayer(playerID int) ([]models.PenaltyWithDetails, float64, error) {
	player, err := s.playerRepo.GetPlayerByID(playerID)
	if err != nil {
		return nil, 0, err
	}
	if player == nil {
		return nil, 0, errors.New("player not found")
	}

	penalties, err := s.penaltyRepo.GetPenaltiesByPlayer(playerID)
	if err != nil {
		return nil, 0, err
	}

	// calculate total outstanding
	var totalOwed float64
	for _, p := range penalties {
		if !p.IsPaid {
			totalOwed += p.AmountOwed
		}
	}

	return penalties, totalOwed, nil
}

func (s *PenaltyService) GetPenaltiesBySession(sessionID int) ([]models.PenaltyWithDetails, float64, error) {
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return nil, 0, err
	}
	if session == nil {
		return nil, 0, errors.New("session not found")
	}

	penalties, err := s.penaltyRepo.GetPenaltiesBySession(sessionID)
	if err != nil {
		return nil, 0, err
	}

	var totalOwed float64
	for _, p := range penalties {
		if !p.IsPaid {
			totalOwed += p.AmountOwed
		}
	}

	return penalties, totalOwed, nil
}

func (s *PenaltyService) MarkAsPaid(penaltyID int) error {
	penalty, err := s.penaltyRepo.GetPenaltyByID(penaltyID)
	if err != nil {
		return err
	}
	if penalty == nil {
		return errors.New("penalty not found")
	}
	if penalty.IsPaid {
		return errors.New("penalty is already marked as paid")
	}
	return s.penaltyRepo.MarkAsPaid(penaltyID)
}