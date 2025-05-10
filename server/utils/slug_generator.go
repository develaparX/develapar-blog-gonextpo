package utils

import (
	"strings"
	"unicode"
)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		if unicode.IsSpace(r) {
			return '-'
		}
		return -1
	}, slug)
	return slug
}
