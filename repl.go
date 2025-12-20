package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	ltext := strings.ToLower(text)
	parts := strings.Fields(ltext)
	return parts
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		input_cleaned := cleanInput(input)

		if len(input_cleaned) == 0 {
			continue
		}
		cmdName := input_cleaned[0]
		args := []string{}
		if len(input_cleaned) > 1 {
			args = input_cleaned[1:]
		}

		if cmd, exists := commands[cmdName]; exists {
			err := cmd.callback(cfg, args)
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", cmdName, err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", cmdName)
		}

	}
}
