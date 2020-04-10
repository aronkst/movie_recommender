package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

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

func menu() {
	clearScreen()
	fmt.Println("Movie Recommender")
	fmt.Println("\n[1] Add a new movie")
	fmt.Println("\n[2] Remove a movie")
	fmt.Println("\n[3] Make the HTML")

	value := input("\nSelect one of the options")
	switch value {
	case "1":
		menuAddNewMovie()
	case "2":
		fmt.Println("2")
	case "3":
		makeHTML()
	default:
		menu()
	}
}

func menuAddNewMovie() {
	clearScreen()

	name := input("What is the name of the movie?")
	movies := getSearchMovies(name)
	for i, movie := range movies {
		fmt.Printf("[%d] %s (%d)\n", i+1, movie.Title, movie.Year)

		if i >= 19 {
			break
		}
	}

	value := stringToInt(input("\nChoose movie from list above or enter [0] to perform a new search."))
	if value <= 0 || value > 20 || value > int64(len(movies)) {
		menuAddNewMovie()
	} else {
		date := input("When did you watch this movie? Enter the date as follows YYYYMMDD")
		if _, err := strconv.Atoi(date); err != nil || len(date) != 8 {
			date = "00000000"
		}

		var like int64
		for {
			like = stringToInt(input("\nDid you like the movie? [0] for No and [1] for Yes"))
			if like == 0 || like == 1 {
				break
			}
		}

		movie := getMovie(movies[value-1].IMDB)
		downloadCover(movie, date, like)
	}
}
