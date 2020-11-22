package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type searchMovie struct {
	Title string
	Year  int64
	IMDb  string
	Image string
}

func getSearchMovies(search string) []searchMovie {
	var movies []searchMovie

	search = strings.ReplaceAll(search, " ", "%20")

	url := fmt.Sprintf("https://www.imdb.com/find?q=%s&s=tt&ttype=ft", search)
	document, err := loadSite(url)
	if err != nil {
		panic(err)
	}

	document.Find("tr[class*='findResult']").Each(func(i int, s *goquery.Selection) {
		searchMovie := searchMovie{
			Title: getTitleFromSiteToSearchMovie(s),
			Year:  getYearFromSiteToSearchMovie(s),
			IMDb:  getIMDBFromSiteToSearchMovie(s),
			Image: getImageFromSiteToSearchMovie(s),
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
	year := getValueFromSiteSelection(selection, "td.result_text", "")
	year = regexReplace(year, "([[0-9]+)", ";$1")

	index := strings.Index(year, "(;") + 2
	if index == -1 {
		return 0
	}

	year = year[index : index+4]
	return stringToInt(year)
}

func getIMDBFromSiteToSearchMovie(selection *goquery.Selection) string {
	imdb := getValueFromSiteSelection(selection, "td.result_text a", "href")
	imdb = regexReplace(imdb, `(\/\?ref_=fn_ft_tt_)([0-9]*).*?`, "")
	imdb = strings.ReplaceAll(imdb, "/title/", "")
	return imdb
}

func getImageFromSiteToSearchMovie(selection *goquery.Selection) string {
	title := getValueFromSiteSelection(selection, "td.primary_photo a img", "src")
	return title
}
