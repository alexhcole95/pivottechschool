package main

import (
	"fmt"
	"log"
	"pivottechschool/cmd/marvel"
)

func main() {
	client := marvel.NewClient(marvel.CharactersURL)
	chars, err := client.GetCharacters(5)
	if err != nil {
		log.Fatal(err)
	}

	for _, char := range chars {
		fmt.Println(char.Name)
		fmt.Println(char.Description)
		fmt.Println("\n")
	}
}
