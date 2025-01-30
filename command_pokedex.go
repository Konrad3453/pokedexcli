package main

import (
	
	"fmt"
	
)


func commandPokedex(cfg *config, args []string, secondArg interface{}) error {
	caughtPokemons := GetCaughtPokemons()
	if len(caughtPokemons) == 0 {
		fmt.Println("Pokedex is empty.")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for pokemonName := range caughtPokemons {
		fmt.Printf("- %s\n", pokemonName)
	}
	return nil
		
	
}