package services

import (
	"errors"
	"sports-session-manager/models"
	"sports-session-manager/repository"
	"time"
)

type SessionService struct {
	sessionRepo *repository.SessionRepository
	playerRepo  *repository.PlayerRepository
}

func NewSessionService(sessionRepo *repository.SessionRepository, playerRepo *repository.PlayerRepository) *SessionService {
	return &SessionService{sessionRepo: sessionRepo, playerRepo: playerRepo}
}

func (s *SessionService) CreateSession(req models.CreateSessionRequest) (*models.Session, error) {
	// verify organiser exists
	organiser, err := s.playerRepo.GetPlayerByID(req.CreatedBy)
	if err != nil {
		return nil, err
	}
	if organiser == nil {
		return nil, errors.New("organiser player not found")
	}
	return s.sessionRepo.CreateSession(req)
}

func (s *SessionService) GetSessionByID(id int) (*models.SessionResponse, error) {
	session, err := s.sessionRepo.GetSessionByID(id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	counts, err := s.sessionRepo.GetSlotCounts(id)
	if err != nil {
		return nil, err
	}

	slots, err := s.sessionRepo.GetSessionSlots(id)
	if err != nil {
		return nil, err
	}

	return &models.SessionResponse{
		Session:        *session,
		ConfirmedCount: counts.TotalConfirmed,
		WaitlistCount:  counts.WaitlistCount,
		Slots:          slots,
	}, nil
}

func (s *SessionService) GetSessions() ([]models.Session, error) {
	return s.sessionRepo.GetSessions()
}

func (s *SessionService) CreatePoll(sessionID int, req models.CreatePollRequest) ([]models.DatePoll, error) {
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}
	if session.Status != "polling" && session.Status != "open" {
		return nil, errors.New("cannot create poll for this session")
	}

	// update status to polling
	err = s.sessionRepo.UpdateSessionStatus(sessionID, "polling")
	if err != nil {
		return nil, err
	}

	return s.sessionRepo.CreatePoll(sessionID, req.Options)
}

func (s *SessionService) GetPolls(sessionID int) ([]models.DatePoll, error) {
	return s.sessionRepo.GetPollsBySession(sessionID)
}

func (s *SessionService) VoteForPoll(sessionID, pollID, playerID int) error {
	// check session exists
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}
	if session.Status != "polling" {
		return errors.New("session is not in polling status")
	}

	// check player exists and is active
	player, err := s.playerRepo.GetPlayerByID(playerID)
	if err != nil {
		return err
	}
	if player == nil {
		return errors.New("player not found")
	}
	if !player.IsActive {
		return errors.New("player is not active")
	}

	// check already voted
	hasVoted, err := s.sessionRepo.HasVoted(pollID, playerID)
	if err != nil {
		return err
	}
	if hasVoted {
		return errors.New("already voted for this option")
	}

	return s.sessionRepo.VoteForPoll(pollID, playerID)
}

func (s *SessionService) LockPollWinner(sessionID, pollID int) error {
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}
	if session.Status != "polling" {
		return errors.New("session is not in polling status")
	}
	return s.sessionRepo.LockPollWinner(sessionID, pollID)
}

func (s *SessionService) JoinSession(sessionID, playerID int) (*models.SessionSlot, error) {
	// check session exists and is open
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}
	if session.Status != "open" {
		return nil, errors.New("session is not open for joining")
	}
	if session.TotalSlots == nil {
		return nil, errors.New("session format not decided yet")
	}

	// check player exists and is active
	player, err := s.playerRepo.GetPlayerByID(playerID)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, errors.New("player not found")
	}
	if !player.IsActive {
		return nil, errors.New("player has been booted and cannot join sessions")
	}

	// check not already in session
	alreadyIn, err := s.sessionRepo.PlayerAlreadyIn(sessionID, playerID)
	if err != nil {
		return nil, err
	}
	if alreadyIn {
		return nil, errors.New("player is already in this session")
	}

	// get current slot counts
	counts, err := s.sessionRepo.GetSlotCounts(sessionID)
	if err != nil {
		return nil, err
	}

	slot, err := s.sessionRepo.JoinSession(sessionID, playerID, player.PlayerType, counts, *session.TotalSlots, player.Position)
	if err != nil {
		return nil, err
	}

	return slot, nil
}

func (s *SessionService) DropFromSession(sessionID, playerID int) error {
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}
	if session.Status == "completed" || session.Status == "cancelled" {
		return errors.New("session is already finished")
	}

	// check if match day dropout
	isMatchDay := false
	if session.MatchDate != nil {
		today := time.Now().Truncate(24 * time.Hour)
		matchDay := session.MatchDate.Truncate(24 * time.Hour)
		isMatchDay = today.Equal(matchDay)
	}

	counts, err := s.sessionRepo.GetSlotCounts(sessionID)
	if err != nil {
		return err
	}

	return s.sessionRepo.DropFromSession(
		sessionID, playerID, isMatchDay,
		session.GroundFee, counts.TotalConfirmed,
	)
}

func (s *SessionService) ConfirmSession(sessionID int) error {
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}
	if session.Status != "open" {
		return errors.New("only open sessions can be confirmed")
	}
	return s.sessionRepo.UpdateSessionStatus(sessionID, "confirmed")
}

func (s *SessionService) CompleteSession(sessionID int) error {
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}
	if session.Status != "confirmed" {
		return errors.New("only confirmed sessions can be completed")
	}
	return s.sessionRepo.UpdateSessionStatus(sessionID, "completed")
}

func (s *SessionService) CheckOutsideNotification(sessionID int) (bool, error) {
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return false, err
	}
	if session == nil {
		return false, errors.New("session not found")
	}
	if session.OutsideNotified {
		return false, nil // already notified
	}
	if session.TotalSlots == nil {
		return false, nil
	}

	counts, err := s.sessionRepo.GetSlotCounts(sessionID)
	if err != nil {
		return false, err
	}

	// if less than 70% filled — suggest calling outsiders
	threshold := int(float64(*session.TotalSlots) * 0.7)
	if counts.TotalConfirmed < threshold {
		err = s.sessionRepo.MarkOutsideNotified(sessionID)
		if err != nil {
			return false, err
		}
		return true, nil // yes, notify organiser
	}

	return false, nil
}