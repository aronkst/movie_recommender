package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", recommendedMovies)

	router.Run()
}

func recommendedMovies(context *gin.Context) {
	context.HTML(http.StatusOK, "list.tmpl", gin.H{
		"title": "Test",
	})
}

func menu() {
	clearScreen()

	fmt.Println("Movie Recommender.")
	fmt.Println("\n[1] Add a new movie.")
	fmt.Println("\n[2] Generate recommended movies HTML file.")

	value := input("\nSelect one of the options.")
	switch value {
	case "1":
		addNewMovie()
	case "2":
		createRecommendedMovies()
	default:
		main()
	}
}

func addNewMovie() {
	clearScreen()

	var name string
	for {
		name = input("What is the name of the movie?")
		if name != "" {
			break
		}
	}

	movies := getSearchMovies(name)
	for i, movie := range movies {
		fmt.Printf("[%d] %s (%d)\n", i+1, movie.Title, movie.Year)
		if i >= 19 {
			break
		}
	}

	countOptions := int64(len(movies))
	if countOptions > 20 {
		countOptions = 20
	}

	var value int64
	for {
		valueString := input("\nChoose movie from list above or enter [0] to perform a new search.")

		if valueString == "0" {
			addNewMovie()
		}

		value = stringToInt(valueString)
		if value >= 1 && value <= countOptions {
			break
		}
	}

	var date string
	for {
		date = input("When did you watch this movie? Enter the date as follows YYYYMMDD. [] for today, [0] for an unknown date or type the date in the format already informed.")

		if date == "" {
			dateTime := time.Now()
			date = dateTime.Format("20060102")
			break
		}

		if date == "0" {
			date = "00000000"
			break
		}

		_, err := time.Parse("20060102", date)
		if err == nil {
			break
		}
	}

	var like int64
	for {
		likeString := input("\nDid you like the movie? [0] for No and [1] for Yes.")
		if likeString == "0" || likeString == "1" {
			like = stringToInt(likeString)
			break
		}
	}

	movie := getMovie(movies[value-1].IMDB)
	downloadCover(movie, date, like)

	os.Exit(0)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func input(message string) string {
	fmt.Printf("%s\n\n", message)

	input := bufio.NewReader(os.Stdin)
	value, err := input.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Println()

	return clearString(value)
}
