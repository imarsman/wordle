package main

import (
	"bufio"
	"bytes"
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

var wordleWords = []string{} // slice to hold words
var triedLetterSet letterSet // each run has an increasing array of letters tried

var maxGuesses int // max guesses - defaults to 6 and can be set

type colourID int // to sort and keep track of colours

const (
	greyColourID   colourID = iota // 0
	yellowColourID                 // 1
	greenColourID                  // 2

	wordLength = 5 // keep static for now
)

func init() {
	triedLetterSet = newEmptyLetterSet()

	// read in words from embedded list
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
	sort.Strings(wordleWords) // sort words
}

// letterItem a letter with a colour
type letterItem struct {
	letter rune
	colour colourID
}

// letterSet a list of letter items (each having a letter and a colourID)
type letterSet struct {
	items *[]letterItem
}

// newEmptyLetterSet get empty letter set
func newEmptyLetterSet() (ls letterSet) {
	ls.items = &[]letterItem{}

	return
}

// newSizedLetterSet new letterset with size
func newFilledLetterSet(word string) (ls letterSet) {
	ls.items = &[]letterItem{}
	// *ls.items = make([]letterItem, size, size)
	for _, v := range word {
		li := new(letterItem)
		li.letter = v
		(*ls.items) = append((*ls.items), *li)
	}

	return
}

// newSizedLetterSet new letterset with size
func newSizedLetterSet(size int) (ls letterSet) {
	ls.items = &[]letterItem{}
	for i := 0; i < size; i++ {
		li := new(letterItem)
		(*ls.items) = append((*ls.items), *li)
	}

	return
}

// String get string output for items
func (ls *letterSet) String() string {
	var buf bytes.Buffer

	for _, v := range *ls.items {
		buf.WriteRune(v.letter)
	}

	return buf.String()
}

// addLetterWithColour add new letter change colour if changed towards green
// The colour change relies on the numbering of the colourID, which is ordinal
// with grey first then yellow then green.
// The function relies on a sorted list to do the search
func (ls *letterSet) addLetterWithColour(letter rune, colour colourID) {
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

	// sort by letter
	sort.Slice(*ls.items, func(i int, j int) bool {
		return ((*ls.items)[i].letter) < ((*ls.items)[j].letter)
	})

	return
}

// fillColourVector fill all items in a letterSet with a colour
func (ls *letterSet) setAllLettersColour(colourID colourID) {
	for i := range *ls.items {
		// (*ls.items)[i].colour = colourID
		(*ls.items)[i].colour = colourID
	}

	return
}

// printLettersWithColour print letters for word with colour
func (ls *letterSet) printLettersWithColour() {
	if len(*ls.items) == 0 {
		return
	}
	for _, l := range *ls.items {
		str := " " + string(l.letter) + " "
		switch l.colour {
		case greenColourID:
			fmt.Print(gchalk.WithBold().BgGreen(str))
		case yellowColourID:
			fmt.Print(gchalk.WithBold().BgYellow(str))
		case greyColourID:
			fmt.Print(gchalk.WithBold().BgGrey(str))
		}
	}
}

// lettersIn count instances of a letter in list
func (ls *letterSet) lettersIn(letter rune) int {
	count := 0
	for _, v := range *ls.items {
		if v.letter == letter {
			count++
		}
	}

	return count
}

// clearBackwards clear backwards to grey for any earlier instances of letter starting
// at position.
func (ls *letterSet) clearBackward(targetLetter rune, startPosition int, maxToClear int) {
	countCleared := 0
	for i := startPosition; i >= 0; i-- {
		currentLetter := (*ls.items)[i].letter
		currentColour := (*ls.items)[i].colour
		// 	if the target letter and letter is currently yellow
		if currentLetter == targetLetter && currentColour == yellowColourID {
			if countCleared <= maxToClear {
				(*ls.items)[i].colour = greyColourID
			}
		}
		// If we in any way are on the target letter, count it as encountered
		if currentLetter == targetLetter {
			countCleared++
		}
	}
}

// printWordLetters print letters for word with colour
func (ls *letterSet) contains(letter rune) (found bool) {
	for _, v := range *(*ls).items {
		if v.letter == letter {
			found = true
			return
		}
	}
	return
}

// printWordLetters print letters for word with colour
func (ls *letterSet) printWordLettersBlank() {
	if len(*ls.items) == 0 {
		return
	}
	for _, l := range *ls.items {
		letter := "   "
		switch l.colour {
		case greenColourID:
			fmt.Print(gchalk.WithBgGreen().Green(string(letter)))
		case yellowColourID:
			fmt.Print(gchalk.WithBgYellow().Yellow(string(letter)))
		case greyColourID:
			fmt.Print(gchalk.WithBgGrey().Grey(string(letter)))
		}
	}
}

// args CLI args
type args struct {
	Tries      int    `arg:"-t" default:"6" help:"number of tries"`
	Show       bool   `arg:"-s" help:"show word"`
	Blank      bool   `arg:"-b" help:"show try results with no letters"`
	HideAnswer bool   `arg:"-H" help:"hide answer at end if not guessed"`
	UseAnswer  string `arg:"-u" help:"use provided answer"`
}

func main() {
	var callArgs args // initialize call args structure
	arg.MustParse(&callArgs)
	maxGuesses = callArgs.Tries

	rand.Seed(time.Now().Unix())

	wordToGuess := wordleWords[rand.Intn(len(wordleWords))]

	if callArgs.UseAnswer != "" {
		if len(callArgs.UseAnswer) != wordLength {
			fmt.Printf("Your manual word %s is not %d letters long. Exiting", callArgs.UseAnswer, wordLength)
			os.Exit(1)
		}
		wordToGuess = callArgs.UseAnswer
	}

	wordToGuess = strings.ToUpper(wordToGuess)

	if callArgs.Show {
		fmt.Println("Selected word", wordToGuess)
	}
	wordToGuessLetterSet := newFilledLetterSet(wordToGuess)

	reader := bufio.NewReader(os.Stdin)

	var guessCount int
	var guessesSet = make([]letterSet, 0, 0)
	var score = 0

tries:
	for guessCount = 0; guessCount < maxGuesses; guessCount++ {
		fmt.Printf("Enter your guess (%v/%v): ", guessCount+1, maxGuesses)

		guessWord, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		guessWord = strings.ToUpper(guessWord[:len(guessWord)-1]) // trim word and uc

		if len(guessWord) != wordLength {
			fmt.Printf("The word you entered was %d letters. You need a word with %d letters\n", len(guessWord), wordLength)
			guessCount--
			continue tries
		}

		guessesLetterSet := newSizedLetterSet(wordLength) // make sized slice

		// fill in letters for items then later fill in colour as needed
		for i, v := range guessWord {
			(*guessesLetterSet.items)[i].letter = v
			(*guessesLetterSet.items)[i].colour = greyColourID
		}

		if guessWord == wordToGuess {
			fmt.Println(gchalk.WithRed().Bold("\nYou guessed right!"))

			guessesLetterSet.setAllLettersColour(greenColourID)
			guessesSet = append(guessesSet, guessesLetterSet)

			fmt.Println("Your wordle matrix is: ")
			for _, guess := range guessesSet {
				if callArgs.Blank {
					guess.printWordLettersBlank()
					fmt.Println()
				} else {
					guess.printLettersWithColour()
					fmt.Println()
				}
			}
			triedNotThere := 0
			for _, v := range *triedLetterSet.items {
				if v.colour == greyColourID {
					triedNotThere++
				}
			}
			// calculate score
			score = len(guessesSet) + len(*triedLetterSet.items)
			fmt.Println()
			fmt.Printf("Your score is %d, %d guesses and %d letters tried\n", score, triedNotThere, len(*triedLetterSet.items))
			break
		} else {
			i := sort.SearchStrings(wordleWords, guessWord)
			if i < len(wordleWords) && wordleWords[i] == guessWord {

				for j, guessLetter := range guessWord {
					for k, letter := range wordToGuess {
						if guessLetter == letter {
							if j == k {
								// Set to green
								(*guessesLetterSet.items)[j].colour = greenColourID
								// Add to the tried letters set
								triedLetterSet.addLetterWithColour(guessLetter, greenColourID)
								break
							} else {
								// Set to yellow
								(*guessesLetterSet.items)[j].colour = yellowColourID
								// Add letter as yellow
								triedLetterSet.addLetterWithColour(guessLetter, yellowColourID)
							}
						}
					}
					// this will have no effect if higher colour already present
					triedLetterSet.addLetterWithColour(guessLetter, greyColourID)
				}

				// keep track of how many times each letter is encountered
				counts := make(map[rune]int)
				// Iterate backwards and decide whether to clear out previous non-green for the same letter
				for l := len(guessWord) - 1; l >= 0; l-- {
					currentLetter := (*guessesLetterSet.items)[l].letter
					currentColour := (*guessesLetterSet.items)[l].colour

					if currentColour == yellowColourID {
						countGuessedWord := guessesLetterSet.lettersIn(currentLetter)
						countWordToGuess := wordToGuessLetterSet.lettersIn(currentLetter)
						count, _ := counts[currentLetter]
						// Initialize count if not initialized
						if count == 0 {
							counts[currentLetter] = 1
							count = 0
						}

						// If we have more of the letter in the guessed word than the word to guess
						// just increment the count
						if (countGuessedWord - count) > countWordToGuess {
							counts[currentLetter] = counts[currentLetter] + 1
						} else {
							// Increment count
							counts[currentLetter] = counts[currentLetter] + 1
							// Figure out how many to reset to grey
							clearCount := countGuessedWord - countWordToGuess - 1
							// Reset to grey
							guessesLetterSet.clearBackward(currentLetter, l, clearCount)
						}
					}
				}
				// We have a set number of guesses
				guessesSet = append(guessesSet, guessesLetterSet)
				fmt.Print(gchalk.WithBold().Paint("Guess "))
				// Print out guess
				guessesLetterSet.printLettersWithColour()
				fmt.Print(gchalk.WithBold().Paint(" Tried "))
				// Print out tried letters
				triedLetterSet.printLettersWithColour()
				fmt.Println()
			} else {
				guessCount--
				fmt.Printf("%s not found in list. Please guess a valid %v letter word from the wordlist\n", guessWord, wordLength)
			}

		}
		// If we've run out of words, end and print out the correct word
		if guessCount+1 == maxGuesses && !callArgs.HideAnswer {
			fmt.Println(gchalk.WithBold().Paint("\nBetter luck next time!"))
			answer := newFilledLetterSet(wordToGuess)
			answer.setAllLettersColour(greenColourID)
			fmt.Print("The correct word is : ")
			answer.printLettersWithColour()
			fmt.Println()
		}
	}
}
