package main

import (
	"strings"
)

func listRecommendedMovies(offset int, title string, summary string, year int64, imdb string, genre string, score float64, metascore int64, order string) ([]movie, int64) {
	var watchedMoviesListIMDb []string
	var watchedMoviesAgainListIMDb []string
	var recommendedMoviesListIMDb []string
	var blockedMoviesListIMDb []string
	var listIMDb []string
	var movies []movie
	var blockedMovies []blockedMovie

	watchedMovies := getWatchedMoviesFromFolders(true)

	for _, watchedMovie := range watchedMovies {
		watchedMoviesListIMDb = append(watchedMoviesListIMDb, watchedMovie.IMDb)
	}
	watchedMoviesListIMDb = uniqueValuesInArrayString(watchedMoviesListIMDb)

	database.Where("imdb IN ?", watchedMoviesListIMDb).Find(&movies)

	for _, movie := range movies {
		recommendedMovieListIMDb := strings.Split(movie.RecommendedMovies, ",")
		for _, IMDb := range recommendedMovieListIMDb {
			recommendedMoviesListIMDb = append(recommendedMoviesListIMDb, IMDb)
		}
	}
	recommendedMoviesListIMDb = uniqueValuesInArrayString(recommendedMoviesListIMDb)
	listIMDb = removeItemInSliceIfExistInSlice(recommendedMoviesListIMDb, watchedMoviesListIMDb)

	database.Find(&blockedMovies)

	for _, blockedMovie := range blockedMovies {
		blockedMoviesListIMDb = append(blockedMoviesListIMDb, blockedMovie.IMDb)
	}
	listIMDb = removeItemInSliceIfExistInSlice(listIMDb, blockedMoviesListIMDb)

	watchedMovies = getWatchedMoviesFromFolders(false)

	for _, watchedMovie := range watchedMovies {
		watchedMoviesAgainListIMDb = append(watchedMoviesAgainListIMDb, watchedMovie.IMDb)
	}
	listIMDb = removeItemInSliceIfExistInSlice(listIMDb, watchedMoviesAgainListIMDb)

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
		query = query.Where("year >= ?", year)
	}

	if imdb != "" {
		query = query.Where("imdb = ?", imdb)
	}

	if genre != "" {
		genre = "%" + genre + "%"
		query = query.Where("genres LIKE ?", genre)
	}

	if score > 0 {
		query = query.Where("score >= ?", score)
	}

	if metascore > 0 {
		query = query.Where("metascore >= ?", metascore)
	}

	var count int64
	query.Model(&movie{}).Count(&count)
	pages := countPages(count)

	query.
		Order("points desc").
		Limit(10).
		Offset(offset).
		Find(&movies)

	return movies, pages
}
