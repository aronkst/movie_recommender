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
	movies, pages := recommendedMovies(page)

	context.HTML(http.StatusOK, "list.tmpl", gin.H{
		"Movies": movies,
		"Pages":  pages,
		"Page":   page,
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
		}

		if date == "0" {
			date = "00000000"
		}

		go addNewMovie(imdb, date, stringToInt(like))

		context.HTML(http.StatusOK, "save.tmpl", gin.H{
			"IMDb": imdb,
		})
	}
}

func addNewMovie(imdb string, date string, like int64) {
	movie := getMovie(imdb)
	go downloadCover(movie, date, like)

	for _, recommendedMovie := range movie.RecommendedMovies {
		go getMovie(recommendedMovie)
	}
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
