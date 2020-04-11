package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadCover(movie movie, date string, like int64) {
	title := regexReplace(movie.Title, "[^a-zA-Z0-9 ]+", "")
	filename := fmt.Sprintf("%s__%s__%s__%d.jpg", date, movie.IMDb, title, like)
	filename = fmt.Sprintf("./%s/%s", date[0:4], filename)
	downloadImage(movie.Cover, date[0:4], filename)
}

func downloadSmallCover(movie movie) {
	filename := fmt.Sprintf("./.covers/%s.jpg", movie.IMDb)
	downloadImage(movie.CoverSmall, "./.covers", filename)
}

func downloadImage(url string, folder string, filename string) {
	createFolderIfNotExists(folder)

	if fileExists(filename) {
		return
	}

	if url == "" {
		return
	}

	image, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer image.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, image.Body)
	if err != nil {
		panic(err)
	}
}
