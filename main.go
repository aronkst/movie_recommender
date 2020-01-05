package main

func main() {
	// menu()

	var newWatchedMovies []movie

	watchedMovies := readWatchedMovies()
	watchedMoviesHTML, recommendedMovies := readHTMLMovies()

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
					recommendedMovies[index].RecommendedBy = append(recommendedMovies[index].RecommendedBy, movie.Title)
				} else {
					recommendedMovie := getMovie(recommendedMovieIMDb)
					recommendedMovie.RecommendedBy = []string{movie.Title}
					recommendedMovies = append(recommendedMovies, recommendedMovie)
				}
			}
		}
	}

	createHTML(watchedMoviesHTML, recommendedMovies)
}
