package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func readMovies() []string {
	var database []string

	folders, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, folder := range folders {
		if folder.IsDir() && folder.Name() != ".git" {
			files, err := ioutil.ReadDir(fmt.Sprintf("./%s/", folder.Name()))
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				imdb := strings.Split(file.Name(), "__")[1]
				database = append(database, imdb)
			}
		}
	}

	return uniqueArrayString(database)
}
