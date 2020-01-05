package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func checkHTML() {
	file := "Recommended Movies.html"
	if _, err := os.Stat(file); os.IsNotExist(err) == false {
		err := os.Remove(file)
		if err != nil {
			panic(err)
		}
	}
}

func createHTML(watchedMovies []string, recommendedMovies []movie) {
	checkHTML()
	bytes := []byte(textHTML(watchedMovies, recommendedMovies))
	err := ioutil.WriteFile("Recommended Movies.html", bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func textHTML(watchedMovies []string, recommendedMovies []movie) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
	<head>
		<title>Recommended Movies</title>
	</head>
	<body>
		<div class="watched-movies">
			%s
		</div>
		<div class="recommended-movies">
%s
		</div>
	</body>
</html>`, textHTMLWatchedMovies(watchedMovies), textHTMLRecommendedMovies(recommendedMovies))
}

func textHTMLWatchedMovies(movies []string) string {
	var html string

	for _, movie := range movies {
		html += fmt.Sprintf("%s, ", movie)
	}

	return html[0 : len(html)-2]
}

func textHTMLRecommendedMovie(movie movie) string {
	return fmt.Sprintf(`			<div class="movie">
				<hr />
				<img class="Cover" src="%s">
				<img class="CoverSmall" src="%s">
				<p class="IMDb">%s</p>
				<h2 class="Title">%s</h2>
				<p class="Year">%d</p>
				<p class="Summary">%s</p>
				<p class="Score">%f</p>
				<p class="AmountOfVotes">%d</p>
				<p class="Metascore">%d</p>
				<p class="Genres">%s</p>
				<p class="RecommendedMovies">%s</p>
			</div>%s`, movie.Cover, movie.CoverSmall, movie.IMDb, movie.Title,
		movie.Year, movie.Summary, movie.Score, movie.AmountOfVotes,
		movie.Metascore, strings.Join(movie.Genres, ", "),
		strings.Join(movie.RecommendedMovies, ", "), "\n")
}

func textHTMLRecommendedMovies(movies []movie) string {
	var html string

	for _, movie := range movies {
		html += textHTMLRecommendedMovie(movie)
	}

	return html[0 : len(html)-1]
}
