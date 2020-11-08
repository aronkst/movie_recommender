package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getMovie(imdb string) (movie, error) {
	var invalidMovie invalidMovie
	whereInvalidMovie := database.Where("imdb = ?", imdb).First(&invalidMovie)

	if whereInvalidMovie.RowsAffected >= 1 {
		return movie{}, errors.New("invalid movie")
	}

	var movie movie
	whereMovie := database.Where("imdb = ?", imdb).First(&movie)

	if whereMovie.RowsAffected <= 0 {
		return getMovieFromSite(imdb)
	}

	return movie, nil
}

func getMovieFromSite(imdb string) (movie, error) {
	fmt.Println("ERROROROROROROR")
	fmt.Println(imdb)

	url := fmt.Sprintf("https://www.imdb.com/title/%s", imdb)
	document, err := loadSite(url)
	if err != nil {
		panic(err)
	}

	score := getScoreFromSiteToMovie(document)
	amountOfVotes := getAmountOfVotesFromSiteToMovie(document)
	metascore := getMetascoreFromSiteToMovie(document)
	urlCoverSmall := getCoverSmallFromSiteToMovie(document)

	var points int64
	if metascore <= 0 {
		points = int64(score * float64(amountOfVotes))
	} else {
		points = int64(float64(float64(score+float64(float64(metascore)/10))/2) * float64(amountOfVotes))
	}

	movie := movie{
		IMDb:              imdb,
		Title:             getTitleFromSiteToMovie(document),
		Year:              getYearFromSiteToMovie(document),
		Summary:           getSummaryFromSiteToMovie(document),
		Score:             score,
		AmountOfVotes:     amountOfVotes,
		Metascore:         metascore,
		Points:            points,
		Genres:            getGenresFromSiteToMovie(document),
		RecommendedMovies: getRecommendedMoviesFromSiteToMovie(document),
		URLCover:          getCoverFromSiteToMovie(document),
		URLCoverSmall:     urlCoverSmall,
		Cover:             "",
	}

	if validMovie(movie) {
		movie.Cover = getImageFromSiteToBase64(urlCoverSmall)
		database.Create(&movie)

		return movie, nil
	}

	invalidMovie := invalidMovie{
		IMDb: imdb,
	}
	database.Create(&invalidMovie)

	return movie, errors.New("invalid movie")
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

func getGenresFromSiteToMovie(document *goquery.Document) string {
	var genres []string

	document.Find("div.see-more.inline.canwrap").Eq(1).Find("a").Each(func(i int, s *goquery.Selection) {
		genres = append(genres, getValueFromSiteInsideSelection(s, ""))
	})

	return strings.Join(genres, ",")
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

func getRecommendedMoviesFromSiteToMovie(document *goquery.Document) string {
	var recommendedMovies []string

	document.Find("div.rec_item").Each(func(i int, s *goquery.Selection) {
		recommendedMovies = append(recommendedMovies, getValueFromSiteInsideSelection(s, "data-tconst"))
	})

	return strings.Join(recommendedMovies, ",")
}

func validMovie(movie movie) bool {
	return movie.URLCover != "" && movie.URLCoverSmall != "" &&
		movie.Score > 0 && movie.Year > 0 && movie.AmountOfVotes > 0 &&
		movie.Genres != "" && movie.Summary != "" &&
		movie.Summary != `Add a Plot »`
}
