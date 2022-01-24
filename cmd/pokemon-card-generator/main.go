package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zignd/pokemon-card-generator/cards"
	"github.com/zignd/pokemon-card-generator/requests"
)

func main() {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	destPath := "/home/zignd/cards"

	if err := GenerateCards(&httpClient, destPath); err != nil {
		panic(err)
	}
}

func GenerateCards(httpClient *http.Client, destPath string) error {
	pokemonsCount, pokemonsCh, err := requests.GetAllPokemons(httpClient)
	if err != nil {
		return fmt.Errorf("failed to initialize the process to retrieve all pokemons: %w", err)
	}

	if err := cards.GenerateAllCards(pokemonsCount, pokemonsCh, destPath); err != nil {
		return fmt.Errorf("failed to make all the expected cards: %w", err)
	}

	return nil
}
