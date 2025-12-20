package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jakshi/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	nextURL     string
	previousURL string
	cache       *pokecache.Cache
}

func commandExit(c *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, cmd := range commands {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(c *config, args []string) error {

	var in_cache bool
	var body []byte

	if body, in_cache = c.cache.Get(c.nextURL); !in_cache {

		res, err := http.Get(c.nextURL)
		if err != nil {
			return fmt.Errorf("failed to fetch data: %v", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			return fmt.Errorf("received non-2xx response code: %d", res.StatusCode)
		}
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		c.cache.Add(c.nextURL, body)
		fmt.Println("data retrieved from API")
	} else {
		fmt.Println("data retrieved from cache")
	}

	var locationData struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	err := json.Unmarshal(body, &locationData)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	for _, location := range locationData.Results {
		fmt.Println(location.Name)
	}
	c.previousURL = locationData.Previous
	c.nextURL = locationData.Next

	return nil
}

func commandMapB(c *config, args []string) error {

	if c.previousURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	var in_cache bool
	var body []byte

	if body, in_cache = c.cache.Get(c.nextURL); !in_cache {
		res, err := http.Get(c.previousURL)
		if err != nil {
			return fmt.Errorf("failed to fetch data: %v", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			return fmt.Errorf("received non-2xx response code: %d", res.StatusCode)
		}
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		c.cache.Add(c.previousURL, body)
		fmt.Println("data retrieved from API")
	} else {
		fmt.Println("data retrieved from cache")
	}

	var locationData struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	err := json.Unmarshal(body, &locationData)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	for _, location := range locationData.Results {
		fmt.Println(location.Name)
	}
	c.previousURL = locationData.Previous
	c.nextURL = locationData.Next

	return nil
}

func commandExplore(c *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide an area name: explore <area-name>")
	}
	areaName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)

	// Check cache first
	if data, ok := c.cache.Get(url); ok {
		return printPokemon(data)
	}

	// Fetch from API
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Cache it
	c.cache.Add(url, body)

	return printPokemon(body)
}

func printPokemon(data []byte) error {
	var areaData struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}

	if err := json.Unmarshal(data, &areaData); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range areaData.PokemonEncounters {
		fmt.Printf("  - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 location areas in the Pokemon world",
			callback:    commandMapB,
		},
		"explore": cliCommand{
			name:        "explore",
			description: "Takes the name of a location area as an argument. List of all the Pokemon located there.",
			callback:    commandExplore,
		},
	}
}
