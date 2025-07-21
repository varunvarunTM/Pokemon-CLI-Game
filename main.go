package main

import (
	"bufio"
	"strings"
	"os"
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"hellogo/internal/pokecache"
	"time"
	"math/rand"
)

type LocationAreaResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"` 
}

type Pokemon struct {
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []Stats `json:"stats"`
	Types []Types `json:"types"`
}

type Stats struct {
	BaseStat int `json:"base_stat"`
	Stat Stat `json:"stat"`
}

type Stat struct {
	Name string `json:"name"`
}

type Types struct {
	Slot int `json:"slot"`
	Type Type `json:"type"`
}

type Type struct {
	Name string `json:"name"`
}

type location struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []result `json:"results"`
}

type result struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type config struct {
	next string
	pokeCache *pokecache.Cache
	previous string
	pokedex map[string]Pokemon
}

type commandRegistry struct {
	name string
	description string
	callback func(*config,string) error
}

func commandMap(cfg *config, locationName string) error {
	nextLink := ""
	
	if cfg.next != "" {
		nextLink = cfg.next 
	} else {
		nextLink = "https://pokeapi.co/api/v2/location-area/"
	}
	
	var body []byte
	
	value,ok := cfg.pokeCache.Get(nextLink)
	if !ok {
		fmt.Println("====new http request====")
		resp,err := http.Get(nextLink)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
	
		body,err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cfg.pokeCache.Add(nextLink,body)

	} else {
		body = value
		fmt.Println("====fetched from pokecache====")
	} 

	var locations location
	if err := json.Unmarshal(body,&locations); err != nil {
		return err
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	for i := 0 ; i < len(locations.Results) && i < 20 ; i++ {
		fmt.Println(locations.Results[i].Name)
	}
	
	return nil
}

func commandMapb(cfg *config, locationName string) error {
	previousLink := ""

	if cfg.previous != "" {
		previousLink = cfg.previous 
	} else {
		fmt.Println("you're on the first page")
		return nil
	}
	
	var body []byte

	value,ok := cfg.pokeCache.Get(previousLink)
	if !ok {
		fmt.Println("====new http request====")
		resp,err := http.Get(previousLink)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		body,err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cfg.pokeCache.Add(previousLink,body)

	} else {
		fmt.Println("====fetched from pokecache====")
		body = value
	}

	var locations location
	if err := json.Unmarshal(body,&locations); err != nil {
		return err
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	for i := 0 ; i < len(locations.Results) && i < 20 ; i++ {
		fmt.Println(locations.Results[i].Name)
	}

	return nil
}

func commandExplore(cfg *config, locationName string) error {
	link := "https://pokeapi.co/api/v2/location-area/" + locationName

	var body []byte
	value,ok := cfg.pokeCache.Get(link)
	if !ok {
		fmt.Println("====new http request====")
		resp,err := http.Get(link)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		body,err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cfg.pokeCache.Add(link,body)

	} else {
		body = value
		fmt.Println("====fetched from pokecache====")
	}

	var locationResponse LocationAreaResponse
	if err := json.Unmarshal(body,&locationResponse); err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n",locationName)
	if len(locationResponse.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
	}

	for i := 0 ; i < len(locationResponse.PokemonEncounters) ; i++ {
		fmt.Println(" - " + locationResponse.PokemonEncounters[i].Pokemon.Name)
	}
	
	return nil
}

func commandCatch(cfg *config, pokemonName string) error {
	link := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	var body []byte
	value,ok := cfg.pokeCache.Get(link)
	if !ok {
		fmt.Println("====new http request====")
		resp,err := http.Get(link)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println("That pokemon does not exist.")
			return nil
		}

		body,err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cfg.pokeCache.Add(link,body)
	} else {
		body = value
		fmt.Println("====fetched from pokecache====")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	var pokemonDetail Pokemon
	if err := json.Unmarshal(body,&pokemonDetail); err != nil {
		return err
	}

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	
	throwPower := r.Intn(1000)

	var TotalStat int
	for _,Stat := range pokemonDetail.Stats {
		TotalStat += Stat.BaseStat
	}

	for i := 3 ; i > 0 ; i-- {
		fmt.Print("\r")
		for j := 0 ; j < 3 ; j++{
			fmt.Print(".")
			time.Sleep(250 * time.Millisecond)
		}
		fmt.Print("\r")
		for j := 0 ; j < 3 ; j++{
			fmt.Print(" ")
			time.Sleep(250 * time.Millisecond)
		}
	}

	fmt.Print("\r")
	if throwPower >= TotalStat {
		fmt.Printf("%s was caught!\n",pokemonName)
		cfg.pokedex[pokemonName] = pokemonDetail
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n",pokemonName)
		fmt.Println("Try Again")
	}

	return nil
}

func commandInspect(cfg *config, pokemonName string) error {
	pokemon,ok := cfg.pokedex[pokemonName]
	if !ok {
		fmt.Println("You have not caught that pokemon.")
		return nil
	}
	
	fmt.Printf("Name: %s\n",pokemonName)

	fmt.Printf("Height: %d\n",pokemon.Height)

	fmt.Printf("Weight: %d\n",pokemon.Weight)

	fmt.Println("Stats:")

	for _,Stat := range pokemon.Stats {
		fmt.Printf("  %s: %d\n",Stat.Stat.Name,Stat.BaseStat)
	}

	fmt.Println("Types:")
	for _,Type := range pokemon.Types {
		fmt.Printf("  - %s\n",Type.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, Name string) error {
	if len(cfg.pokedex) < 1 {
		fmt.Println("Empty Pokedex! Catch some pokemons.")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for key,_ := range cfg.pokedex {
		fmt.Printf(" - %s\n",key)
	}

	return nil
}

func commandHelp(cfg *config, Name string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	commandOrder := []string{"help","map","mapb","explore","catch","inspect","exit"}
	for _,commandName := range commandOrder {
		command := Registry[commandName]
		fmt.Printf("%s: %s\n", command.name , command.description )
	}
	return nil
}

func commandExit(cfg *config, Name string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

var Registry map[string]commandRegistry

func main() {
	
	cfg := &config{}
	cfg.pokedex = make(map[string]Pokemon)
	cfg.pokeCache = pokecache.NewCache(10 * time.Second)

	Registry = map[string]commandRegistry{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "Display's next 20 locations",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Display's previous 20 locations",
			callback: commandMapb,
		},
		"explore": {
			name: "explore",
			description: "list of all the Pokémon located there",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "throw pokeball at pokemon to catch",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "see details about a Pokemon caught",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "list of pokemons caught",
			callback: commandPokedex,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan(){
			break
		}
		word := strings.Fields(strings.ToLower(strings.TrimSpace(scanner.Text())))

		if len(word) == 0 {
			fmt.Println("No command entered. Type 'help' for options.")
			continue
		}

		command,ok := Registry[word[0]]
		if !ok {
			fmt.Println("Unknown command. Type 'help' for options")
			continue
		}

		switch len(word) {
		case 1:
			command.callback(cfg,"")
		case 2:
			Name := word[1]
			command.callback(cfg,Name)
		default:
			fmt.Println("Enter a valid name. Type 'help' for options.")
		}
	}
}