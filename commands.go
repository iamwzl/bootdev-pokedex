package main

import(
	"fmt"
	"os"
	"sort"
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
	        description:	"Exit the Pokedex",
	        callback:    	cmdExit}

	commands["help"] = pokedexCommands{
	        command:		"help",
	        description:	"Displays this help message!",
	        callback:    	cmdHelp}

	commands["map"] = pokedexCommands{
	        command:		"map",
	        description:	"View the next page of 20 Location Areas",
	        callback:    	cmdMap}

	commands["mapb"] = pokedexCommands{
	        command:		"mapb",
	        description:	"View the previous page of 20 Location Areas",
	        callback:    	cmdMapB}

	
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
		fmt.Printf("  %s - %s\n", command.command, command.description)
	}
	return nil
}

func cmdMap(client *pokeapi.PokeAPIClient, args ...string) error{
	locations,err := client.GetNextLocationAreas()
	if err != nil{
		return err
	}
	for _,location := range locations{
		fmt.Println(location.Name)
	} 
	return nil
}

func cmdMapB(client *pokeapi.PokeAPIClient, args ...string) error{
	locations,err := client.GetPrevLocationAreas()
	if err != nil{
		return err
	}
	for _,location := range locations{
		fmt.Println(location.Name)
	} 
	return nil
}