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

func main() {
	database, errDatabase = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if errDatabase != nil {
		panic(errDatabase)
	}

	database.AutoMigrate(&movie{})
	database.AutoMigrate(&notWatch{})
	database.AutoMigrate(&invalidMovie{})

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/watched-movies", getWatchedMovies)
	router.POST("/watched-movies", postWatchedMovies)
	router.GET("/recommended-movies", getRecommendedMovies)
	router.GET("/not-watch", getNotWatch)
	router.POST("/not-watch", postNotWatch)
	router.DELETE("/not-watch", deleteNotWatch)
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

	movies := listWatchedMovies(offset, title, summary, year, imdb, genre, score, metascore, order)

	context.JSON(200, movies)
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

	movies := listRecommendedMovies(offset, title, summary, year, imdb, genre, score, metascore, order)

	context.JSON(200, movies)
}

func getNotWatch(context *gin.Context) {
	offset := pageParams(context)
	title, summary, year, imdb, genre, score, metascore, order := searchParams(context)

	movies := listNotWatch(offset, title, summary, year, imdb, genre, score, metascore, order)

	context.JSON(200, movies)
}

func postNotWatch(context *gin.Context) {
	imdb := context.PostForm("imdb")

	if imdb == "" {
		jsonError(context, "invalid values")
		return
	}

	json := addNotWatch(imdb)

	context.JSON(200, json)
}

func deleteNotWatch(context *gin.Context) {
	imdb := context.PostForm("imdb")

	if imdb == "" {
		jsonError(context, "invalid values")
		return
	}

	json := removeNotWatch(imdb)

	context.JSON(200, json)
}

func getSearch(context *gin.Context) {
	title := context.DefaultQuery("title", "")

	if title == "" {
		jsonError(context, "invalid values")
	}

	movies := getSearchMovies(title)

	context.JSON(200, movies)
}

func pageParams(context *gin.Context) int {
	pageString := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageString) // TODO alterar metodos de conversão para valor final, stringToInt -> stringToInt64
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
