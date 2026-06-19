package apifootball

import "github.com/martupiru/scouting-primera-nacional/internal/models"

// ToDomain convierte un PlayerEntry (formato API-Football) en nuestros
// modelos de dominio: un Player y su lista de PlayerStats.
func ToDomain(entry PlayerEntry) (models.Player, []models.PlayerStats) {
	player := models.Player{
		ID:          entry.PlayerData.ID,
		Name:        entry.PlayerData.Name,
		FirstName:   entry.PlayerData.FirstName,
		LastName:    entry.PlayerData.LastName,
		Age:         entry.PlayerData.Age,
		Nationality: entry.PlayerData.Nationality,
		Photo:       entry.PlayerData.Photo,
	}

	var statsList []models.PlayerStats
	for _, raw := range entry.Statistics {
		statsList = append(statsList, models.PlayerStats{
			TeamID:       raw.Team.ID,
			TeamName:     raw.Team.Name,
			Season:       raw.League.Season,
			Position:     raw.Games.Position,
			Appearences:  raw.Games.Appearences,
			Lineups:      raw.Games.Lineups,
			Minutes:      raw.Games.Minutes,
			GoalsTotal:   raw.Goals.Total,
			GoalsAssists: raw.Goals.Assists,
			CardsYellow:  raw.Cards.Yellow,
			CardsRed:     raw.Cards.Red,
		})
	}

	return player, statsList
}
