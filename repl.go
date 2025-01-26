package main

import(
	"bufio"
	"fmt"
	"os"
	"strings"
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
	for{
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Println(os.Stderr, "reading standard input:", err)
		}
		cleanedInput := cleanInput(scanner.Text())
		if len(cleanedInput)>0{
			if command, ok := commands[cleanedInput[0]]; ok {
					err := command.callback(cleanedInput[1:]...)
					if err != nil{
						fmt.Println(err)
					}
			}else{
				fmt.Println("Unknown command")
			}
		}
	}
}