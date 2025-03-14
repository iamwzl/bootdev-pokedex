package main

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func formatPokemonName(name string) string {
	name = strings.ReplaceAll(name, "-", " ")
	titleCaser := cases.Title(language.Und)
	return titleCaser.String(name)
}

func cleanPokemonName(name string) string {
	runeReplace := map[rune]rune{'é': 'e', '♂': 'm', '♀': 'f', ' ': '-'}

	var cleanedNameBuilder strings.Builder
	for _, r := range strings.ToLower(name) {
		if val, ok := runeReplace[r]; ok {
			r = val
		}
		if (r >= 'a' && r <= 'z') || r == '-' {
			cleanedNameBuilder.WriteRune(r)
		}
	}
	return cleanedNameBuilder.String()
}
