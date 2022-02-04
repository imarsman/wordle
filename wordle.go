package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

//go:embed sgb-words.txt
var words string

// var wordList = []string{}
var wordleWords = []string{}

// const WORDS_URL =
// "https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt"
const (
	wordLength = 5
	maxGuesses = 6
	grey       = "Grey"
	green      = "Green"
	yellow     = "Yellow"
)

func init() {
	scanner := bufio.NewScanner(strings.NewReader(words))
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == wordLength {
			wordleWords = append(wordleWords, strings.TrimSpace(strings.ToUpper(word)))
		}
	}
	err := scanner.Err()
	if err != nil {
		panic(err)
	}
	sort.Strings(wordleWords) // sorting is probably unecessary
}

func getFilledColourVector(color string) [wordLength]string {
	colourVector := [wordLength]string{}
	for i := range colourVector {
		colourVector[i] = color
	}
	return colourVector
}

func displayWord(word string, colourVector [wordLength]string) {
	for i, c := range word {
		switch colourVector[i] {
		case green:
			fmt.Print("\033[42m\033[1;30m")
		case yellow:
			fmt.Print("\033[43m\033[1;30m")
		case grey:
			fmt.Print("\033[40m\033[1;37m")
		}
		fmt.Printf(" %c ", c)
		fmt.Print("\033[m\033[m")
	}
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().Unix())

	selectedWord := wordleWords[rand.Intn(len(wordleWords))]

	reader := bufio.NewReader(os.Stdin)
	guesses := []map[string][wordLength]string{}
	var guessCount int
	for guessCount = 0; guessCount < maxGuesses; guessCount++ {
		fmt.Printf("Enter your guess (%v/%v): ", guessCount+1, maxGuesses)

		guessWord, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		guessWord = strings.ToUpper(guessWord[:len(guessWord)-1])

		if guessWord == selectedWord {
			fmt.Println("You guessed right!")
			colourVector := getFilledColourVector(green)

			guesses = append(guesses, map[string][wordLength]string{guessWord: colourVector})

			fmt.Println("Your wordle matrix is: ")
			for _, guess := range guesses {
				for guessWord, colourVector := range guess {
					displayWord(guessWord, colourVector)
				}
			}
			break
		} else {
			i := sort.SearchStrings(wordleWords, guessWord)
			if i < len(wordleWords) && wordleWords[i] == guessWord {
				colourVector := getFilledColourVector(grey)
				for j, guessLetter := range guessWord {
					for k, letter := range selectedWord {
						if guessLetter == letter {
							if j == k {
								colourVector[j] = green
								break
							} else {
								colourVector[j] = yellow
							}
						}
					}
				}
				guesses = append(guesses, map[string][wordLength]string{guessWord: colourVector})
				displayWord(guessWord, colourVector)
			} else {
				guessCount--
				fmt.Printf("Please guess a valid %v letter word from the wordlist", wordLength)
				fmt.Println()
			}
		}
	}

	if guessCount == maxGuesses {
		fmt.Println("Better luck next time!")
		colourVector := getFilledColourVector("Green")
		fmt.Print("The correct word is : ")
		displayWord(selectedWord, colourVector)
	}
}
