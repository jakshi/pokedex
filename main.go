package main

import (
	"time"

	"github.com/jakshi/pokedex/internal/pokecache"
)

func main() {
	cfg := &config{
		nextURL:     "https://pokeapi.co/api/v2/location-area/",
		previousURL: "",
		cache:       pokecache.NewCache(5 * time.Minute),
	}

	startRepl(cfg)
}
