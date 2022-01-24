package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zignd/pokemon-card-generator/entities"
)

func GetPokemons(httpClient *http.Client, limit, offset int) (*entities.Pokemons, error) {
	url := fmt.Sprintf("%s/pokemon?limit=%d&offset=%d", baseURL, limit, offset)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send the request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %w", err)
	}

	var pokemons entities.Pokemons
	if err := json.Unmarshal(b, &pokemons); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %w", err)
	}

	return &pokemons, nil
}

type GetAllPokemonsItemResult struct {
	Pokemon entities.PokemonsResultsItem
	Error   error
}

func GetAllPokemons(httpClient *http.Client) (int, <-chan GetAllPokemonsItemResult, error) {
	resultsCh := make(chan GetAllPokemonsItemResult, 50)

	pokemons, err := GetPokemons(httpClient, 1, 0)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to perform initial request to count available pokemons: %w", err)
	}
	pokemonsCount := pokemons.Count

	go func() {
		const limit = 100
		pages := (pokemonsCount / limit) + 1

		for page := 0; page < pages; page++ {
			offset := page * limit
			fmt.Println("limit", limit, "offset", offset)
			pokemons, err := GetPokemons(httpClient, limit, offset)
			if err != nil {
				resultsCh <- GetAllPokemonsItemResult{
					Error: fmt.Errorf("failed to retrieve the pokemons: %w", err),
				}
			}
			for _, p := range pokemons.Results {
				resultsCh <- GetAllPokemonsItemResult{
					Pokemon: p,
				}
			}
		}
		close(resultsCh)
	}()

	return pokemonsCount, resultsCh, nil
}
