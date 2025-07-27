# Pokemon CLI Game

A command-line Pokemon game built in Go that allows users to explore locations, catch Pokemon, and manage their Pokedex using the PokeAPI. Features an internal caching system for improved performance.

## Features

- 🗺️ **Location Exploration** - Navigate through Pokemon world locations with map/mapb commands
- ⚡ **Pokemon Catching** - Throw Pokeballs at wild Pokemon with probability-based catch mechanics
- 📚 **Pokedex Management** - Inspect caught Pokemon and view your collection
- 🚀 **Intelligent Caching** - Custom pokecache package with automatic cleanup for faster responses
- 🎮 **Interactive CLI** - Clean command-line interface with help system

## Project Structure
```
POKEMON CLI GAME/
├── internal/
│ └── pokecache/
│ ├── pokecache.go # Cache implementation with mutex
│ └── pokecache_test.go # Unit tests for cache functionality
├── go.mod # Go module definition
├── main.go # Main CLI interface and commands
└── repl.log # REPL session log
```

## Installation

1. Clone this repository:
```bash
git clone https://github.com/yourusername/pokemon-cli-game.git
cd pokemon-cli-game
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the game:
```bash
go run .
```

## Usage

### Starting the Game

Run the application and use the interactive prompt:
```bash
go run .
```

The game will start with a `Pokedex > ` prompt where you can enter commands.

### Available Commands

Get help and see all commands
```
Pokedex > help
```

Navigate locations
```
Pokedex > map # Show next 20 locations
Pokedex > mapb # Show previous 20 locations
```

Explore a specific location
```
Pokedex > explore <location-name>

Example: Pokedex > explore great-marsh-area-1
Catch Pokemon
Pokedex > catch <pokemon-name>

Example: Pokedex > catch magikarp
Inspect caught Pokemon
Pokedex > inspect <pokemon-name>

Example: Pokedex > inspect magikarp
View your collection
Pokedex > pokedex
```

Exit the game
```
Pokedex > exit
```

### Example Game Session
```
Pokedex > map
====new http request====
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex > explore great-marsh-area-1
====new http request====
Exploring great-marsh-area-1...
Found Pokemon:

arbok

psyduck

magikarp

gyarados
...

Pokedex > catch magikarp
====new http request====
Throwing a Pokeball at magikarp...
...
magikarp was caught!
You may now inspect it with the inspect command.

Pokedex > inspect magikarp
Name: magikarp
Height: 9
Weight: 100
Stats:
hp: 20
attack: 10
defense: 55
...
```

## How It Works

1. **API Integration**: Fetches data from PokeAPI (https://pokeapi.co/) for locations and Pokemon details
2. **Caching System**: Custom pokecache package stores API responses with automatic expiration
3. **Catch Mechanics**: Probability-based catching system using Pokemon stats and random number generation
4. **State Management**: Maintains user's current location context and Pokedex collection
5. **Concurrent Safety**: Thread-safe cache implementation using mutex locks

## Technical Details

- **Go Version**: 1.21+
- **API Source**: PokeAPI v2 (https://pokeapi.co/api/v2/)
- **Cache Duration**: 10 seconds with automatic cleanup
- **Concurrency**: Goroutine-based cache reaping with ticker
- **Data Format**: JSON parsing for API responses

## Dependencies

This project uses only Go standard library packages:
- `bufio` - Buffered I/O for user input
- `encoding/json` - JSON parsing for API responses
- `net/http` - HTTP client for API requests
- `sync` - Mutex for thread-safe operations
- `time` - Time-based operations and cache expiration
- `math/rand` - Random number generation for catch mechanics
- `strings` - String manipulation and parsing
- `os` - Operating system interface

## Cache Implementation

The pokecache package provides:
```bash
type Cache struct {
cacheEntries map[string]cacheEntry
mu sync.Mutex
interval time.Duration
}
```

// Key methods
```bash
func NewCache(interval time.Duration) *Cache
func (c *Cache) Add(key string, val []byte)
func (c *Cache) Get(key string) ([]byte, bool)
```

**Features:**
- Automatic entry expiration based on configurable interval
- Thread-safe operations with mutex protection
- Background cleanup goroutine (reapLoop)
- Memory-efficient byte slice storage

## Error Handling

The application handles various scenarios:
- Invalid Pokemon names with graceful error messages
- Network connectivity issues with retry capability
- Malformed API responses with proper error reporting
- Unknown commands with help suggestions
- Empty Pokedex states with user guidance

## Testing

Run the test suite:
```bash
go test ./internal/pokecache
```

Test coverage includes:
- Cache add/get operations
- Automatic entry expiration
- Concurrent access safety
- Time-based cleanup functionality

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- Uses PokeAPI for Pokemon data and location information
- Inspired by classic Pokemon games and CLI applications
- Built with Go's excellent standard library for performance
- Custom caching implementation for educational purposes
