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

type colourID int

const (
	greyColourID colourID = iota
	yellowColourID
	greenColourID
)

func getColourName(c colourID) (result string) {
	switch c {
	case greyColourID:
		result = greyColour
	case yellowColourID:
		result = yellowColour
	case greenColourID:
		result = greenColour
	}

	return
}

type foundItem struct {
	colour colourID
	letter rune
}

type foundItemSet struct {
	items *[]foundItem
}

func (fis *foundItemSet) String() (str string) {
	for _, v := range *fis.items {
		str = str + string(v.letter)
	}

	return
}

func (fis *foundItemSet) add(letter rune, colour colourID) {
	i := sort.Search(len(*fis.items), func(pos int) bool {
		return (*fis.items)[pos].letter >= letter
	})

	if len(*fis.items) > 0 {
		if i < len(*fis.items) && (*fis.items)[i].letter == letter {
			if (*fis.items)[i].colour < colour {
				(*fis.items)[i].colour = colour
			}
		} else {
			*fis.items = append(*fis.items, foundItem{letter: letter, colour: colour})
		}
	} else {
		*fis.items = append(*fis.items, foundItem{letter: letter, colour: colour})
	}

	sort.Slice(*fis.items, func(i int, j int) bool {
		return ((*fis.items)[i].letter) < ((*fis.items)[j].letter)
	})
	return
}

func (fis *foundItemSet) colourVector() (vector []string) {
	for _, v := range *fis.items {
		vector = append(vector, getColourName(v.colour))
	}
	return
}

var foundItems foundItemSet

//go:embed sgb-words.txt
var words string

// var wordList = []string{}
var wordleWords = []string{}

// const WORDS_URL =
// "https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt"
const (
	wordLength   = 5
	maxGuesses   = 6
	greyColour   = "Grey"
	yellowColour = "Yellow"
	greenColour  = "Green"
)

func init() {
	foundItems = foundItemSet{}
	foundItems.items = &[]foundItem{}

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
	sort.Strings(wordleWords)
}

func getFilledColourVector(color string) []string {
	colourVector := make([]string, 5, 5)
	for i := range colourVector {
		colourVector[i] = color
	}
	return colourVector
}

// func displayWord(word string, colourVector [wordLength]string) {
func displayWord(word string, colourVector []string) {
	if len(word) == 0 {
		return
	}
	for i, c := range word {
		switch colourVector[i] {
		case greenColour:
			fmt.Print("\033[42m\033[1;30m")
		case yellowColour:
			fmt.Print("\033[43m\033[1;30m")
		case greyColour:
			fmt.Print("\033[40m\033[1;37m")
		}
		fmt.Printf(" %c ", c)
		fmt.Print("\033[m\033[m")
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	selectedWord := wordleWords[rand.Intn(len(wordleWords))]
	// selectedWord := "looks"
	// fmt.Println("selectedWord", selectedWord)

	reader := bufio.NewReader(os.Stdin)
	guesses := []map[string][]string{}
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
			colourVector := getFilledColourVector(greenColour)

			guesses = append(guesses, map[string][]string{guessWord: colourVector})

			fmt.Println("Your wordle matrix is: ")
			for _, guess := range guesses {
				for guessWord, colourVector := range guess {
					displayWord(guessWord, colourVector)
					fmt.Println()
				}
			}
			break
		} else {
			i := sort.SearchStrings(wordleWords, guessWord)
			if i < len(wordleWords) && wordleWords[i] == guessWord {
				colourVector := getFilledColourVector(greyColour)
				for j, guessLetter := range guessWord {
					var cid colourID
					cid = greyColourID
					for k, letter := range selectedWord {
						if guessLetter == letter {
							if j == k {
								cid = greenColourID
								colourVector[j] = greenColour
								break
							} else {
								cid = yellowColourID
								colourVector[j] = yellowColour
							}
						}
					}
					foundItems.add(guessLetter, cid)
				}
				guesses = append(guesses, map[string][]string{guessWord: colourVector})
				displayWord(guessWord, colourVector)
				fmt.Print(" [")
				letters := foundItems.String()
				vector := foundItems.colourVector()
				displayWord(letters, vector)
				fmt.Print("]")
				fmt.Println()
			} else {
				guessCount--
				fmt.Printf("%s not found in list. Please guess a valid %v letter word from the wordlist", guessWord, wordLength)
				fmt.Println()
			}
		}
	}

	if guessCount == maxGuesses {
		fmt.Println("Better luck next time!")
		colourVector := getFilledColourVector("Green")
		fmt.Print("The correct word is : ")
		displayWord(selectedWord, colourVector)
		fmt.Println()
	}
}
