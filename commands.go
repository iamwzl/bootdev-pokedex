package main

import(
	"fmt"
	"os"
	"sort"
	"strings"
	"math"
	"math/rand"
	"github.com/StupidWeasel/bootdev-pokedex/internal/pokeapi"
)
type pokedexCommands struct {
	command			string
	description 	string
	callback    	func(client *pokeapi.PokeAPIClient, args ...string) error
}

var commands = make(map[string]pokedexCommands,0)
var commandsSorted = make([]string, 0)
func populateCommands(){

	commands["exit"] = pokedexCommands{
	        command:		"exit",
	        description:	"-- Exit the Pokedex",
	        callback:    	cmdExit}

	commands["help"] = pokedexCommands{
	        command:		"help",
	        description:	"{optional_command} -- Displays this help message!",
	        callback:    	cmdHelp}

	commands["map"] = pokedexCommands{
	        command:		"map",
	        description:	"-- View the next page of 20 Location Areas",
	        callback:    	cmdMap}

	commands["mapb"] = pokedexCommands{
	        command:		"mapb",
	        description:	"-- View the previous page of 20 Location Areas",
	        callback:    	cmdMapB}

	commands["explore"] = pokedexCommands{
	        command:		"explore",
	        description:	"{location} -- List all pokemon in this location",
	        callback:    	cmdExplore}

	commands["catch"] = pokedexCommands{
			command:		"catch",
	    	description:	"{pokemon} -- Attempt to catch the named pokemon",
	    	callback:    	cmdCatch}
	
	commands["inspect"] = pokedexCommands{
			command:		"inspect",
	    	description:	"{pokemon} -- Inspect a specific pokemon you have caught",
	    	callback:    	cmdInspect}

	commands["pokedex"] = pokedexCommands{
			command:		"pokedex",
	    	description:	"List all the pokemon you have caught",
	    	callback:    	cmdPokedex}

	
	for k,_ := range commands{
		commandsSorted = append(commandsSorted, k)
	}
	sort.Strings(commandsSorted)
}

func cmdExit(_ *pokeapi.PokeAPIClient, args ...string) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cmdHelp(_ *pokeapi.PokeAPIClient, args ...string) error{

	if len(args)>0{
		if command, ok := commands[args[0]]; ok {
		    fmt.Printf("  %s - %s\n", command.command, command.description)
			return nil
		}
	}
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _,key := range commandsSorted{
		command := commands[key]
		fmt.Printf("  %s %s\n", command.command, command.description)
	}
	return nil
}

func cmdMap(client *pokeapi.PokeAPIClient, args ...string) error{
	results,err := client.GetNextLocationAreas()
	if err != nil{
		return err
	}
	for _,location := range results.Results{
		fmt.Println(location.Name)
	} 
	return nil
}

func cmdMapB(client *pokeapi.PokeAPIClient, args ...string) error{
	response,err := client.GetPrevLocationAreas()
	if err != nil{
		return err
	}
	for _,location := range response.Results{
		fmt.Println(location.Name)
	} 
	return nil
}

func cmdExplore(client *pokeapi.PokeAPIClient, args ...string) error{
	
	if len(args)==0{
		return fmt.Errorf("No location specified")
	}

	response,err := client.GetNamedLocation(args[0])
	if err != nil{
		return err
	}
	fmt.Printf("Exploring %s...\n", args[0])
	for _, encounter := range response.PokemonEncounters{
		fmt.Println(" - " + encounter.Pokemon.Name)
	} 
	return nil
}

func cmdCatch(client *pokeapi.PokeAPIClient, args ...string) error{
	
	if len(args)==0{
		return fmt.Errorf("No pokemon specified")
	}

	target := strings.Join(args,"-")
	response,err := client.GetNamedPokemon(cleanPokemonName(target))
	if err != nil{
		return err
	}
	friendlyName := formatPokemonName(response.Name)
	fmt.Printf("Throwing a Pokeball at %s...\n", friendlyName)
	chance := int(math.Ceil(2 * float64(response.BaseExperience) / 70))
	if chance < 2 {
	    chance = 2
	}
	fmt.Printf("%.2f%% chance to catch %s\n", 100/float64(chance), friendlyName)
	if rand.Intn(chance) == 0 {
		client.Pokedex[target] = response
		fmt.Printf("%s was caught!\n", friendlyName)
	}else{
		fmt.Printf("%s esaped!\n", friendlyName)
	}
	return nil
}

func cmdInspect(client *pokeapi.PokeAPIClient, args ...string) error{

if len(args)==0{
		return fmt.Errorf("No pokemon specified")
	}

	target := strings.Join(args,"-")
	pokemon,ok := client.Pokedex[cleanPokemonName(target)]
	if !ok{
		return fmt.Errorf("You have not caugt one yet!")
	}

	fmt.Printf("Name: %s\n", formatPokemonName(pokemon.Name))
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _,stat := range pokemon.Stats{
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _,thisType := range pokemon.Types{
		fmt.Printf(" -%s\n", thisType.Type.Name)
	}
	return nil
}

func cmdPokedex(client *pokeapi.PokeAPIClient, args ...string) error{

	if len(client.Pokedex)==0{
		return fmt.Errorf("You have no pokemon, use \"catch {pokemon}\" to catch some")
	}
	fmt.Println("Your pokedex:")
	for _,pokemon := range client.Pokedex{
		fmt.Printf(" - %s\n", formatPokemonName(pokemon.Name))

	}
	return nil
}