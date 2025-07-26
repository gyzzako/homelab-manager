package utils

import "strings"

func ReplaceMany(input string, replacements map[string]string) string {
	for key, val := range replacements {
		input = strings.ReplaceAll(input, key, val)
	}
	return input
}
