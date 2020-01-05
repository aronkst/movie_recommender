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
				imdb := strings.Split(file.Name(), "__")[1]
				database = append(database, imdb)
			}
		}
	}

	return uniqueArrayString(database)
}

func readHTML() string {
	file, err := ioutil.ReadFile("Recommended Movies.html")
	if err != nil {
		panic(err)
	}
	return string(file)
}

func readRecommendedMovies() ([]string, []movie) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(readHTML()))
	if err != nil {
		panic(err)
	}

	var watchedMovies []string
	var recommendedMovies []movie

	watchedMovies = strings.Split(getValueFromSiteDocument(document, "watched-movies", ""), ", ")

	document.Find("div.recommended-movies div.movie").Each(func(i int, s *goquery.Selection) {
		movie := movie{
			IMDb:              getValueFromSiteSelection(s, "p.IMDb", ""),
			Title:             getValueFromSiteSelection(s, "h2.Title", ""),
			Year:              stringToInt(getValueFromSiteSelection(s, "p.Year", "")),
			Summary:           getValueFromSiteSelection(s, "p.Summary", ""),
			Score:             stringToFloat(getValueFromSiteSelection(s, "p.Score", "")),
			AmountOfVotes:     stringToInt(getValueFromSiteSelection(s, "p.AmountOfVotes", "")),
			Metascore:         stringToInt(getValueFromSiteSelection(s, "p.Metascore", "")),
			Points:            stringToInt(getValueFromSiteSelection(s, "p.Points", "")),
			Genres:            strings.Split(getValueFromSiteSelection(s, "p.Genres", ""), ", "),
			Cover:             getValueFromSiteSelection(s, "img.Cover", "src"),
			CoverSmall:        getValueFromSiteSelection(s, "img.CoverSmall", "src"),
			RecommendedMovies: strings.Split(getValueFromSiteSelection(s, "p.RecommendedMovies", ""), ", "),
		}
		recommendedMovies = append(recommendedMovies, movie)
	})

	return watchedMovies, recommendedMovies
}
