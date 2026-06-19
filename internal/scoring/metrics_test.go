package scoring

import (
	"testing"

	"github.com/martupiru/scouting-primera-nacional/internal/models"
)

func TestAggregateAndMetrics(t *testing.T) {
	yellow, red, minutes, lineups, appearences := 1, 1, 900, 10, 10

	stats := []models.PlayerStats{
		{Season: 2024, Position: "Goalkeeper"}, // bloque sin datos (otro equipo)
		{
			Season:      2024,
			Position:    "Goalkeeper",
			Minutes:     &minutes,
			Lineups:     &lineups,
			Appearences: &appearences,
			CardsYellow: &yellow,
			CardsRed:    &red,
		},
	}

	agg := Aggregate(stats, 2024)

	if agg.Minutes != 900 {
		t.Errorf("esperaba 900 minutos, llegó %d", agg.Minutes)
	}
	if !agg.IsEligible() {
		t.Errorf("esperaba que fuera elegible (900 >= 450)")
	}
	if agg.StartRatio() != 1.0 {
		t.Errorf("esperaba StartRatio 1.0, llegó %f", agg.StartRatio())
	}

	expectedCardsPer90 := (1.0 + 1.0*2) / 900.0 * 90
	diff := agg.CardsPer90() - expectedCardsPer90
	if diff > 0.0001 || diff < -0.0001 {
		t.Errorf("esperaba CardsPer90 %f, llegó %f", expectedCardsPer90, agg.CardsPer90())
	}
}
