package entities

type Pokemons struct {
	Count    int
	Next     string
	Previous string
	Results  []PokemonsResultsItem
}

type PokemonsResultsItem struct {
	Name string
	Url  string
}
