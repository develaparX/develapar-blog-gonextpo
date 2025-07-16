package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func GenerateSlug(title string) string {
	if title == "" {
		return ""
	}

	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace spaces and special characters with hyphens
	slug = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		if unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			return '-'
		}
		return -1
	}, slug)
	
	// Remove multiple consecutive hyphens
	re := regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")
	
	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")
	
	// If slug is empty after processing, return a default
	if slug == "" {
		return "untitled"
	}
	
	return slug
}
