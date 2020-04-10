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
		if folder.IsDir() && folder.Name() != ".git" {
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

func readHTML() string {
	file, err := ioutil.ReadFile("Recommended Movies.html")
	if err != nil {
		return ""
	}
	return replaceCommentHTML(string(file))
}

func readHTMLMovies() ([]movie, []movie) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(readHTML()))
	if err != nil {
		panic(err)
	}

	var watchedMovies, recommendedMovies []movie

	watchedMovies = loadMovieFromHTML(document, "div#WatchedMovies div#Movie")
	recommendedMovies = loadMovieFromHTML(document, "div#RecommendedMovies div#Movie")

	return watchedMovies, recommendedMovies
}

func loadMovieFromHTML(document *goquery.Document, selector string) []movie {
	var movies []movie

	document.Find(selector).Each(func(i int, s *goquery.Selection) {
		movie := movie{
			IMDb:              getValueFromSiteSelection(s, "p#IMDb", ""),
			Title:             getValueFromSiteSelection(s, "h2#Title", ""),
			Year:              stringToInt(getValueFromSiteSelection(s, "p#Year", "")),
			Summary:           getValueFromSiteSelection(s, "p#Summary", ""),
			Score:             stringToFloat(getValueFromSiteSelection(s, "p#Score", "")),
			AmountOfVotes:     stringToInt(getValueFromSiteSelection(s, "p#AmountOfVotes", "")),
			Metascore:         stringToInt(getValueFromSiteSelection(s, "p#Metascore", "")),
			Points:            stringToInt(getValueFromSiteSelection(s, "p#Points", "")),
			Genres:            strings.Split(getValueFromSiteSelection(s, "p#Genres", ""), ", "),
			Cover:             getValueFromSiteSelection(s, "p#Cover", ""),
			CoverSmall:        getValueFromSiteSelection(s, "img#CoverSmall", "src"),
			RecommendedMovies: strings.Split(getValueFromSiteSelection(s, "p#RecommendedMovies", ""), ", "),
			RecommendedBy:     strings.Split(getValueFromSiteSelection(s, "p#RecommendedBy", ""), ", "),
		}
		movies = append(movies, movie)
	})

	return movies
}
