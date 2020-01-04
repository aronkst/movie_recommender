package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func input() string {
	input := bufio.NewReader(os.Stdin)
	value, err := input.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return clearString(value)
}

func menu() {
	clearScreen()
	fmt.Println("Movie Recommender")
	fmt.Println()
	fmt.Println("[1] Add a new movie")
	fmt.Println("[2] Remove a movie")
	fmt.Println()
	fmt.Println("Select one of the options:")

	value := input()
	switch value {
	case "1":
		menuAddNewMovie()
	case "2":
		fmt.Println("2")
	}
}

func menuAddNewMovie() {
	clearScreen()
	fmt.Println("What is the name of the movie?")
	fmt.Println()
	name := input()
	fmt.Println()

	movies := getSearchMovies(name)
	for _, movie := range movies {
		fmt.Println(movie.Title)
	}
}
