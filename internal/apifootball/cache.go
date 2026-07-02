package apifootball

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Cache guarda y recupera respuestas de la API en disco
type Cache struct {
	dir string
}

// NewCache crea un caché que guarda archivos en la carpeta indicada
func NewCache(dir string) (*Cache, error) {
	// MkdirAll crea la carpeta (y las intermedias si hacen falta)
	// El 0755 es el permiso Unix estándar para carpetas
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error creando carpeta de caché: %w", err)
	}
	return &Cache{dir: dir}, nil
}

// filename genera el nombre del archivo para una página dada
func (c *Cache) filename(leagueID, season, page int) string {
	return filepath.Join(c.dir, fmt.Sprintf("players_%d_%d_page%d.json", leagueID, season, page))
}

// Get intenta recuperar una página del caché
// Devuelve nil, nil si no existe todavía (no es un error, simplemente no está)
func (c *Cache) Get(leagueID, season, page int) (*PlayersResponse, error) {
	path := c.filename(leagueID, season, page)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Trayendo página %d desde la API...\n", page)
			return nil, nil // caché miss, va a buscar a la API
		}
		return nil, fmt.Errorf("error leyendo caché: %w", err)
	}

	// Si llegamos acá, era un hit de caché — no imprimimos nada
	var resp PlayersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error parseando caché: %w", err)
	}

	return &resp, nil
}

// Set guarda una respuesta en el cache
func (c *Cache) Set(leagueID, season, page int, resp *PlayersResponse) error {
	path := c.filename(leagueID, season, page)

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("error serializando para caché: %w", err)
	}

	// 0644 = permiso estándar para archivos (lectura/escritura para el dueño)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo caché: %w", err)
	}

	return nil
}
