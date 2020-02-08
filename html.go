package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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

	sort.Slice(watchedMovies, func(i, j int) bool {
		return watchedMovies[i].Points > watchedMovies[j].Points
	})
	sort.Slice(recommendedMovies, func(i, j int) bool {
		return recommendedMovies[i].Points > recommendedMovies[j].Points
	})

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
		<!-- <div class="watched-movies">
%s
		</div> -->
		<div class="recommended-movies">
%s
		</div>
	</body>
</html>`, textHTMLMovies(watchedMovies), textHTMLMovies(recommendedMovies))
}

func textHTMLRecommendedMovie(movie movie) string {
	return fmt.Sprintf(`			<div class="movie">
				<hr />
				<p class="Cover">%s</p>
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
				<p class="RecommendedBy">%s</p>
			</div>%s`, movie.Cover, movie.CoverSmall, movie.IMDb, movie.Title,
		movie.Year, movie.Summary, movie.Score, movie.AmountOfVotes,
		movie.Metascore, movie.Points, strings.Join(movie.Genres, ", "),
		strings.Join(movie.RecommendedMovies, ", "),
		strings.Join(movie.RecommendedBy, ", "), "\n")
}

func textHTMLMovies(movies []movie) string {
	var html string

	for _, movie := range movies {
		html += textHTMLRecommendedMovie(movie)
	}

	if len(html) <= 0 {
		return ""
	}

	return html[0 : len(html)-1]
}
