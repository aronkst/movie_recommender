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
		<style>
			h2 {
				margin: 10px 0;
			}
			hr {
				margin-bottom: 30px;
			}
			.hide {
				display: none;
			}
			.tableTd {
				font-weight: bold;
				padding: 0 10px 10px 0;
			}
			.tableTdValue {
				padding-bottom: 10px;
			}
			.imgRecommendedBy {
				margin: 5px 5px 0 0;
				width: 70px;
			}
		</style> 
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
				<p id="IMDb" class="hide">%s</p>
				<p id="Cover" class="hide">%s</p>
				<p id="CoverSmall" class="hide">%s</p>
				<p id="Points" class="hide">%d</p>
				<p id="RecommendedMovies" class="hide">%s</p>
				<p id="RecommendedBy" class="hide">%s</p>
				<p id="RecommendedByTitles" class="hide">%s</p>
				<img id="CoverLocal" src="./.covers/%s.jpg" />
				<a href="https://www.imdb.com/title/%s/" target="_blank">Go to IMDb</a>
				<h2 id="Title">%s</h2>
				<table>
					<tr>
						<td class="tableTd">Summary:</td>
						<td id="Summary" class="tableTdValue">%s</td>
					</tr>
					<tr>
						<td class="tableTd">Year:</td>
						<td id="Year" class="tableTdValue">%d</td>
					</tr>
					<tr>
						<td class="tableTd">IMDb Score:</td>
						<td id="Score" class="tableTdValue">%f</td>
					</tr>
					<tr>
						<td class="tableTd">IMDb Amount of Votes:</td>
						<td id="AmountOfVotes" class="tableTdValue">%d</td>
					</tr>
					<tr>
						<td class="tableTd">Metascore:</td>
						<td id="Metascore" class="tableTdValue">%d</td>
					</tr>
					<tr>
						<td class="tableTd">Genres:</td>
						<td id="Genres" class="tableTdValue">%s</td>
					</tr>
					<tr>
						<td class="tableTd">Recommended by:</td>
						<td></td>
					</tr>
				</table>
				<div>%s</div>
				<br />
			</div>%s`, movie.IMDb, movie.Cover, movie.CoverSmall, movie.Points, strings.Join(movie.RecommendedMovies, ", "),
		strings.Join(movie.RecommendedBy, ", "), strings.Join(movie.RecommendedByTitles, ", "), movie.IMDb, movie.IMDb,
		movie.Title, movie.Summary, movie.Year, movie.Score, movie.AmountOfVotes, movie.Metascore,
		strings.Join(movie.Genres, ", "), coversRecommendedBy(movie), "\n")
}

func coversRecommendedBy(movie movie) string {
	var imgCovers string
	for index, recommendedMovie := range movie.RecommendedBy {
		imgCovers += fmt.Sprintf(`<img src="./.covers/%s.jpg" title="%s" class="imgRecommendedBy" />`,
			recommendedMovie, movie.RecommendedByTitles[index])
	}
	return imgCovers
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
	watchedMoviesHTML, recommendedMovies := readWatchedAndRecommendedMoviesFromHTML()

	for _, imdb := range watchedMovies {
		if contains, _ := findMovieIMDb(watchedMoviesHTML, imdb); contains == false {
			movie := getMovie(imdb)
			newWatchedMovies = append(newWatchedMovies, movie)
		}
	}

	watchedMoviesHTML = append(watchedMoviesHTML, newWatchedMovies...)

	for _, movie := range newWatchedMovies {
		for _, recommendedMovieIMDb := range movie.RecommendedMovies {
			if contains, _ := findMovieIMDb(watchedMoviesHTML, recommendedMovieIMDb); contains == false {
				if contains, index := findMovieIMDb(recommendedMovies, recommendedMovieIMDb); contains {
					recommendedMovies[index].Points += movie.Points
					recommendedMovies[index].RecommendedBy = append(recommendedMovies[index].RecommendedBy, movie.IMDb)
					recommendedMovies[index].RecommendedByTitles = append(recommendedMovies[index].RecommendedByTitles, movie.Title)
				} else {
					recommendedMovie := getMovie(recommendedMovieIMDb)
					if validMovie(recommendedMovie) {
						recommendedMovie.RecommendedBy = []string{movie.IMDb}
						recommendedMovie.RecommendedByTitles = []string{movie.Title}
						recommendedMovies = append(recommendedMovies, recommendedMovie)
					}
				}
			}
		}
	}

	createHTML(watchedMoviesHTML, recommendedMovies)
}
