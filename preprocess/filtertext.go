package preprocess

import (
	"strings"
)

func FilterText(text string) string {
	var concatenated string = strings.ReplaceAll(text, "\n", " ")
	return concatenated
}