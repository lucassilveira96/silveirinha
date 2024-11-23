package util

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

func ReplaceSpecialChars(input string) string {
	normalized := norm.NFD.String(input)

	var result strings.Builder
	for _, r := range normalized {
		if r <= 127 {
			result.WriteRune(r)
		}
	}

	return result.String()
}
