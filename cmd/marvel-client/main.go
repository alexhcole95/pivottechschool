package main

import (
	"fmt"
	"log"
	"pivottechschool/packages/marvel"
)

func main() {
	client := marvel.NewClient(marvel.BaseURL)
	chars, err := client.GetCharacters(5)
	if err != nil {
		log.Fatal(err)
	}

	for _, char := range chars {
		fmt.Printf("Name: %v | Description: %v\n", char.Name, char.Description)
	}
}
