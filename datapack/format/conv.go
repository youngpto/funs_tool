package format

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var conv = cases.Title(language.Und)

func Title(input string) string {
	return conv.String(input)
}
