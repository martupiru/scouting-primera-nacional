package apifootball

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://v3.football.api-sports.io"

// Client es el cliente HTTP para API-Football
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient crea un nuevo cliente con la API key dada por el sitio
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// GetPlayers trae una página de jugadores para una liga y temporada dadas
func (c *Client) GetPlayers(leagueID, season, page int) (*PlayersResponse, error) {
	url := fmt.Sprintf("%s/players?league=%d&season=%d&page=%d", baseURL, leagueID, season, page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando el request: %w", err)
	}

	// API-Football requiere la key en este header específico
	req.Header.Set("x-apisports-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando el request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("la API respondió con status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo el body: %w", err)
	}

	var result PlayersResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parseando el JSON: %w", err)
	}

	return &result, nil
}

// GetPlayersCached trae jugadores usando el cache en disco
// Si la página ya fue descargada antes, la lee del disco sin tocar la API
func (c *Client) GetPlayersCached(cache *Cache, leagueID, season, page int) (*PlayersResponse, error) {
	// Primero buscamos en el caché.
	cached, err := cache.Get(leagueID, season, page)
	if err != nil {
		return nil, err
	}
	if cached != nil {
		return cached, nil // caché hit, devolvemos sin llamar a la API
	}

	// Caché miss: vamos a la API y guardamos el resultado.
	resp, err := c.GetPlayers(leagueID, season, page)
	if err != nil {
		return nil, err
	}

	if err := cache.Set(leagueID, season, page, resp); err != nil {
		// Si no se pudo guardar en caché, no es fatal — avisamos pero seguimos.
		fmt.Printf("Advertencia: no se pudo guardar en caché la página %d: %v\n", page, err)
	}

	return resp, nil
}
