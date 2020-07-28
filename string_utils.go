package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func containsTrigger(message string, word string) bool {
	if word == "" {
		return false
	}
	n := strings.Count(message, word) + 1
	if n == 1 {
		return false
	}

	a := make([]string, n)
	n--
	i := 0
	for i < n {
		m := strings.Index(message, word)
		if m < 0 {
			return false
		}
		previousLetterCheck := true
		nextLetterCheck := true
		//is the trigger the start of a word
		if m == 0 { //the trigger word starts this block
			if i != 0 { // it's not the first block
				previousLetter, _ := utf8.DecodeLastRuneInString(a[i])
				if !isSeparator(previousLetter) {
					previousLetterCheck = false
				}
			}
		} else { //there are characters prior to the trigger
			previousLetter, _ := utf8.DecodeLastRuneInString(message[:m])
			if !isSeparator(previousLetter) {
				previousLetterCheck = false
			}
		}
		//is the trigger the end of the word
		if len(message) > m+len(word) {
			nextLetter, _ := utf8.DecodeRuneInString(message[m+len(word):])
			if !isSeparator(nextLetter) {
				previousLetterCheck = false
			}
		}
		if previousLetterCheck && nextLetterCheck {
			return true
		}

		a[i] = message[:m]
		message = message[m+len(word):]
		i++
	}
	a[i] = message
	return false
}

//Copied from "strings" library
func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}
