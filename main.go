package main

import (
	"strconv"
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

	router.GET("/watched-movies", getWatchedMovies)
	router.POST("/watched-movies", postWatchedMovies)
	router.GET("/recommended-movies", getRecommendedMovies)

	router.Run()
}

func returnJSONError(context *gin.Context, message string) {
	context.JSON(200, gin.H{
		"error": message,
	})
}

func getWatchedMovies(context *gin.Context) {
	offset := pageParam(context)
	title, summary, year, imdb, genre, score, metascore, order := searchParams(context)
	movies := listWatchedMovies(offset, title, summary, year, imdb, genre, score, metascore, order)
	context.JSON(200, movies)
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

func getRecommendedMovies(context *gin.Context) {
	// TODO
	// Pegar todos os filmes assistidos
	// montar todos em um array, com os seus IMDBs
	// Fazer uma consulta no banco, pegando a coluna RecommendedMovies desses filmes
	// Criar um array de IMDBs, será usados para pegar os filmes recomendados
	// Fazer loop em todos os RecommendedMovies, colocando cada imdb no array
	// Remover IMDBs que estão tambem na lista de filmes assistidos (não recomendar um filme que eu ja assisti)
	// Pegar todos os filmes que eu não quero assistir, do banco de dados
	// Remover do array de filmes recomendados, os itens que estão em filmes para não assistir
	// Remover IMDBs duplicados
	// Com esse array, utilizar no filtro da query abaixo
	// Depois criar a paginação
	// E filtro de pesquisa
	// TODO

	offset := pageParam(context)
	title, summary, year, imdb, genre, score, metascore, order := searchParams(context)

	movies := listRecommendedMovies(offset, title, summary, year, imdb, genre, score, metascore, order)

	context.JSON(200, movies)
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

func pageParam(context *gin.Context) int {
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
