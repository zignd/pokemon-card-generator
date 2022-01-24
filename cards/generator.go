package cards

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/zignd/pokemon-card-generator/requests"
)

func GenerateCardFile(card *Card, destPath string) error {
	content := fmt.Sprintf(`Name: %s
Image: [supressed]
Atk: %d
Def: %d`, card.Name, card.Attack, card.Defense)
	if err := ioutil.WriteFile(destPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write to file at destPath: %w", err)
	}
	return nil
}

func GenerateAllCards(pokemonsCount int, pokemonsCh <-chan requests.GetAllPokemonsItemResult, destPath string) error {
	errorsCh := make(chan error)

	go func() {
		for p := range pokemonsCh {
			filePath := path.Join(destPath, fmt.Sprintf("%s.txt", p.Pokemon.Name))
			if err := GenerateCardFile(&Card{
				Name:    p.Pokemon.Name,
				Picture: nil,
				Attack:  100,
				Defense: 80,
			}, filePath); err != nil {
				errorsCh <- fmt.Errorf("failed to generate the card: %w", err)
				close(errorsCh)
			}
		}
		close(errorsCh)
	}()

	err := <-errorsCh
	if err != nil {
		return err
	}
	return nil
}
