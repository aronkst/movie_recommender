package main

import (
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

func replacePointsAndCommas(value string) string {
	value = strings.Replace(value, ".", "", -1)
	return strings.Replace(value, ",", "", -1)
}

func stringIsNumeric(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func uniqueValuesInArrayString(array []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range array {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func pagination(page int) int {
	if page <= 1 {
		return 0
	}
	return (page * 10) - 10
}
