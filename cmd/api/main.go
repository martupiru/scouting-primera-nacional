package main

import (
	"fmt"
	"os"

	"github.com/martupiru/scouting-primera-nacional/internal/apifootball"
	"github.com/martupiru/scouting-primera-nacional/internal/models"
	"github.com/martupiru/scouting-primera-nacional/internal/scoring"
)

const (
	leagueID = 129
	season   = 2024
	topN     = 10
)

func main() {
	apiKey := os.Getenv("API_FOOTBALL_KEY")
	if apiKey == "" {
		fmt.Println("Error: falta la variable de entorno API_FOOTBALL_KEY")
		os.Exit(1)
	}

	client := apifootball.NewClient(apiKey)
	cache, err := apifootball.NewCache("cache")
	if err != nil {
		fmt.Printf("Error iniciando caché: %v\n", err)
		os.Exit(1)
	}

	firstPage, err := client.GetPlayersCached(cache, leagueID, season, 1)
	if err != nil {
		fmt.Printf("Error en página 1: %v\n", err)
		os.Exit(1)
	}

	totalPages := firstPage.Paging.Total
	allEntries := firstPage.Response

	for page := 2; page <= totalPages; page++ {
		resp, err := client.GetPlayersCached(cache, leagueID, season, page)
		if err != nil {
			fmt.Printf("Error en página %d: %v\n", page, err)
			os.Exit(1)
		}
		allEntries = append(allEntries, resp.Response...)
	}

	// Convertir al formato que espera RankByPosition.
	entries := make([]struct {
		Player models.Player
		Stats  []models.PlayerStats
	}, len(allEntries))

	for i, e := range allEntries {
		player, stats := apifootball.ToDomain(e)
		entries[i].Player = player
		entries[i].Stats = stats
	}

	ranking := scoring.RankByPosition(entries, season, topN)

	positions := []string{"Attacker", "Midfielder", "Defender", "Goalkeeper"}
	labels := map[string]string{
		"Attacker":   "DELANTEROS",
		"Midfielder": "MEDIOCAMPISTAS",
		"Defender":   "DEFENSORES",
		"Goalkeeper": "ARQUEROS",
	}

	for _, pos := range positions {
		players, ok := ranking[pos]
		if !ok || len(players) == 0 {
			continue
		}
		fmt.Printf("\n=== %s ===\n", labels[pos])
		fmt.Printf("%-4s %-20s %-6s %-10s %-10s %-10s %-8s\n",
			"#", "Jugador", "Min", "Goles/90", "Asist/90", "Tarj/90", "Score")
		fmt.Println("--------------------------------------------------------------------")
		for i, p := range players {
			fmt.Printf("%-4d %-20s %-6d %-10.2f %-10.2f %-10.2f %-8.1f\n",
				i+1,
				p.Player.Name,
				p.Stats.Minutes,
				p.Stats.GoalsPer90(),
				p.Stats.AssistsPer90(),
				p.Stats.CardsPer90(),
				p.Score,
			)
		}
	}
}
