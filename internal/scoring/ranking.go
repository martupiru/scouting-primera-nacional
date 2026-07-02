package scoring

import (
	"sort"

	"github.com/martupiru/scouting-primera-nacional/internal/models"
)

// RankedPlayer es un jugador con su score final calculado
type RankedPlayer struct {
	Player models.Player
	Stats  AggregatedStats
	Score  float64
}

// weights define los pesos de cada métrica por posición
// Los pesos suman 1.0 en cada posición
type weights struct {
	goalsPer90   float64
	assistsPer90 float64
	startRatio   float64
	cardsPer90   float64 // resta puntos — se invierte al normalizar
	minutesTotal float64
}

var positionWeights = map[string]weights{
	"Attacker": {
		goalsPer90:   0.50,
		assistsPer90: 0.15,
		startRatio:   0.25,
		cardsPer90:   0.10,
		minutesTotal: 0.00,
	},
	"Midfielder": {
		goalsPer90:   0.30,
		assistsPer90: 0.30,
		startRatio:   0.25,
		cardsPer90:   0.15,
		minutesTotal: 0.00,
	},
	"Defender": {
		goalsPer90:   0.15,
		assistsPer90: 0.10,
		startRatio:   0.30,
		cardsPer90:   0.20,
		minutesTotal: 0.25,
	},
	"Goalkeeper": {
		goalsPer90:   0.00,
		assistsPer90: 0.00,
		startRatio:   0.40,
		cardsPer90:   0.10,
		minutesTotal: 0.50,
	},
}

// normalize lleva un valor al rango 0-100 dado un min y max
// Si min == max (todos tienen el mismo valor), devuelve 100 para todos
func normalize(value, min, max float64) float64 {
	if max == min {
		return 100
	}
	return (value - min) / (max - min) * 100
}

// minMax encuentra el min y max de una slice de float64
func minMax(values []float64) (float64, float64) {
	min, max := values[0], values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

// RankByPosition toma jugadores con sus stats, filtra elegibles,
// agrupa por posicion y devuelve los top N de cada posicion rankeados.
func RankByPosition(entries []struct {
	Player models.Player
	Stats  []models.PlayerStats
}, season, topN int) map[string][]RankedPlayer {

	// Agrupar jugadores elegibles por posicion
	groups := make(map[string][]RankedPlayer)
	for _, e := range entries {
		agg := Aggregate(e.Stats, season)
		if !agg.IsEligible() {
			continue
		}
		pos := agg.Position
		if _, ok := positionWeights[pos]; !ok {
			continue // posicion desconocida, la ignoramos
		}
		groups[pos] = append(groups[pos], RankedPlayer{
			Player: e.Player,
			Stats:  agg,
		})
	}

	result := make(map[string][]RankedPlayer)

	for pos, players := range groups {
		if len(players) < 2 {
			players[0].Score = 100.0
			result[pos] = players
			continue
		}

		w := positionWeights[pos]

		// Extraer los valores de cada metrica para poder normalizar
		goals := make([]float64, len(players))
		assists := make([]float64, len(players))
		starts := make([]float64, len(players))
		cards := make([]float64, len(players))
		minutes := make([]float64, len(players))

		for i, p := range players {
			goals[i] = p.Stats.GoalsPer90()
			assists[i] = p.Stats.AssistsPer90()
			starts[i] = p.Stats.StartRatio()
			cards[i] = p.Stats.CardsPer90()
			minutes[i] = float64(p.Stats.Minutes)
		}

		goalsMin, goalsMax := minMax(goals)
		assistsMin, assistsMax := minMax(assists)
		startsMin, startsMax := minMax(starts)
		cardsMin, cardsMax := minMax(cards)
		minutesMin, minutesMax := minMax(minutes)

		// Calcular el score de cada jugador
		for i := range players {
			goalsScore := normalize(goals[i], goalsMin, goalsMax)
			assistsScore := normalize(assists[i], assistsMin, assistsMax)
			startsScore := normalize(starts[i], startsMin, startsMax)
			// Las tarjetas penalizan: invertimos el score (más tarjetas -> peor)
			cardsScore := 100 - normalize(cards[i], cardsMin, cardsMax)
			minutesScore := normalize(minutes[i], minutesMin, minutesMax)

			players[i].Score = goalsScore*w.goalsPer90 +
				assistsScore*w.assistsPer90 +
				startsScore*w.startRatio +
				cardsScore*w.cardsPer90 +
				minutesScore*w.minutesTotal
		}

		// Ordenar de mayor a menor score
		sort.Slice(players, func(i, j int) bool {
			return players[i].Score > players[j].Score
		})

		// Tomar los top N
		if len(players) > topN {
			players = players[:topN]
		}
		result[pos] = players
	}

	return result
}
