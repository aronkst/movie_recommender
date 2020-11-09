package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type watchedMovie struct {
	Date  int64
	IMDb  string
	Title string
	Like  int64
}

func getWatchedMoviesFromFolders() []watchedMovie {
	var watchedMovies []watchedMovie

	folders, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, folder := range folders {
		if folder.IsDir() && stringIsNumeric(folder.Name()) {
			files, err := ioutil.ReadDir(fmt.Sprintf("./%s/", folder.Name()))
			if err != nil {
				panic(err)
			}

			for _, file := range files {
				values := strings.Split(file.Name(), "__")

				likeString := values[3]
				likeString = likeString[:len(likeString)-4]
				like := stringToInt(likeString)

				if like == 1 {
					watchedMovie := watchedMovie{
						Date:  stringToInt(values[0]),
						IMDb:  values[1],
						Title: values[2],
						Like:  like,
					}

					watchedMovies = append(watchedMovies, watchedMovie)
				}
			}
		}
	}

	watchedMovies = uniqueWatchedMovies(watchedMovies)
	return watchedMovies
}

func uniqueWatchedMovies(array []watchedMovie) []watchedMovie {
	keys := make(map[string]bool)
	list := []watchedMovie{}
	for _, entry := range array {
		if _, value := keys[entry.IMDb]; !value {
			keys[entry.IMDb] = true
			list = append(list, entry)
		}
	}
	return list
}

func listWatchedMovies() []movie {
	var listIMDb []string
	var movies []movie

	watchedMovies := getWatchedMoviesFromFolders()
	for _, watchedMovie := range watchedMovies {
		listIMDb = append(listIMDb, watchedMovie.IMDb)
	}

	// TODO Paginação
	database.
		Where("imdb IN ?", listIMDb).
		Order("points desc").
		Limit(20).
		Offset(0).
		Find(&movies)

	return movies
}

func addWatchedMovies(imdb string, date string, like int64) map[string]string {
	json := make(map[string]string)

	movie, err := getMovie(imdb)
	if err != nil {
		json["error"] = "invalid movie"
		return json
	}

	downloadCover(movie, date, like)

	recommendedMovies := strings.Split(movie.RecommendedMovies, ",")
	for _, recommendedMovie := range recommendedMovies {
		getMovie(recommendedMovie)
	}

	json["imdb"] = movie.IMDb
	json["title"] = movie.Title

	return json
}
