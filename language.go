package kkcity

import (
	"errors"
	"strings"
)

// languages used to generate db column
var languages []string

// ErrLanguageIndex the language index is wrong.
var ErrLanguageIndex = errors.New("Language index wrong.")

// setupLanguage to setup current language.
// langs must follow ISO-639-1 (https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)
func setupLanguage(langs []string) {
	languages = nil
	for _, one := range langs {
		if len(one) != 2 {
			panic("Language length is not 2.")
		}
		languages = append(languages, strings.ToLower(one))
	}
}

// getLanguage a certain language with index.
func getLanguage(tp int) (string, error) {
	llength := len(languages)
	if tp < 0 || tp >= llength {
		return "", ErrLanguageIndex
	}
	return languages[tp], nil
}

// getAll to get all the languages.
func getAll() []string {
	return languages
}
