package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func getFolderDownloadImage(date string) string {
	return fmt.Sprintf("/tmp/%s", date[0:4])
}

func checkFolderDownloadImage(date string) {
	folder := getFolderDownloadImage(date)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, os.ModePerm)
	}
}

func fileNameDownloadCover(date string, movie movie) string {
	regex, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		panic(err)
	}
	movieTitle := regex.ReplaceAllString(movie.Title, "")
	return fmt.Sprintf("%s__%s__%s.jpg", date, movie.IMDb, movieTitle)
}

func downloadCover(movie movie, date string) {
	image, err := http.Get(movie.Cover)
	if err != nil {
		panic(err)
	}
	defer image.Body.Close()

	checkFolderDownloadImage(date)

	filename := fmt.Sprintf("%s/%s", getFolderDownloadImage(date), fileNameDownloadCover(date, movie))
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
