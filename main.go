package main

import (
	"time"

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

	router.POST("/watched-movies", postWatchedMovies)

	router.Run()
}

func returnJSONError(context *gin.Context, message string) {
	context.JSON(200, gin.H{
		"error": message,
	})
}

func postWatchedMovies(context *gin.Context) {
	imdb := context.PostForm("imdb")
	date := context.PostForm("date")
	likeString := context.PostForm("like")

	if imdb == "" || likeString == "" {
		returnJSONError(context, "invalid values")
	} else {
		date = formatDate(date)
		like := stringToInt(likeString)

		movie := addWatchedMovies(imdb, date, like)

		context.JSON(200, movie)
	}
}

func formatDate(date string) string {
	if date == "" {
		dateTime := time.Now()
		return dateTime.Format("20060102")
	} else if date == "0" {
		date = "00000000"
	}
	return date
}
