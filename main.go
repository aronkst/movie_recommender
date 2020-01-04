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
	fmt.Println("Movie Reccomender")
	fmt.Println()
	fmt.Println("[1] Add new movie")
	fmt.Println("[2] Remove a movie")
	fmt.Println()
	fmt.Println("Select option:")
	value := input()
	fmt.Println()
	switch value {
	case "1":
		fmt.Println("one")
	case "2":
		fmt.Println("two")
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
