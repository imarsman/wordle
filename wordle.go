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

var wordleWords = []string{} // slice to hold words
var triedItems letterSet     // each run has an increasing array of letters tried

func init() {
	triedItems = letterSet{}
	triedItems.items = &[]letterItem{}

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

type colourID int // to sort and keep track of colours

type letterItem struct {
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

type letterSet struct {
	items *[]letterItem
}

// newEmptyLetterSet get empty letter set
func newEmptyLetterSet() (ls letterSet) {
	ls.items = &[]letterItem{}

	return
}

// newSizedLetterSet new letterset with size
func newSizedLetterSet(size int) (ls letterSet) {
	ls.items = &[]letterItem{}
	*ls.items = make([]letterItem, size, size)
	return
}

func (ls *letterSet) String() (str string) {
	for _, v := range *ls.items {
		str = str + string(v.letter)
	}

	return
}

// add letter if it's new or change its colour if it's changed towards green
func (ls *letterSet) add(letter rune, colour colourID) {
	i := sort.Search(len(*ls.items), func(pos int) bool {
		return (*ls.items)[pos].letter >= letter
	})

	if len(*ls.items) > 0 {
		if i < len(*ls.items) && (*ls.items)[i].letter == letter {
			if (*ls.items)[i].colour < colour {
				(*ls.items)[i].colour = colour
			}
		} else {
			*ls.items = append(*ls.items, letterItem{letter: letter, colour: colour})
		}
	} else {
		*ls.items = append(*ls.items, letterItem{letter: letter, colour: colour})
	}

	sort.Slice(*ls.items, func(i int, j int) bool {
		return ((*ls.items)[i].letter) < ((*ls.items)[j].letter)
	})

	return
}

// // colourVector get a vector (slice) of colour names for the list of tried letters
// func (ls *letterSet) colourVector() (vector []string) {
// 	for _, v := range *ls.items {
// 		vector = append(vector, getColourName(v.colour))
// 	}

// 	return
// }

func (ls *letterSet) fillWithColour(colour colourID) {
	for _, v := range *ls.items {
		v.colour = colour
	}

	return
}

// func getFilledLetterSet(colourID colourID) (ls letterSet) {
// 	ls = newSizedLetterSet(wordLength)
// 	for i := range *ls.items {
// 		(*ls.items)[i].colour = colourID
// 	}

// 	return ls
// }

// func getFilledColourVector(color string) []string {
// 	colourVector := make([]string, wordLength, wordLength)
// 	for i := range colourVector {
// 		colourVector[i] = color
// 	}

// 	return colourVector
// }

func (ls *letterSet) filledColourVector(colourID colourID) {
	for _, v := range *ls.items {
		v.colour = colourID
	}

	return
}

// printWordLetters print letters for word with colour
func (ls *letterSet) printWordLetters() {
	if len(*ls.items) == 0 {
		return
	}
	for _, l := range *ls.items {
		switch l.colour {
		case greenColourID:
			fmt.Print("\033[42m\033[1;30m")
		case yellowColourID:
			fmt.Print("\033[43m\033[1;30m")
		case greyColourID:
			fmt.Print("\033[40m\033[1;37m")
		}
		fmt.Printf(" %c ", l.letter)
		fmt.Print("\033[m\033[m")
	}
}

// // func displayWord(word string, colourVector [wordLength]string) {
// func displayWord(word string, colourVector []string) {
// 	if len(word) == 0 {
// 		return
// 	}
// 	for i, c := range word {
// 		switch colourVector[i] {
// 		case greenColour:
// 			fmt.Print("\033[42m\033[1;30m")
// 		case yellowColour:
// 			fmt.Print("\033[43m\033[1;30m")
// 		case greyColour:
// 			fmt.Print("\033[40m\033[1;37m")
// 		}
// 		fmt.Printf(" %c ", c)
// 		fmt.Print("\033[m\033[m")
// 	}
// }

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

	guessesLetters := newSizedLetterSet(maxGuesses)
	// guessesLetters := letterSet{}
	// items := make([]letterItem, maxGuesses, maxGuesses)
	// guessesLetters.items = &items

	var guessCount int
	var guessesSet = make([]letterSet, 6, 6)

	var evaluate = func() {
		for guessCount = 0; guessCount < maxGuesses; guessCount++ {
			fmt.Printf("Enter your guess (%v/%v): ", guessCount+1, maxGuesses)

			guessWord, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			guessWord = strings.ToUpper(guessWord[:len(guessWord)-1])
			// fill in letters for items then later fill in colour as needed
			for i, v := range guessWord {
				(*guessesLetters.items)[i].letter = v
			}

			if guessWord == selectedWord {
				fmt.Println("You guessed right!")

				guessesLetters.filledColourVector(greenColourID)
				fmt.Println("Your wordle matrix is: ")
				for _, guess := range guessesSet {
					guess.printWordLetters()
					fmt.Println()
				}
				break
			} else {
				i := sort.SearchStrings(wordleWords, guessWord)
				if i < len(wordleWords) && wordleWords[i] == guessWord {
					for j, guessLetter := range guessWord {
						var cid colourID
						cid = greyColourID
						for k, letter := range selectedWord {
							if guessLetter == letter {
								if j == k {
									cid = greenColourID
									(*guessesLetters.items)[j].colour = greenColourID
									break
								} else {
									cid = yellowColourID
									(*guessesLetters.items)[j].colour = yellowColourID
								}
							}
						}
						triedItems.add(guessLetter, cid)
					}
					guessesSet = append(guessesSet, guessesLetters)
					guessesLetters.printWordLetters()
					fmt.Print(" [")
					triedItems.printWordLetters()
					fmt.Print("]")
					fmt.Println()
				} else {
					guessCount--
					fmt.Printf("%s not found in list. Please guess a valid %v letter word from the wordlist", guessWord, wordLength)
					fmt.Println()
				}
			}
			if guessCount == maxGuesses {
				fmt.Println("Better luck next time!")
				guessesLetters.fillWithColour(greenColourID)
				fmt.Print("The correct word is : ")
				guessesLetters.printWordLetters()
				fmt.Println()
			}
		}
	}

	evaluate()
}
