package main

import (
	"strings"
)

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
