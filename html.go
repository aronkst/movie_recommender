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

func createHTML(watchedMovies []movie, recommendedMovies []movie) {
	checkHTML()
	bytes := []byte(textHTML(watchedMovies, recommendedMovies))
	err := ioutil.WriteFile("Recommended Movies.html", bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func textHTML(watchedMovies []movie, recommendedMovies []movie) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
	<head>
		<title>Recommended Movies</title>
	</head>
	<body>
		<h1>Recommended Movies</h1>
		<div class="watched-movies" style="display: none;">
%s
		</div>
		<div class="recommended-movies">
%s
		</div>
	</body>
</html>`, textHTMLMovies(watchedMovies), textHTMLMovies(recommendedMovies))
}

func textHTMLRecommendedMovie(movie movie) string {
	return fmt.Sprintf(`			<div class="movie">
				<hr />
				<p class="Cover>%s</p>
				<img class="CoverSmall" src="%s" />
				<p class="IMDb">%s</p>
				<h2 class="Title">%s</h2>
				<p class="Year">%d</p>
				<p class="Summary">%s</p>
				<p class="Score">%f</p>
				<p class="AmountOfVotes">%d</p>
				<p class="Metascore">%d</p>
				<p class="Points">%d</p>
				<p class="Genres">%s</p>
				<p class="RecommendedMovies">%s</p>
			</div>%s`, movie.Cover, movie.CoverSmall, movie.IMDb, movie.Title,
		movie.Year, movie.Summary, movie.Score, movie.AmountOfVotes,
		movie.Metascore, movie.Points, strings.Join(movie.Genres, ", "),
		strings.Join(movie.RecommendedMovies, ", "), "\n")
}

func textHTMLMovies(movies []movie) string {
	var html string

	for _, movie := range movies {
		html += textHTMLRecommendedMovie(movie)
	}

	return html[0 : len(html)-1]
}
