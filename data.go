package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func readWatchedMovies() []movie {
	var movies []movie

	folders, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, folder := range folders {
		if folder.IsDir() && folder.Name() != ".git" && folder.Name() != ".covers" && folder.Name() != ".data" {
			files, err := ioutil.ReadDir(fmt.Sprintf("./%s/", folder.Name()))
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				like := strings.Split(file.Name(), "__")[3]
				if like[0:1] == "1" {
					imdb := strings.Split(file.Name(), "__")[1]
					movie := getMovie(imdb)
					movies = append(movies, movie)
				}
			}
		}
	}

	return uniqueMovies(movies)
}

func uniqueMovies(array []movie) []movie {
	keys := make(map[string]bool)
	list := []movie{}
	for _, entry := range array {
		if _, value := keys[entry.IMDb]; !value {
			keys[entry.IMDb] = true
			list = append(list, entry)
		}
	}
	return list
}
