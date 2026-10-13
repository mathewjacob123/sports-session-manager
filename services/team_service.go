package services

import (
	"errors"
	"sort"
	"sports-session-manager/models"
	"sports-session-manager/repository"
)

type TeamService struct {
	sessionRepo *repository.SessionRepository
	playerRepo  *repository.PlayerRepository
}

func NewTeamService(sessionRepo *repository.SessionRepository, playerRepo *repository.PlayerRepository) *TeamService {
	return &TeamService{sessionRepo: sessionRepo, playerRepo: playerRepo}
}

func (s *TeamService) GenerateTeams(sessionID int) (*models.GeneratedTeams, error) {

	// get session
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}
	if session.Status != "confirmed" && session.Status != "open" {
		return nil, errors.New("teams can only be generated for open or confirmed sessions")
	}

	// get confirmed slots
	slots, err := s.sessionRepo.GetSessionSlots(sessionID)
	if err != nil {
		return nil, err
	}

	// collect confirmed players only
	var players []models.TeamPlayer
	for _, slot := range slots {
		if slot.Status != "confirmed" || slot.PlayerID == nil {
			continue
		}
		player, err := s.playerRepo.GetPlayerByID(*slot.PlayerID)
		if err != nil || player == nil {
			continue
		}
		players = append(players, models.TeamPlayer{
			PlayerID:    player.ID,
			Name:        player.Name,
			Position:    player.Position,
			SkillRating: player.SkillRating,
			PlayerType:  player.PlayerType,
		})
	}

	if len(players) < 2 {
		return nil, errors.New("not enough confirmed players to generate teams")
	}

	return generateBalancedTeams(players), nil
}

func generateBalancedTeams(players []models.TeamPlayer) *models.GeneratedTeams {
	var goalkeepers []models.TeamPlayer
	var defenders   []models.TeamPlayer
	var midfielders []models.TeamPlayer
	var forwards    []models.TeamPlayer

	// separate by position
	for _, p := range players {
		switch p.Position {
		case "goalkeeper":
			goalkeepers = append(goalkeepers, p)
		case "defender":
			defenders = append(defenders, p)
		case "midfielder":
			midfielders = append(midfielders, p)
		case "forward":
			forwards = append(forwards, p)
		}
	}

	// sort each group by skill descending
	sortBySkill := func(group []models.TeamPlayer) {
		sort.Slice(group, func(i, j int) bool {
			return group[i].SkillRating > group[j].SkillRating
		})
	}
	sortBySkill(goalkeepers)
	sortBySkill(defenders)
	sortBySkill(midfielders)
	sortBySkill(forwards)

	var teamA, teamB []models.TeamPlayer

	// Step 1 — split goalkeepers one each
	// highest rated GK to team with lower skill (start with A)
	for i, gk := range goalkeepers {
		if i%2 == 0 {
			teamA = append(teamA, gk)
		} else {
			teamB = append(teamB, gk)
		}
	}

snakeDraft := func(group []models.TeamPlayer) {
    for i, p := range group {
        if i%2 == 0 {
            teamA = append(teamA, p)
        } else {
            teamB = append(teamB, p)
        }
    }
}

	snakeDraft(defenders) 
	snakeDraft(midfielders)
	snakeDraft(forwards)

	// Step 3 — calculate skill totals
	skillA := totalSkill(teamA)
	skillB := totalSkill(teamB)

	// Step 4 — if difference > 2, try to reduce gap by swapping
	diff := skillA - skillB
	if diff < 0 {
		diff = -diff
	}

	if diff > 2 {
		teamA, teamB = balanceSwap(teamA, teamB, skillA, skillB)
		skillA = totalSkill(teamA)
		skillB = totalSkill(teamB)
		diff = skillA - skillB
		if diff < 0 {
			diff = -diff
		}
	}

	return &models.GeneratedTeams{
		TeamA:      teamA,
		TeamB:      teamB,
		TeamASkill: skillA,
		TeamBSkill: skillB,
		SkillDiff:  diff,
	}
}

func totalSkill(team []models.TeamPlayer) int {
	total := 0
	for _, p := range team {
		total += p.SkillRating
	}
	return total
}

func balanceSwap(teamA, teamB []models.TeamPlayer, skillA, skillB int) ([]models.TeamPlayer, []models.TeamPlayer) {
	// find a pair of same-position players to swap that reduces the skill gap
	bestDiff := abs(skillA - skillB)
	swapA, swapB := -1, -1

	for i, pa := range teamA {
		for j, pb := range teamB {
			// only swap same position
			if pa.Position != pb.Position {
				continue
			}
			// calculate new skill totals after swap
			newSkillA := skillA - pa.SkillRating + pb.SkillRating
			newSkillB := skillB - pb.SkillRating + pa.SkillRating
			newDiff := abs(newSkillA - newSkillB)

			if newDiff < bestDiff {
				bestDiff = newDiff
				swapA = i
				swapB = j
			}
		}
	}

	// perform the best swap found
	if swapA != -1 && swapB != -1 {
		teamA[swapA], teamB[swapB] = teamB[swapB], teamA[swapA]
	}

	return teamA, teamB
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}