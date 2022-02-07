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
	"github.com/jwalton/gchalk"
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
	triedItems = newEmptyLetterSet()

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
	// *ls.items = make([]letterItem, size, size)
	for i := 0; i < size; i++ {
		li := new(letterItem)
		(*ls.items) = append((*ls.items), *li)
	}
	return
}

func (ls *letterSet) String() (str string) {
	for _, v := range *ls.items {
		str = str + string(v.letter)
	}

	return
}

// addWithColour letter if it's new or change its colour if it's changed towards green
func (ls *letterSet) addWithColour(letter rune, colour colourID) {
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

func (ls *letterSet) filledColourVector(colourID colourID) {
	for i := range *ls.items {
		(*ls.items)[i].colour = colourID
	}

	return
}

// printWordLetters print letters for word with colour
func (ls *letterSet) printWordLetters() {
	if len(*ls.items) == 0 {
		return
	}
	for _, l := range *ls.items {
		str := " " + string(l.letter) + " "
		switch l.colour {
		case greenColourID:
			fmt.Print(gchalk.BgGreen(str))
		case yellowColourID:
			fmt.Print(gchalk.BgYellow(str))
		case greyColourID:
			fmt.Print(gchalk.BgGrey(str))
		}
	}
}

// printWordLetters print letters for word with colour
func (ls *letterSet) printWordLettersBlank() {
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
		fmt.Printf(" %c ", ' ')
		fmt.Print("\033[m\033[m")
	}
}

// args CLI args
type args struct {
	Tries int  `arg:"-t" default:"6" help:"number of tries"`
	Show  bool `arg:"-s" help:"show word"`
	Blank bool `arg:"-b" help:"show try results with no letters"`
}

func main() {
	var callArgs args // initialize call args structure
	arg.MustParse(&callArgs)
	maxGuesses = callArgs.Tries

	rand.Seed(time.Now().Unix())

	selectedWord := wordleWords[rand.Intn(len(wordleWords))]

	if callArgs.Show {
		fmt.Println(selectedWord)
	}

	reader := bufio.NewReader(os.Stdin)

	var guessCount int
	var guessesSet = make([]letterSet, 0, 0)

tries:
	for guessCount = 0; guessCount < maxGuesses; guessCount++ {
		fmt.Printf("Enter your guess (%v/%v): ", guessCount+1, maxGuesses)

		guessWord, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		guessWord = strings.ToUpper(guessWord[:len(guessWord)-1]) // trim word and uc

		if len(guessWord) != wordLength {
			// fmt.Println(guessWord)
			fmt.Printf("The word you entered was %d letters. You need a word with %d letters\n", len(guessWord), wordLength)
			guessCount--
			continue tries
		}

		guessesLetters := newSizedLetterSet(wordLength)
		// fill in letters for items then later fill in colour as needed
		for i, v := range guessWord {
			(*guessesLetters.items)[i].letter = v
			(*guessesLetters.items)[i].colour = greyColourID
		}

		if guessWord == selectedWord {
			fmt.Println("You guessed right!")

			guessesLetters.filledColourVector(greenColourID)
			guessesSet = append(guessesSet, guessesLetters)
			// guessesSet[guessCount] = guessesLetters
			fmt.Println("Your wordle matrix is: ")
			for _, guess := range guessesSet {
				if callArgs.Blank {
					guess.printWordLettersBlank()
					fmt.Println()
				} else {
					guess.printWordLetters()
					fmt.Println()
				}
			}
			break
		} else {
			i := sort.SearchStrings(wordleWords, guessWord)
			if i < len(wordleWords) && wordleWords[i] == guessWord {
				for j, guessLetter := range guessWord {
					for k, letter := range selectedWord {
						if guessLetter == letter {
							if j == k {
								(*guessesLetters.items)[j].colour = greenColourID
								triedItems.addWithColour(guessLetter, greenColourID)
								break
							} else {
								(*guessesLetters.items)[j].colour = yellowColourID
								triedItems.addWithColour(guessLetter, yellowColourID)
							}
						}
					}
					// this will have no effect if higher colour already present
					triedItems.addWithColour(guessLetter, greyColourID)
				}
				guessesSet = append(guessesSet, guessesLetters)
				fmt.Print(gchalk.WithBold().Paint("Guess "))
				guessesLetters.printWordLetters()
				fmt.Print(gchalk.WithBold().Paint(" Tried "))
				triedItems.printWordLetters()
				fmt.Println()
			} else {
				guessCount--
				fmt.Printf("%s not found in list. Please guess a valid %v letter word from the wordlist\n", guessWord, wordLength)
			}
		}
		if guessCount == maxGuesses {
			fmt.Println("Better luck next time!")
			guessesLetters.filledColourVector(greenColourID)
			fmt.Print("The correct word is : ")
			guessesLetters.printWordLetters()
			fmt.Println()
		}
	}
}
