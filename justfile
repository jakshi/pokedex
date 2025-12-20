# Default recipe - list all available commands
default:
    @just --list

# Run the pokedex CLI
run:
    go run .

# Run pokecache tests
test-pokecache:
    go test ./internal/pokecache

# Run pokecache tests with verbose output
test-pokecache-v:
    go test -v ./internal/pokecache

# Run all tests
test:
    go test ./...
