package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type searchMovie struct {
	Title string
	Year  int64
	IMDB  string
}

func getSearchMovies(search string) []searchMovie {
	search = strings.ReplaceAll(search, " ", "%20")
	url := fmt.Sprintf("https://www.imdb.com/find?q=%s&s=tt&ttype=ft", search)
	document, err := loadSite(url)
	if err != nil {
		panic(err)
	}

	var movies []searchMovie

	document.Find("tr[class*='findResult']").Each(func(i int, s *goquery.Selection) {
		searchMovie := searchMovie{
			Title: getTitleFromSiteToSearchMovie(s),
			Year:  getYearFromSiteToSearchMovie(s),
			IMDB:  getIMDBFromSiteToSearchMovie(s),
		}
		movies = append(movies, searchMovie)
	})

	return movies
}

func getTitleFromSiteToSearchMovie(selection *goquery.Selection) string {
	title := getValueFromSiteSelection(selection, "td.result_text a", "")
	return title
}

func getYearFromSiteToSearchMovie(selection *goquery.Selection) int64 {
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

func getIMDBFromSiteToSearchMovie(selection *goquery.Selection) string {
	imdb := getValueFromSiteSelection(selection, "td.result_text a", "href")
	imdb = regexReplace(imdb, `(\/\?ref_=fn_ft_tt_)([0-9]*).*?`, "")
	imdb = strings.ReplaceAll(imdb, "/title/", "")
	return imdb
}
