package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func readWatchedMovies() []string {
	var database []string

	folders, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, folder := range folders {
		if folder.IsDir() && folder.Name() != ".git" && folder.Name() != ".covers" && folder.Name() != ".data" {
			files, err := ioutil.ReadDir(fmt.Sprintf("./%s/", folder.Name()))
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				like := strings.Split(file.Name(), "__")[3]
				if like[0:1] == "1" {
					imdb := strings.Split(file.Name(), "__")[1]
					database = append(database, imdb)
				}
			}
		}
	}

	return uniqueArrayString(database)
}

func loadRecommendedMoviesHTMLFile() string {
	file, err := ioutil.ReadFile("Recommended Movies.html")
	if err != nil {
		return ""
	}
	return replaceCommentHTML(string(file))
}

func readWatchedAndRecommendedMoviesFromHTML() ([]movie, []movie) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(loadRecommendedMoviesHTMLFile()))
	if err != nil {
		panic(err)
	}

	var watchedMovies, recommendedMovies []movie

	watchedMovies = getMoviesFromHTML(document, "div#WatchedMovies div#Movie")
	recommendedMovies = getMoviesFromHTML(document, "div#RecommendedMovies div#Movie")

	return watchedMovies, recommendedMovies
}

func getMoviesFromHTML(document *goquery.Document, selector string) []movie {
	var movies []movie

	document.Find(selector).Each(func(i int, s *goquery.Selection) {
		movie := movie{
			IMDb:                getValueFromSiteSelection(s, "p#IMDb", ""),
			Title:               getValueFromSiteSelection(s, "h2#Title", ""),
			Year:                stringToInt(getValueFromSiteSelection(s, "td#Year", "")),
			Summary:             getValueFromSiteSelection(s, "td#Summary", ""),
			Score:               stringToFloat(getValueFromSiteSelection(s, "td#Score", "")),
			AmountOfVotes:       stringToInt(getValueFromSiteSelection(s, "td#AmountOfVotes", "")),
			Metascore:           stringToInt(getValueFromSiteSelection(s, "td#Metascore", "")),
			Points:              stringToInt(getValueFromSiteSelection(s, "p#Points", "")),
			Genres:              strings.Split(getValueFromSiteSelection(s, "td#Genres", ""), ", "),
			Cover:               getValueFromSiteSelection(s, "p#Cover", ""),
			CoverSmall:          getValueFromSiteSelection(s, "p#CoverSmall", ""),
			RecommendedMovies:   strings.Split(getValueFromSiteSelection(s, "p#RecommendedMovies", ""), ", "),
			RecommendedBy:       strings.Split(getValueFromSiteSelection(s, "p#RecommendedBy", ""), ", "),
			RecommendedByTitles: strings.Split(getValueFromSiteSelection(s, "p#RecommendedByTitles", ""), ", "),
		}
		movies = append(movies, movie)
	})

	return movies
}
