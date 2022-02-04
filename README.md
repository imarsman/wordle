# Wordle

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/compatibility-emacs.svg)](https://forthebadge.com)

A terminal Wordle written in Go.

Modified a bit to use standard go style, to use an embedded word list, to avoid
duplication of word lists. The logic to do the colouring and the storing of the
colour data is great and I did not really touch that.

![Example](assets/example.png)

## Running

If you couldn't tell from the image, you need Go installed and then run:

```
go run wordle.go
```

## Pros

- Choose your wordlists
- If you're good you can choose longer words
- If you're bad you can increase number of guesses

## Cons

- Sometimes you get weird words with the default wordlist
