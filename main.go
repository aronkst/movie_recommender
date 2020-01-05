package main

func main() {
	// menu()

	watchedMovies := readWatchedMovies()
	watchedMoviesHTML, recommendedMovies := readHTMLMovies()

	for _, imdb := range watchedMovies {
		if contains, _ := findMovieIMDb(watchedMoviesHTML, imdb); contains == false {
			movie := getMovie(imdb)
			watchedMoviesHTML = append(watchedMoviesHTML, movie)

			for _, recommendedMovieIMDb := range movie.RecommendedMovies {
				if contains, index := findMovieIMDb(recommendedMovies, recommendedMovieIMDb); contains {
					recommendedMovies[index].Points += movie.Points
				} else {
					recommendedMovie := getMovie(recommendedMovieIMDb)
					recommendedMovies = append(recommendedMovies, recommendedMovie)
				}
			}
		}
	}

	createHTML(watchedMoviesHTML, recommendedMovies)
}
