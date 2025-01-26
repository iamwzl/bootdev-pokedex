package main

import(
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/StupidWeasel/bootdev-pokedex/internal/pokeapi"
	"time"
)

func cleanInput(text string) []string{
	output := make([]string, 0)
	for _,word := range strings.Fields(text){
		output = append(output, strings.ToLower(word))
	}
	return output
}

func runREPL(){
	scanner := bufio.NewScanner(os.Stdin)
	client := pokeapi.NewPokeAPIClient(5 * time.Minute)
	for{
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Println("error reading standard input")
		}
		cleanedInput := cleanInput(scanner.Text())
		if len(cleanedInput)>0{
			if command, ok := commands[cleanedInput[0]]; ok {
					err := command.callback(client, cleanedInput[1:]...)
					if err != nil{
						fmt.Println(err)
					}
			}else{
				fmt.Println("Unknown command")
			}
		}
	}
}