package utils

import "unicode/utf8"

// written with help of ChatGPT

func IsOnlyRepeats(s, char string) bool {
	if utf8.RuneCountInString(char) != 1 || char == "" {
		return false // char must be exactly one rune
	}

	for _, r := range s {
		if string(r) != char {
			return false
		}
	}
	return len(s) > 0
}
