package main

import (
	"regexp"
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
	Points            int64
	Genres            []string
	Cover             string
	CoverSmall        string
	RecommendedMovies []string
	RecommendedBy     []string
}

func getMovie(imdb string) movie {
	url := urlIMDB(imdb)
	document, err := loadSite(url)
	if err != nil {
		panic(err)
	}

	var points int64

	score := getScoreToMovie(document)
	amountOfVotes := getAmountOfVotesToMovie(document)
	metascore := getMetascoreToMovie(document)
	if metascore <= 0 {
		points = int64(float64((score+float64(metascore/10))/2) * float64(amountOfVotes))
	} else {
		points = int64(score * float64(amountOfVotes))
	}

	return movie{
		IMDb:              imdb,
		Title:             getTitleToMovie(document),
		Year:              getYearToMovie(document),
		Summary:           getSummaryToMovie(document),
		Score:             score,
		AmountOfVotes:     amountOfVotes,
		Metascore:         metascore,
		Points:            points,
		Genres:            getGenresToMovie(document),
		Cover:             getCoverToMovie(document),
		CoverSmall:        getCoverSmallToMovie(document),
		RecommendedMovies: getRecommendedMoviesToMovie(document),
		RecommendedBy:     []string{},
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

func getCoverToMovie(document *goquery.Document) string {
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

func getCoverSmallToMovie(document *goquery.Document) string {
	return getValueFromSiteDocument(document, "div.poster a img", "src")
}

func getRecommendedMoviesToMovie(document *goquery.Document) []string {
	var recommendedMovies []string

	document.Find("div.rec_item").Each(func(i int, s *goquery.Selection) {
		recommendedMovies = append(recommendedMovies, getValueFromSiteSelectionInside(s, "data-tconst"))
	})

	return recommendedMovies
}
