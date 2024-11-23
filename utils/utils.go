package utils

import (
	"strings"
	"unicode"
)

// ToPascalCase converts a string to PascalCase
func ToPascalCase(input string) string {
	// Capitalize the first letter of the input and keep the rest as is
	if len(input) > 0 {
		input = strings.ToUpper(string(input[0])) + input[1:]
	}
	return input
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(input string) string {
	// Make the first letter lowercase (for camelCase)
	if len(input) > 0 {
		input = strings.ToLower(string(input[0])) + input[1:]
	}
	return input
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(input string) string {
	var result strings.Builder
	for i, r := range input {
		// If the rune is uppercase and not the first one, insert an underscore
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		// Add the current rune as lowercase
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// ToUrlCase converts a string to url-case
func ToUrlCase(input string) string {
	var result strings.Builder
	for i, r := range input {
		// If the rune is uppercase and not the first one, insert an underscore
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('-')
		}
		// Add the current rune as lowercase
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
