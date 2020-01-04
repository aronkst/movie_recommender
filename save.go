package main

import (
	"io"
	"net/http"
	"os"
)

func downloadCover(movie movie) {
	image, err := http.Get(movie.Cover)
	if err != nil {
		panic(err)
	}
	defer image.Body.Close()

	file, err := os.Create("asdf.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, image.Body)
	if err != nil {
		panic(err)
	}
}
