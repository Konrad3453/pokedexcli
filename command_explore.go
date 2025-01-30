package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Encounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

type LocationArea struct {
	PokemonEncounters []Encounter `json:"pokemon_encounters"`
	Name              string      `json:"name"`
	URL               string      `json:"url"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationList struct {
	Results []Location `json:"results"`
}

type LocationAreaList struct {
	Results []LocationArea `json:"results"`
}



func getLocationAreaID(locationAreaName string) (int, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch location areas: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	var locationAreaList LocationAreaList
	err = json.Unmarshal(body, &locationAreaList)
	if err != nil {
		return 0, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	for _, locationArea := range locationAreaList.Results {
		if locationArea.Name == locationAreaName {
			var locationAreaID int
			fmt.Sscanf(locationArea.URL, "https://pokeapi.co/api/v2/location-area/%d/", &locationAreaID)
			return locationAreaID, nil
		}
	}

	return 0, fmt.Errorf("location area not found")
}

func listLocationAreas() ([]string, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location areas: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var locationAreaList LocationAreaList
	err = json.Unmarshal(body, &locationAreaList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	var locationAreaNames []string
	for _, locationArea := range locationAreaList.Results {
		locationAreaNames = append(locationAreaNames, locationArea.Name)
	}

	return locationAreaNames, nil
}

func commandExplore(cfg *config, args []string, secondArg interface{}) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a location to explore")
	}
	location := args[0]

	
	fmt.Printf("Exploring the location: %s...\n", location)

	locationAreaID, err := getLocationAreaID(location)
	if err != nil {
		
		locationAreas, listErr := listLocationAreas()
		if listErr != nil {
			return listErr
		}
		fmt.Println("Available location areas:")
		for _, area := range locationAreas {
			fmt.Println(area)
		}
		return err
	}

	
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationAreaID)

	
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch location data: %v", err)
	}
	defer resp.Body.Close()

	// parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal 
	var locationData LocationArea
	err = json.Unmarshal(body, &locationData)
	if err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if len(locationData.PokemonEncounters) == 0 {
		fmt.Println("No Pokémon found in this location.")
		return nil
	}

	fmt.Println("Pokémon found in this location:")
	for _, encounter := range locationData.PokemonEncounters {
		fmt.Printf("- %s\n ", encounter.Pokemon.Name)
	}

	return nil
}
