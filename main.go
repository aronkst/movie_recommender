package main

import (
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB
var errDatabase error

type response struct {
	Movies []movie
	Pages  int64
}

func main() {
	database, errDatabase = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if errDatabase != nil {
		panic(errDatabase)
	}

	database.AutoMigrate(&movie{})
	database.AutoMigrate(&blockedMovie{})
	database.AutoMigrate(&invalidMovie{})

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/watched-movies", getWatchedMovies)
	router.POST("/watched-movies", postWatchedMovies)
	router.GET("/recommended-movies", getRecommendedMovies)
	router.GET("/blocked-movies", getBlockedMovies)
	router.POST("/blocked-movies", postBlockedMovies)
	router.DELETE("/blocked-movies", deleteBlockedMovies)
	router.GET("/search", getSearch)

	router.Run()
}

func jsonError(context *gin.Context, message string) {
	context.JSON(200, gin.H{
		"error": message,
	})
}

func getWatchedMovies(context *gin.Context) {
	offset := pageParams(context)
	title, summary, year, imdb, genre, score, metascore, order := searchParams(context)

	movies, pages := listWatchedMovies(offset, title, summary, year, imdb, genre, score, metascore, order)

	context.JSON(200, response{
		Movies: movies,
		Pages:  pages,
	})
}

func postWatchedMovies(context *gin.Context) {
	imdb := context.PostForm("imdb")
	date := context.PostForm("date")
	likeString := context.PostForm("like")

	if imdb == "" || likeString == "" {
		jsonError(context, "invalid values")
		return
	}

	date = formatDate(date)
	like := stringToInt(likeString)
	movie := addWatchedMovies(imdb, date, like)

	context.JSON(200, movie)
}

func getRecommendedMovies(context *gin.Context) {
	offset := pageParams(context)
	title, summary, year, imdb, genre, score, metascore, order := searchParams(context)

	movies, pages := listRecommendedMovies(offset, title, summary, year, imdb, genre, score, metascore, order)

	context.JSON(200, response{
		Movies: movies,
		Pages:  pages,
	})
}

func getBlockedMovies(context *gin.Context) {
	offset := pageParams(context)
	title, summary, year, imdb, genre, score, metascore, order := searchParams(context)

	movies, pages := listBlockedMovies(offset, title, summary, year, imdb, genre, score, metascore, order)

	context.JSON(200, response{
		Movies: movies,
		Pages:  pages,
	})
}

func postBlockedMovies(context *gin.Context) {
	imdb := context.PostForm("imdb")

	if imdb == "" {
		jsonError(context, "invalid values")
		return
	}

	json := addBlockedMovie(imdb)

	context.JSON(200, json)
}

func deleteBlockedMovies(context *gin.Context) {
	imdb := context.PostForm("imdb")

	if imdb == "" {
		jsonError(context, "invalid values")
		return
	}

	json := removeBlockedMovie(imdb)

	context.JSON(200, json)
}

func getSearch(context *gin.Context) {
	title := context.DefaultQuery("title", "")

	if title == "" {
		jsonError(context, "invalid values")
		return
	}

	movies := getSearchMovies(title)

	context.JSON(200, movies)
}

func pageParams(context *gin.Context) int {
	pageString := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageString)
	offset := pagination(page)

	return offset
}

func searchParams(context *gin.Context) (string, string, int64, string, string, float64, int64, string) {
	title := context.DefaultQuery("title", "")
	summary := context.DefaultQuery("summary", "")
	yearString := context.DefaultQuery("year", "0")
	year := stringToInt(yearString)
	imdb := context.DefaultQuery("imdb", "")
	genre := context.DefaultQuery("genre", "")
	scoreString := context.DefaultQuery("score", "0")
	score := stringToFloat(scoreString)
	metascoreString := context.DefaultQuery("metascore", "0")
	metascore := stringToInt(metascoreString)
	order := context.DefaultQuery("order", "points")

	return title, summary, year, imdb, genre, score, metascore, order
}
