package service

import (
	"regexp"
	"strings"
)

// Match all non letter characters (including multi byte letters) and all non numbers
var reNonPalindromeChars = regexp.MustCompile(`[^\p{L}0-9]`)

// IsPalindrome checks if some text is a palindrome
// A palindrome is a word, number, phrase, or other sequence of characters which reads the same backward as forward, such as madam, racecar.
// There are also numeric palindromes, including date/time stamps using short digits 11/11/11 11:11 and long digits 02/02/2020.
// Sentence-length palindromes ignore capitalization, punctuation, and word boundaries.
// More details on: https://en.wikipedia.org/wiki/Palindrome
func IsPalindrome(text string) bool {
	text = strings.ToLower(text)
	runes := []rune(reNonPalindromeChars.ReplaceAllString(text, ""))
	length := len(runes)
	elementsToCheck := int(length / 2)

	result := true
	for i := 0; i < elementsToCheck; i++ {
		if runes[i] != runes[length - (i + 1)] {
			result = false
			break
		}
	}

	return result
}
