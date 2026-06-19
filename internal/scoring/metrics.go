package scoring

import "github.com/martupiru/scouting-primera-nacional/internal/models"

// MinMinutesThreshold es el mínimo de minutos jugados para que un jugador
// sea considerado en el ranking (evita que muestras chicas distorsionen el resultado
const MinMinutesThreshold = 450

// AggregatedStats junta las estadísticas de un jugador en una temporada
// sumando todos los equipos en los que jugó
type AggregatedStats struct {
	Position     string
	Appearences  int
	Lineups      int
	Minutes      int
	GoalsTotal   int
	GoalsAssists int
	CardsYellow  int
	CardsRed     int
}

// Aggregate suma los bloques de estadísticas de una temporada puntual
func Aggregate(statsList []models.PlayerStats, season int) AggregatedStats {
	var agg AggregatedStats
	for _, s := range statsList {
		if s.Season != season {
			continue
		}
		if agg.Position == "" {
			agg.Position = s.Position
		}
		agg.Appearences += intOrZero(s.Appearences)
		agg.Lineups += intOrZero(s.Lineups)
		agg.Minutes += intOrZero(s.Minutes)
		agg.GoalsTotal += intOrZero(s.GoalsTotal)
		agg.GoalsAssists += intOrZero(s.GoalsAssists)
		agg.CardsYellow += intOrZero(s.CardsYellow)
		agg.CardsRed += intOrZero(s.CardsRed)
	}
	return agg
}

// intOrZero convierte un puntero nullable a int en 0 si es nil.
// Acá ya nos sirve tratarlo como 0: estamos sumando, y "sin dato" en un
// bloque puntual no debería invalidar la suma total del jugador
func intOrZero(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// metricas por 90 minutos!!!

// IsEligible determina si el jugador tiene suficientes minutos para entrar al ranking
func (a AggregatedStats) IsEligible() bool {
	return a.Minutes >= MinMinutesThreshold
}

// GoalsPer90 calcula los goles cada 90 minutos jugados
func (a AggregatedStats) GoalsPer90() float64 {
	if a.Minutes == 0 {
		return 0
	}
	return float64(a.GoalsTotal) / float64(a.Minutes) * 90
}

// AssistsPer90 ídem para asistencias
func (a AggregatedStats) AssistsPer90() float64 {
	if a.Minutes == 0 {
		return 0
	}
	return float64(a.GoalsAssists) / float64(a.Minutes) * 90
}

// CardsPer90 pondera la tarjeta roja el doble que la amarilla
func (a AggregatedStats) CardsPer90() float64 {
	if a.Minutes == 0 {
		return 0
	}
	weighted := float64(a.CardsYellow) + float64(a.CardsRed)*2
	return weighted / float64(a.Minutes) * 90
}

// StartRatio es la proporción de partidos jugados como titular
func (a AggregatedStats) StartRatio() float64 {
	if a.Appearences == 0 {
		return 0
	}
	return float64(a.Lineups) / float64(a.Appearences)
}
