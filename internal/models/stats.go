package models

// PlayerStats representa las estadísticas de un jugador en una liga/temporada/equipo
// específico. Los campos son punteros porque la API frecuentemente los devuelve
// como null quiere decir "no tenemos ese dato", no "es cero".
type PlayerStats struct {
	TeamID   int    `json:"-"` //no viene directo
	TeamName string `json:"-"`
	Season   int    `json:"-"`
	Position string `json:"-"`

	Appearences *int `json:"-"`
	Lineups     *int `json:"-"`
	Minutes     *int `json:"-"`

	GoalsTotal   *int `json:"-"`
	GoalsAssists *int `json:"-"`

	CardsYellow *int `json:"-"`
	CardsRed    *int `json:"-"`
}
