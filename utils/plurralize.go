package utils

import "github.com/gedex/inflector"

func ToPlural(word string) string {
	return inflector.Pluralize(word)
}

func ToSingular(word string) string {
	return inflector.Singularize(word)
}
