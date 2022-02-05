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

	"github.com/alexflint/go-arg"
)

//go:embed sgb-words.txt
var words string

var maxGuesses = 6

const (
	greyColourID colourID = iota
	yellowColourID
	greenColourID

	wordLength   = 5
	greyColour   = "Grey"
	yellowColour = "Yellow"
	greenColour  = "Green"
)

func init() {
	triedItems = triedItemSet{}
	triedItems.items = &[]triedItem{}

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

var triedItems triedItemSet // each run has an increasing array of letters tried

var wordleWords = []string{} // slice to hold words

type colourID int // to sort and keep track of colours

type triedItem struct {
	colour colourID
	letter rune
}

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

type triedItemSet struct {
	items *[]triedItem
}

func (fis *triedItemSet) String() (str string) {
	for _, v := range *fis.items {
		str = str + string(v.letter)
	}

	return
}

// add letter if it's new or change its colour if it's changed towards green
func (fis *triedItemSet) add(letter rune, colour colourID) {
	i := sort.Search(len(*fis.items), func(pos int) bool {
		return (*fis.items)[pos].letter >= letter
	})

	if len(*fis.items) > 0 {
		if i < len(*fis.items) && (*fis.items)[i].letter == letter {
			if (*fis.items)[i].colour < colour {
				(*fis.items)[i].colour = colour
			}
		} else {
			*fis.items = append(*fis.items, triedItem{letter: letter, colour: colour})
		}
	} else {
		*fis.items = append(*fis.items, triedItem{letter: letter, colour: colour})
	}

	sort.Slice(*fis.items, func(i int, j int) bool {
		return ((*fis.items)[i].letter) < ((*fis.items)[j].letter)
	})

	return
}

// colourVector get a vector (slice) of colour names for the list of tried letters
func (fis *triedItemSet) colourVector() (vector []string) {
	for _, v := range *fis.items {
		vector = append(vector, getColourName(v.colour))
	}

	return
}

func getFilledColourVector(color string) []string {
	colourVector := make([]string, wordLength, wordLength)
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

// args CLI args
type args struct {
	Tries int `arg:"-t" default:"6" help:"number of tries"`
}

func main() {
	var callArgs args // initialize call args structure
	arg.MustParse(&callArgs)
	maxGuesses = callArgs.Tries

	rand.Seed(time.Now().Unix())

	selectedWord := wordleWords[rand.Intn(len(wordleWords))]

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
					triedItems.add(guessLetter, cid)
				}
				guesses = append(guesses, map[string][]string{guessWord: colourVector})
				displayWord(guessWord, colourVector)
				fmt.Print(" [")
				letters := triedItems.String()
				vector := triedItems.colourVector()
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
