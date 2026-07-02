package main

import (
	"fmt"
	"os"

	"github.com/martupiru/scouting-primera-nacional/internal/apifootball"
	"github.com/martupiru/scouting-primera-nacional/internal/scoring"
)

const (
	leagueID = 129
	season   = 2024
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

	fmt.Println("Trayendo jugadores de Primera Nacional 2024...")

	firstPage, err := client.GetPlayersCached(cache, leagueID, season, 1)
	if err != nil {
		fmt.Printf("Error en página 1: %v\n", err)
		os.Exit(1)
	}

	totalPages := firstPage.Paging.Total
	fmt.Printf("Total de páginas: %d\n\n", totalPages)

	allEntries := firstPage.Response

	for page := 2; page <= totalPages; page++ {
		resp, err := client.GetPlayersCached(cache, leagueID, season, page)
		if err != nil {
			fmt.Printf("Error en página %d: %v\n", page, err)
			os.Exit(1)
		}
		allEntries = append(allEntries, resp.Response...)
	}

	fmt.Printf("\nTotal jugadores recibidos: %d\n", len(allEntries))
	fmt.Println("Jugadores elegibles (≥90 min):\n")

	for _, entry := range allEntries {
		player, statsList := apifootball.ToDomain(entry)
		agg := scoring.Aggregate(statsList, season)

		if !agg.IsEligible() {
			continue
		}

		fmt.Printf("%-20s | %-12s | %4d min | %.2f goles/90 | %.2f asist/90 | %.2f tarj/90\n",
			player.Name,
			agg.Position,
			agg.Minutes,
			agg.GoalsPer90(),
			agg.AssistsPer90(),
			agg.CardsPer90(),
		)
	}
}
