package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

// choices between cars, dogs, and US Presidents.
var countryData = []string{"ghana", "togo", "benin", "kenya", "uganda"}

// the average game of hangman has 7 turns.
const lives = 7

// function that returns user input
func getInput() string {
	var userInput string
	fmt.Println("")

	data, err := fmt.Scanf("%s \n", &userInput)
	if err != nil || data != 1 {
		return "Error"
	}

	return userInput
}

// this function will need to be repeated three times. Dogs, Cars, and US Presidents.
func getRandomCountry() string {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(countryData)
	return countryData[rand.Intn(max-min)]
}

// this function takes a slice and looks for an element in it. If found it will
// return its key, otherwise it will return -1 and a bool of false.
func findInput(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func getCharacterPositions(char string, word string) []int {
	var pos []int

	for i, ch := range word {
		if char == fmt.Sprintf("%c", ch) {
			pos = append(pos, i)
		}
	}

	return pos
}

func guessedResult(slice []string, word string) []string {
	length := utf8.RuneCountInString(word)
	array := make([]string, length)

	// fill the entire array with "_"
	for id := range array {
		array[id] = "_"
	}

	// for each letter in slice get the positions
	for i := 0; i < len(slice); i++ {
		positions := getCharacterPositions(slice[i], word)

		for j := 0; j < len(positions); j++ {
			array[positions[j]] = slice[i]
		}
	}

	return array
}

func countUnique(word string) int {
	var characters []string

	// create an array to store all values characters
	for _, ch := range word {

		// convert rune to string
		s := fmt.Sprintf("%c", ch)

		// check if character exists
		k, _ := findInput(characters, s)
		if k == -1 {
			// add if yes
			characters = append(characters, s)
		} else {
			// skip if no
			continue
		}
	}

	return len(characters)
}

func main() {
	var usedChars []string
	var correctGuess []string
	var incorrectGuess []string
	fmt.Println("WELCOME TO HANGMAN! PLEASE PRESS 'S' AND THEN ENTER TO BEGIN!")

	// get input to start
	data := getInput()

	// if returned input is not s keep asking
	for data != "s" {
		fmt.Println("I'M SORRY, THAT IS NOT A VALID INPUT.\nPLEASE PRESS 'S' AND THEN ENTER TO BEGIN!")
		data = getInput()
	}

	// if input is s
	if data == "s" {
		fmt.Println("GAME BEGUN - GOOD LUCK!")
	}

	// sets the word to be guessed
	guess := getRandomCountry()

	// count the unique items
	uc := countUnique(guess)

	fmt.Println("GUESS THE COUNTRY")

	for {
		fmt.Println(guessedResult(correctGuess, guess))

		value := getInput()

		k, _ := findInput(usedChars, value)
		if k == -1 {
			fmt.Println("")
			usedChars = append(usedChars, value)
			fmt.Println("")
		} else {
			fmt.Println("VALUE ALREADY IN USE, PLEASE TRY AGAIN")
			continue
		}

		if strings.Contains(guess, value) {
			correctGuess = append(correctGuess, value) // if guess is correct store it
			fmt.Printf(" %s IS A CORRECT LETTER, GOOD JOB!", value)
		} else {
			incorrectGuess = append(incorrectGuess, value) // if guess is correct store it
			fmt.Printf(" %s IS NOT A CORRECT GUESS, UH-OH!", value)
		}
		fmt.Println("")
		fmt.Println("")

		fmt.Print("NUMBER OF LIVES LEFT: ")
		fmt.Println(lives - len(incorrectGuess))

		fmt.Print("INCORRECT GUESSES: ")
		fmt.Println(incorrectGuess)

		fmt.Print("CORRECT GUESSES: ")
		fmt.Println(correctGuess)

		if len(incorrectGuess) == lives {
			break
		}

		if len(correctGuess) == uc {
			break
		}
	}

	if len(incorrectGuess) == lives {
		fmt.Println("")
		fmt.Println("----------------YOU DIED!----------------")
		fmt.Printf("THE COUNTRY WAS %s", guess)
	}

	if len(correctGuess) == uc {
		fmt.Println("")
		fmt.Println("----------------YOU LIVED!----------------")
		fmt.Printf("THE COUNTRY WAS %s", guess)
	}
}
