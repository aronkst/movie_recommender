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
		<!-- <div id="WatchedMovies">
%s
		</div> -->
		<div id="RecommendedMovies">
%s
		</div>
	</body>
</html>`, textHTMLMovies(watchedMovies), textHTMLMovies(recommendedMovies))
}

func textHTMLStructure(movie movie) string {
	return fmt.Sprintf(`			<div id="Movie">
				<hr />
				<p id="Cover" style="display: none;" title="Cover">%s</p>
				<p id="CoverSmall" style="display: none;" title="Small Cover">%s</p>
				<p id="RecommendedMovies" style="display: none;" title="Recommended Movies">%s</p>
				<p id="IMDb" title="IMDb">%s</p>
				<img id="CoverLocal" src="./.covers/%s.jpg" title="Cover" />
				<h2 id="Title" title="Title">%s</h2>
				<p id="Summary" title="Summary">%s</p>
				<p id="Year" title="Year">%d</p>
				<p id="Score" title="Score">%f</p>
				<p id="AmountOfVotes" title="Amount Of Votes">%d</p>
				<p id="Metascore" title="Metascore">%d</p>
				<p id="Points" title="Points">%d</p>
				<p id="Genres" title="Genres">%s</p>
				<p id="RecommendedBy" title="Recommended By">%s</p>
			</div>%s`, movie.Cover, movie.CoverSmall, strings.Join(movie.RecommendedMovies, ", "),
		movie.IMDb, movie.IMDb, movie.Title, movie.Summary, movie.Year,
		movie.Score, movie.AmountOfVotes, movie.Metascore, movie.Points,
		strings.Join(movie.Genres, ", "), strings.Join(movie.RecommendedBy, ", "),
		"\n")
}

func textHTMLMovies(movies []movie) string {
	var html string

	for _, movie := range movies {
		html += textHTMLStructure(movie)
	}

	if len(html) <= 0 {
		return ""
	}

	return html[0 : len(html)-1]
}

func makeHTML() {
	var newWatchedMovies []movie

	watchedMovies := readWatchedMovies()
	watchedMoviesHTML, recommendedMovies := readHTMLMovies()

	for _, imdb := range watchedMovies {
		if contains, _ := findMovieIMDb(watchedMoviesHTML, imdb); contains == false {
			movie := getMovie(imdb)
			downloadSmallCover(movie)
			newWatchedMovies = append(newWatchedMovies, movie)
		}
	}

	watchedMoviesHTML = append(watchedMoviesHTML, newWatchedMovies...)

	for _, movie := range newWatchedMovies {
		for _, recommendedMovieIMDb := range movie.RecommendedMovies {
			if contains, _ := findMovieIMDb(watchedMoviesHTML, recommendedMovieIMDb); contains == false {
				if contains, index := findMovieIMDb(recommendedMovies, recommendedMovieIMDb); contains {
					recommendedMovies[index].Points += movie.Points
					recommendedMovies[index].RecommendedBy = append(recommendedMovies[index].RecommendedBy, movie.Title)
				} else {
					recommendedMovie := getMovie(recommendedMovieIMDb)
					downloadSmallCover(recommendedMovie)
					recommendedMovie.RecommendedBy = []string{movie.Title}
					recommendedMovies = append(recommendedMovies, recommendedMovie)
				}
			}
		}
	}

	createHTML(watchedMoviesHTML, recommendedMovies)
}
