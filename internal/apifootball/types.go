package apifootball

// PlayersResponse representa la respuesta completa del endpoint /players.
type PlayersResponse struct {
	Response []PlayerEntry `json:"response"`
	Paging   struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
}

// PlayerEntry agrupa los datos básicos de un jugador con sus estadísticas.
type PlayerEntry struct {
	PlayerData RawPlayer       `json:"player"`
	Statistics []RawStatistics `json:"statistics"`
}

// RawPlayer son los datos personales del jugador, tal como vienen en el JSON.
type RawPlayer struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Photo       string `json:"photo"`
}

// RawStatistics son las estadísticas de un jugador para un equipo/temporada.
// Solo modelamos los campos que vimos poblados en la respuesta real
// (games, goals, cards). El resto (shots, passes, tackles, duels, dribbles,
// fouls, penalty) viene null para esta liga, así que no los necesitamos
// todavía -- si en el futuro hace falta, los agregamos.
type RawStatistics struct {
	Team struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"team"`

	League struct {
		Season int `json:"season"`
	} `json:"league"`

	Games struct {
		Appearences *int   `json:"appearences"`
		Lineups     *int   `json:"lineups"`
		Minutes     *int   `json:"minutes"`
		Position    string `json:"position"`
	} `json:"games"`

	Goals struct {
		Total   *int `json:"total"`
		Assists *int `json:"assists"`
	} `json:"goals"`

	Cards struct {
		Yellow *int `json:"yellow"`
		Red    *int `json:"red"`
	} `json:"cards"`
}
