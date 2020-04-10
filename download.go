package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func getFolderDownloadCover(date string) string {
	return fmt.Sprintf("./%s", date[0:4])
}

func checkFolderDownloadCover(date string) {
	folder := getFolderDownloadCover(date)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, os.ModePerm)
	}
}

func checkFolderDownloadSmallCover() {
	folder := "./covers"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, os.ModePerm)
	}
}

func fileNameDownloadCover(date string, movie movie, like int64) string {
	regex, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		panic(err)
	}
	movieTitle := regex.ReplaceAllString(movie.Title, "")
	return fmt.Sprintf("%s__%s__%s__%d.jpg", date, movie.IMDb, movieTitle, like)
}

func downloadCover(movie movie, date string, like int64) {
	checkFolderDownloadCover(date)

	filename := fmt.Sprintf("%s/%s", getFolderDownloadCover(date), fileNameDownloadCover(date, movie, like))

	downloadImage(movie.Cover, filename)
}

func downloadSmallCover(movie movie, date string, like int64) {
	checkFolderDownloadSmallCover()

	filename := fmt.Sprintf("./covers/%s.jpg", movie.IMDb)

	downloadImage(movie.CoverSmall, filename)
}

func downloadImage(url string, fileString string) {
	image, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer image.Body.Close()

	file, err := os.Create(fileString)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, image.Body)
	if err != nil {
		panic(err)
	}
}
