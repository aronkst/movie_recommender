package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func clearString(value string) string {
	newValue := strings.TrimSpace(value)
	newValue = strings.ReplaceAll(newValue, "\n", "")
	return strings.ReplaceAll(newValue, "\r", "")
}

func urlIMDB(imdb string) string {
	return fmt.Sprintf("https://www.imdb.com/title/%s", imdb)
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
	newValue := strings.Replace(value, ".", "", -1)
	return strings.Replace(newValue, ",", "", -1)
}

func formatSearchValue(search string) string {
	newSearch := strings.TrimSpace(search)
	newSearch = strings.ReplaceAll(newSearch, " ", "+")
	return newSearch
}

func urlIMDBSearch(search string) string {
	searchFormated := formatSearchValue(search)
	return fmt.Sprintf("https://www.imdb.com/find?q=%s&s=tt&ttype=ft", searchFormated)
}
