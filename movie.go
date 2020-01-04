package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type movie struct {
	IMDb              string
	Title             string
	Year              int64
	Summary           string
	Score             float64
	AmountOfVotes     int64
	Metascore         int64
	Genres            []string
	Image             string
	RecommendedMovies []string
}

func getMovie(imdb string) movie {
	url := urlIMDB(imdb)
	document, err := loadSite(url)
	if err != nil {
		panic(err)
	}

	return movie{
		IMDb:              imdb,
		Title:             getTitleToMovie(document),
		Year:              getYearToMovie(document),
		Summary:           getSummaryToMovie(document),
		Score:             getScoreToMovie(document),
		AmountOfVotes:     getAmountOfVotesToMovie(document),
		Metascore:         getMetascoreToMovie(document),
		Genres:            getGenresToMovie(document),
		Image:             getImageToMovie(document),
		RecommendedMovies: getRecommendedMoviesToMovie(document),
	}
}

func getTitleToMovie(document *goquery.Document) string {
	title := getValueFromSiteDocument(document, "meta[property='og:title']", "content")
	title = strings.Replace(title, " - IMDb", "", 1)
	return regexReplace(title, `\s*\(.+\)$`, "")
}

func getYearToMovie(document *goquery.Document) int64 {
	year := getValueFromSiteDocument(document, "h1 span#titleYear a", "")
	return stringToInt(year)
}

func getSummaryToMovie(document *goquery.Document) string {
	return getValueFromSiteDocument(document, "div.summary_text", "")
}

func getScoreToMovie(document *goquery.Document) float64 {
	score := getValueFromSiteDocument(document, "div.ratingValue strong span[itemprop='ratingValue']", "")
	return stringToFloat(score)
}

func getAmountOfVotesToMovie(document *goquery.Document) int64 {
	amountOfVotes := getValueFromSiteDocument(document, "div.imdbRating a span[itemprop='ratingCount']", "")
	amountOfVotes = replacePointsAndCommas(amountOfVotes)
	return stringToInt(amountOfVotes)
}

func getMetascoreToMovie(document *goquery.Document) int64 {
	metascore := getValueFromSiteDocument(document, "div.metacriticScore.titleReviewBarSubItem span", "")
	return stringToInt(metascore)
}

func getGenresToMovie(document *goquery.Document) []string {
	var genres []string

	document.Find("div.see-more.inline.canwrap").Eq(1).Find("a").Each(func(i int, s *goquery.Selection) {
		genres = append(genres, getValueFromSiteSelectionInside(s, ""))
	})

	return genres
}

func getImageToMovie(document *goquery.Document) string {
	return getValueFromSiteDocument(document, "div.poster a img", "src")
}

func getRecommendedMoviesToMovie(document *goquery.Document) []string {
	var recommendedMovies []string

	document.Find("div.rec_item").Each(func(i int, s *goquery.Selection) {
		recommendedMovies = append(recommendedMovies, getValueFromSiteSelectionInside(s, "data-tconst"))
	})

	return recommendedMovies
}
