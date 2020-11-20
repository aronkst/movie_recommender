package main

import (
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func formatDate(date string) string {
	if date == "" {
		dateTime := time.Now()
		return dateTime.Format("20060102")
	} else if date == "0" {
		date = "00000000"
	}

	return date
}

func removeItemInSliceIfExistInSlice(slice1 []string, slice2 []string) []string {
	var final []string
	var exists bool

	for _, s1 := range slice1 {
		exists = false
		for _, s2 := range slice2 {
			if s1 == s2 {
				exists = true
			}
		}
		if !exists {
			final = append(final, s1)
		}
	}

	return final
}

func countPages(count int64) int64 {
	pages := float64(count) / float64(10)
	return int64(math.Ceil(pages))
}
