package models

// Player representa los datos básicos de un jugador (no cambian por temporada).
type Player struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Photo       string `json:"photo"`
}
