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

func listWatchedMovies(offset int, title string, summary string, year int64, imdb string, genre string, score float64, metascore int64, order string) []movie {
	var listIMDb []string
	var movies []movie

	watchedMovies := getWatchedMoviesFromFolders()

	for _, watchedMovie := range watchedMovies {
		listIMDb = append(listIMDb, watchedMovie.IMDb)
	}
	listIMDb = uniqueValuesInArrayString(listIMDb)

	query := database.Where("imdb IN ?", listIMDb)

	if title != "" {
		title = "%" + title + "%"
		query = query.Where("title LIKE ?", title)
	}

	if summary != "" {
		summary = "%" + summary + "%"
		query = query.Where("summary LIKE ?", summary)
	}

	if year > 0 {
		query = query.Where("year > ?", year)
	}

	if imdb != "" {
		query = query.Where("imdb = ?", imdb)
	}

	if genre != "" {
		genre = "%" + genre + "%"
		query = query.Where("genres LIKE ?", genre)
	}

	if score > 0 {
		query = query.Where("score > ?", score)
	}

	if metascore > 0 {
		query = query.Where("metascore > ?", metascore)
	}

	query.
		Order("points desc").
		Limit(10).
		Offset(offset).
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

	return json
}
