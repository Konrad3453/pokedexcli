package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Pokemon struct {
	Name string `json:"name"`
}

func checkPokemonExists(pokemonName string) (bool, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to fetch Pokémon data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil // Pokémon does not exist
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var pokemon Pokemon
	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		return false, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return true, nil // Pokémon exists
}

func commandCatch(cfg *config, args []string, secondArg interface{}) error {
	if len(args) == 0 {
		return fmt.Errorf("no Pokémon specified")
	}

	pokemon := args[0]

	// Check if the Pokémon exists
	exists, err := checkPokemonExists(pokemon)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Pokémon not found: %s", pokemon)
	}

	// Pokémon exists, attempt to catch
	if catchFunction(pokemon) {
		fmt.Printf("%s was caught!", pokemon)
	} else {
		fmt.Printf("%s escaped!", pokemon)
	}
	
	return nil
}

func chance() int {
	rand.Seed(time.Now().UnixNano())
	percentage := rand.Intn(101)
	return percentage
}

func catchFunction(pokemon string) bool {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	time.Sleep(5 * time.Second)
	if chance() > 50 {
		return true // Pokémon caught
	} else {
		return false // Pokémon escaped
	}
}
