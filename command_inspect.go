package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)


type Stat struct {
	Stat struct {
		Name string `json:"name"`
	} `json:"stat"`
	BaseStat int `json:"base_stat"`
}

type PokemonDetails struct {
	Name     string `json:"name"`
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Stats    []Stat `json:"stats"`
	Types    []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}



 
func GetCaughtPokemons() map[string]Pokemon {
	return caughtPokemons
}


func commandInspect(cfg *config, args []string, secondArg interface{}) error {
	caughtPokemons := GetCaughtPokemons()
	if len(caughtPokemons) == 0 {
		fmt.Println("No Pokémon caught yet.")
		return nil
	}
	if len(args) == 0 {
		fmt.Println("Specify a Pokémon to inspect:")
		for pokemonName := range caughtPokemons {
			fmt.Printf("- %s\n", pokemonName)
		}
		return nil
		
		
	}
	
	pokemonName := args[0]
	pokemon, exists := caughtPokemons[pokemonName]
	if !exists {
		fmt.Printf("Pokémon not caught: %s\n", pokemonName)
	} else {

		details, err := fetchPokemonDetails(pokemonName)
		if err != nil {
			return err
		}

		
		fmt.Printf("Pokémon caught: %s\n", pokemon.Name)
		fmt.Printf("Name: %s\n", details.Name)
		fmt.Printf("Height: %d\n", details.Height)
		fmt.Printf("Weight: %d\n", details.Weight)

		
		fmt.Println("Stats:")
		for _, stat := range details.Stats {
			fmt.Printf("  -%s: %d\n", strings.ReplaceAll(stat.Stat.Name, "-", " "), stat.BaseStat)
		}

		
		fmt.Println("Types:")
		for _, t := range details.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
	}

	return nil
}


func fetchPokemonDetails(pokemonName string) (*PokemonDetails, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pokémon data: %v", err)
	}
	defer resp.Body.Close()

	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get valid response for %s: status code %d", pokemonName, resp.StatusCode)
	}

	var details PokemonDetails
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Pokémon data: %v", err)
	}

	return &details, nil
}
