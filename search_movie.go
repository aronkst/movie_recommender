package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type searchMovie struct {
	Title string
	Year  int64
	IMDB  string
}

func getSearchMovies(seach string) []searchMovie {
	url := urlIMDBSearch(seach)
	document, err := loadSite(url)
	if err != nil {
		panic(err)
	}

	var movies []searchMovie

	document.Find("tr[class*='findResult']").Each(func(i int, s *goquery.Selection) {
		searchMovie := searchMovie{
			Title: getTitleToSearchMovie(s),
			Year:  getYearToSearchMovie(s),
			IMDB:  getIMDBToSearchMovie(s),
		}
		movies = append(movies, searchMovie)
	})

	return movies
}

func getTitleToSearchMovie(selection *goquery.Selection) string {
	title := getValueFromSiteSelection(selection, "td.result_text a", "")
	return title
}

func getYearToSearchMovie(selection *goquery.Selection) int64 {
	yearString := getValueFromSiteSelection(selection, "td.result_text", "")
	yearString = regexReplace(yearString, "([[0-9]+)", ";$1")
	index := strings.Index(yearString, "(;") + 2
	if index == -1 {
		return 0
	}
	yearString = yearString[index : index+4]
	year := stringToInt(yearString)
	return year
}

func getIMDBToSearchMovie(selection *goquery.Selection) string {
	imdb := getValueFromSiteSelection(selection, "td.result_text a", "href")
	imdb = imdb[0:16]
	imdb = strings.ReplaceAll(imdb, "/title/", "")
	imdb = strings.ReplaceAll(imdb, "/", "")
	return imdb
}
