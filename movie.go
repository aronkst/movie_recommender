package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type movie struct {
	IMDb                string   `json:"imdb"`
	Title               string   `json:"title"`
	Year                int64    `json:"year"`
	Summary             string   `json:"summary"`
	Score               float64  `json:"score"`
	AmountOfVotes       int64    `json:"amount_of_votes"`
	Metascore           int64    `json:"metascore"`
	Points              int64    `json:"points"`
	Genres              []string `json:"genres"`
	Cover               string   `json:"cover"`
	CoverSmall          string   `json:"cover_small"`
	RecommendedMovies   []string `json:"recommended_movies"`
	RecommendedBy       []string `json:"recommended_by"`
	RecommendedByTitles []string `json:"recommended_by_titles"`
}

func getMovie(imdb string) movie {
	jsonFile := fmt.Sprintf("./.data/%s.json", imdb)
	if fileExists(jsonFile) {
		return getMovieFromJSON(jsonFile)
	}

	return getMovieFromSite(imdb)
}

func getMovieFromSite(imdb string) movie {
	url := fmt.Sprintf("https://www.imdb.com/title/%s", imdb)
	document, err := loadSite(url)
	if err != nil {
		panic(err)
	}

	score := getScoreFromSiteToMovie(document)
	amountOfVotes := getAmountOfVotesFromSiteToMovie(document)
	metascore := getMetascoreFromSiteToMovie(document)

	var points int64
	if metascore <= 0 {
		points = int64(score * float64(amountOfVotes))
	} else {
		points = int64(float64(float64(score+float64(float64(metascore)/10))/2) * float64(amountOfVotes))
	}

	movie := movie{
		IMDb:                imdb,
		Title:               getTitleFromSiteToMovie(document),
		Year:                getYearFromSiteToMovie(document),
		Summary:             getSummaryFromSiteToMovie(document),
		Score:               score,
		AmountOfVotes:       amountOfVotes,
		Metascore:           metascore,
		Points:              points,
		Genres:              getGenresFromSiteToMovie(document),
		Cover:               getCoverFromSiteToMovie(document),
		CoverSmall:          getCoverSmallFromSiteToMovie(document),
		RecommendedMovies:   getRecommendedMoviesFromSiteToMovie(document),
		RecommendedBy:       []string{},
		RecommendedByTitles: []string{},
	}

	createMovieJSONData(movie)
	downloadSmallCover(movie)

	return movie
}

func getMovieFromJSON(jsonFile string) movie {
	var movie movie

	file, err := os.Open(jsonFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &movie)

	return movie
}

func createMovieJSONData(movie movie) {
	file, err := json.MarshalIndent(movie, "", "  ")
	if err != nil {
		panic(err)
	}

	createFolderIfNotExists("./.data")

	jsonFile := fmt.Sprintf("./.data/%s.json", movie.IMDb)

	err = ioutil.WriteFile(jsonFile, file, 0644)
	if err != nil {
		panic(err)
	}
}

func getTitleFromSiteToMovie(document *goquery.Document) string {
	title := getValueFromSiteDocument(document, "meta[property='og:title']", "content")
	title = strings.Replace(title, " - IMDb", "", 1)
	return regexReplace(title, `\s*\(.+\)$`, "")
}

func getYearFromSiteToMovie(document *goquery.Document) int64 {
	year := getValueFromSiteDocument(document, "h1 span#titleYear a", "")
	return stringToInt(year)
}

func getSummaryFromSiteToMovie(document *goquery.Document) string {
	return getValueFromSiteDocument(document, "div.summary_text", "")
}

func getScoreFromSiteToMovie(document *goquery.Document) float64 {
	score := getValueFromSiteDocument(document, "div.ratingValue strong span[itemprop='ratingValue']", "")
	return stringToFloat(score)
}

func getAmountOfVotesFromSiteToMovie(document *goquery.Document) int64 {
	amountOfVotes := getValueFromSiteDocument(document, "div.imdbRating a span[itemprop='ratingCount']", "")
	amountOfVotes = replacePointsAndCommas(amountOfVotes)
	return stringToInt(amountOfVotes)
}

func getMetascoreFromSiteToMovie(document *goquery.Document) int64 {
	metascore := getValueFromSiteDocument(document, "div.metacriticScore.titleReviewBarSubItem span", "")
	return stringToInt(metascore)
}

func getGenresFromSiteToMovie(document *goquery.Document) []string {
	var genres []string

	document.Find("div.see-more.inline.canwrap").Eq(1).Find("a").Each(func(i int, s *goquery.Selection) {
		genres = append(genres, getValueFromSiteInsideSelection(s, ""))
	})

	return genres
}

func getCoverFromSiteToMovie(document *goquery.Document) string {
	regex, err := regexp.Compile(`"image": ".*?",`)
	if err != nil {
		panic(err)
	}

	html, err := document.Html()
	if err != nil {
		panic(err)
	}

	value := regex.FindString(html)
	value = strings.ReplaceAll(value, `"image": "`, "")
	return strings.ReplaceAll(value, `",`, "")
}

func getCoverSmallFromSiteToMovie(document *goquery.Document) string {
	return getValueFromSiteDocument(document, "div.poster a img", "src")
}

func getRecommendedMoviesFromSiteToMovie(document *goquery.Document) []string {
	var recommendedMovies []string

	document.Find("div.rec_item").Each(func(i int, s *goquery.Selection) {
		recommendedMovies = append(recommendedMovies, getValueFromSiteInsideSelection(s, "data-tconst"))
	})

	return recommendedMovies
}

func replacePointsAndCommas(value string) string {
	value = strings.Replace(value, ".", "", -1)
	return strings.Replace(value, ",", "", -1)
}
