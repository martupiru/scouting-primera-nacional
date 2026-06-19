package apifootball

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseAndTransform(t *testing.T) {
	//leo el JSON de prueba (con dos jugadores, uno sin datos y otro con dos equipos)
	data, err := os.ReadFile("testdata/sample_players.json")
	if err != nil {
		t.Fatalf("no pude leer el archivo de prueba: %v", err)
	}

	var resp PlayersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("error al parsear el JSON: %v", err)
	}

	if len(resp.Response) != 2 {
		t.Fatalf("esperaba 2 jugadores, llegaron %d", len(resp.Response))
	}

	// Caso 1: jugador sin datos reales (todo null)
	player1, stats1 := ToDomain(resp.Response[0])
	if player1.Name != "G. Conti" {
		t.Errorf("esperaba nombre 'G. Conti', llegó '%s'", player1.Name)
	}
	if stats1[0].Minutes != nil {
		t.Errorf("esperaba Minutes nil, llegó %v", *stats1[0].Minutes)
	}

	// Caso 2: jugador con dos equipos, el segundo con datos reales
	player2, stats2 := ToDomain(resp.Response[1])
	if player2.Name != "M. Monllor" {
		t.Errorf("esperaba nombre 'M. Monllor', llegó '%s'", player2.Name)
	}
	if len(stats2) != 2 {
		t.Fatalf("esperaba 2 bloques de estadísticas, llegaron %d", len(stats2))
	}
	if stats2[1].Minutes == nil || *stats2[1].Minutes != 900 {
		t.Errorf("esperaba Minutes=900 en el segundo equipo")
	}
	if stats2[1].CardsYellow == nil || *stats2[1].CardsYellow != 1 {
		t.Errorf("esperaba CardsYellow=1 en el segundo equipo")
	}
}
