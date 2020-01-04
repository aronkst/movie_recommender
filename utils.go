package main

import "strings"

func clearString(value string) string {
	newValue := strings.TrimSpace(value)
	newValue = strings.ReplaceAll(newValue, "\n", "")
	return strings.ReplaceAll(newValue, "\r", "")
}
