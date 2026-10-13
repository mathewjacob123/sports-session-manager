package models

type TeamPlayer struct {
	PlayerID    int    `json:"player_id"`
	Name        string `json:"name"`
	Position    string `json:"position"`
	SkillRating int    `json:"skill_rating"`
	PlayerType  string `json:"player_type"`
}

type GeneratedTeams struct {
	TeamA      []TeamPlayer `json:"team_a"`
	TeamB      []TeamPlayer `json:"team_b"`
	TeamASkill int          `json:"team_a_total_skill"`
	TeamBSkill int          `json:"team_b_total_skill"`
	SkillDiff  int          `json:"skill_difference"`
}