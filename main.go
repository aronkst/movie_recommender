package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	menu()
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
	fmt.Println()
	switch value {
	case "1":
		menuAddNewMovie()
	case "2":
		fmt.Println("2")
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func input() string {
	var value string
	_, err := fmt.Scanf("%s", &value)
	if err != nil {
		panic(err)
	}
	return value
}

func menuAddNewMovie() {
	clearScreen()
	fmt.Println("What is the name of the movie?")
	fmt.Println()
	name := input()
	fmt.Println()
	fmt.Println(name) // TODO REMOVE
}
