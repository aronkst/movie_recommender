package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

func clearString(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\n", "")
	return strings.ReplaceAll(value, "\r", "")
}

func regexReplace(value string, regexOld string, regexNew string) string {
	regex := regexp.MustCompile(regexOld)
	return regex.ReplaceAllString(value, regexNew)
}

func stringToInt(value string) int64 {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}

	return number
}

func stringToFloat(value string) float64 {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}

	return number
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

func createFolderIfNotExists(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, os.ModePerm)
	}
}
