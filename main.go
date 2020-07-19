package main

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"htmlGenres":            htmlGenres,
		"htmlIMDbTitle":         htmlIMDbTitle,
		"htmlExistPreviousPage": htmlExistPreviousPage,
		"htmlPreviousPage":      htmlPreviousPage,
		"htmlExistNextPage":     htmlExistNextPage,
		"htmlNextPage":          htmlNextPage,
	})
	router.LoadHTMLFiles(
		"./templates/list.tmpl",
		"./templates/search.tmpl",
		"./templates/add.tmpl",
		"./templates/save.tmpl",
	)
	router.StaticFS("./.covers", http.Dir("./.covers"))

	router.GET("/", getRecommendedMovies)
	router.GET("/search", getSearch)
	router.GET("/add", getAdd)
	router.GET("/save", getSave)

	router.Run()
}

func getRecommendedMovies(context *gin.Context) {
	pageString := context.DefaultQuery("page", "1")
	page := stringToInt(pageString)
	title, summary, year, imdb, genre, score := searchParams(context)
	movies, pages := recommendedMovies(page, title, summary, year, imdb, genre, score)

	context.HTML(http.StatusOK, "list.tmpl", gin.H{
		"Movies":  movies,
		"Pages":   pages,
		"Page":    page,
		"Title":   title,
		"Summary": summary,
		"Year":    year,
		"IMDb":    imdb,
		"Genre":   genre,
		"Score":   score,
	})
}

func getSearch(context *gin.Context) {
	var movies []searchMovie

	search := context.DefaultQuery("search", "")
	if search != "" {
		movies = getSearchMovies(search)
	}

	context.HTML(http.StatusOK, "search.tmpl", gin.H{
		"Search": search,
		"Movies": movies,
	})
}

func getAdd(context *gin.Context) {
	imdb := context.DefaultQuery("imdb", "")

	if imdb == "" {
		context.Redirect(http.StatusMovedPermanently, "/search")
	} else {
		context.HTML(http.StatusOK, "add.tmpl", gin.H{
			"IMDb": imdb,
		})
	}
}

func getSave(context *gin.Context) {
	imdb := context.DefaultQuery("imdb", "")
	date := context.DefaultQuery("date", "")
	like := context.DefaultQuery("like", "")

	if imdb == "" || like == "" {
		context.Redirect(http.StatusMovedPermanently, "/search")
	} else {
		if date == "" {
			dateTime := time.Now()
			date = dateTime.Format("20060102")
		} else if date == "0" {
			date = "00000000"
		}

		go addNewMovie(imdb, date, stringToInt(like))

		context.HTML(http.StatusOK, "save.tmpl", gin.H{
			"IMDb": imdb,
		})
	}
}

func searchParams(context *gin.Context) (string, string, int64, string, string, float64) {
	title := context.DefaultQuery("title", "")
	summary := context.DefaultQuery("summary", "")
	yearString := context.DefaultQuery("year", "0")
	year := stringToInt(yearString)
	imdb := context.DefaultQuery("imdb", "")
	genre := context.DefaultQuery("genre", "")
	scoreString := context.DefaultQuery("score", "0")
	score := stringToFloat(scoreString)

	return title, summary, year, imdb, genre, score
}

func htmlGenres(genres []string) string {
	return strings.Join(genres, ", ")
}

func htmlIMDbTitle(imdb string) string {
	return getMovie(imdb).Title
}

func htmlExistPreviousPage(page int64) bool {
	return page >= 2
}

func htmlPreviousPage(page int64) int64 {
	return page - 1
}

func htmlExistNextPage(page int64, pages []int) bool {
	return int(page) < len(pages)
}

func htmlNextPage(page int64) int64 {
	return page + 1
}
