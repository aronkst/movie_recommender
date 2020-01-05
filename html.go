package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
				<h2>%s (%d)</h2>
				<img src="%s">
				<p>%s</p>
				<p>%f</p>
			</div>%s`, movie.Title, movie.Year, movie.CoverSmall, movie.Summary, movie.Score, "\n")
}

func textHTMLRecommendedMovies(movies []movie) string {
	var html string

	for _, movie := range movies {
		html += textHTMLRecommendedMovie(movie)
	}

	return html[0 : len(html)-1]
}
