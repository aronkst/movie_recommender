package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

func readWatchedMovies() []movie {
	var movies []movie

	folders, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, folder := range folders {
		if folder.IsDir() && folder.Name() != ".git" && folder.Name() != ".covers" && folder.Name() != ".data" && folder.Name() != "templates" {
			files, err := ioutil.ReadDir(fmt.Sprintf("./%s/", folder.Name()))
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				like := strings.Split(file.Name(), "__")[3]
				if like[0:1] == "1" {
					imdb := strings.Split(file.Name(), "__")[1]
					movie := getMovie(imdb)
					movies = append(movies, movie)
				}
			}
		}
	}

	movies = uniqueMovies(movies)
	return movies
}

func createRecommendedMovies() []movie {
	var recommendedMovies []movie

	watchedMovies := readWatchedMovies()
	for _, movie := range watchedMovies {
		for _, recommendedMovieIMDb := range movie.RecommendedMovies {
			if contains, _ := findMovieByIMDb(watchedMovies, recommendedMovieIMDb); contains == false {
				if contains, index := findMovieByIMDb(recommendedMovies, recommendedMovieIMDb); contains {
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

	return recommendedMovies
}

func recommendedMovies(page int64, title string, summary string, year int64, imdb string, genre string, score float64, metascore int64, order string) ([]movie, []int) {
	recommendedMovies := createRecommendedMovies()
	recommendedMovies = moviesSearch(recommendedMovies, title, summary, year, imdb, genre, score, metascore)
	recommendedMovies = sortMovies(recommendedMovies, order)
	pages := moviesPages(recommendedMovies)
	recommendedMovies = moviesPagination(recommendedMovies, page)
	return recommendedMovies, pages
}

func watchedMovies(page int64) ([]movie, []int) {
	watchedMovies := readWatchedMovies()
	return moviesPagination(watchedMovies, page), moviesPages(watchedMovies)
}

func uniqueMovies(array []movie) []movie {
	keys := make(map[string]bool)
	list := []movie{}
	for _, entry := range array {
		if _, value := keys[entry.IMDb]; !value {
			keys[entry.IMDb] = true
			list = append(list, entry)
		}
	}
	return list
}

func moviesPagination(data []movie, page int64) []movie {
	if page <= 0 {
		page = 1
	}

	start := (page - 1) * 10
	stop := start + 10
	count := int64(len(data))

	if start > count {
		return nil
	}

	if stop > count {
		stop = count
	}

	return data[start:stop]
}

func moviesHowManyPages(data []movie) int64 {
	howManyMovies := len(data)

	if howManyMovies%10 == 0 {
		return int64(howManyMovies / 10)
	}

	return int64((howManyMovies / 10) + 1)
}

func moviesPages(data []movie) []int {
	min := 1
	max := int(moviesHowManyPages(data))

	pages := make([]int, max-min+1)
	for i := range pages {
		pages[i] = min + i
	}
	return pages
}

func sortMovies(data []movie, order string) []movie {
	switch order {
	case "score":
		sort.Slice(data, func(i, j int) bool {
			return data[i].Score > data[j].Score
		})
	case "metascore":
		sort.Slice(data, func(i, j int) bool {
			return data[i].Metascore > data[j].Metascore
		})
	default:
		sort.Slice(data, func(i, j int) bool {
			return data[i].Points > data[j].Points
		})
	}

	return data
}

func findMovieByIMDb(movies []movie, imdb string) (bool, int) {
	for index, movie := range movies {
		if movie.IMDb == imdb {
			return true, index
		}
	}
	return false, -1
}

func validMovie(movie movie) bool {
	return movie.Cover != "" && movie.CoverSmall != "" && movie.Score > 0 &&
		movie.Year > 0 && movie.AmountOfVotes > 0 && len(movie.Genres) > 0 &&
		movie.Summary != "" && movie.Summary != `Add a Plot »`
}

func addNewMovie(imdb string, date string, like int64) {
	movie := getMovie(imdb)
	go downloadCover(movie, date, like)

	for _, recommendedMovie := range movie.RecommendedMovies {
		go getMovie(recommendedMovie)
	}
}

func moviesSearch(movies []movie, title string, summary string, year int64, imdb string, genre string, score float64, metascore int64) []movie {
	var searchedMovies []movie

	if title != "" || summary != "" || year > 0 || imdb != "" || genre != "" || score > 0 || metascore > 0 {
		for _, movie := range movies {
			if title != "" {
				lowerMovieTitle := strings.ToLower(movie.Title)
				lowerTitle := strings.ToLower(title)
				if strings.Contains(lowerMovieTitle, lowerTitle) == false {
					continue
				}
			}

			if summary != "" {
				lowerMovieSummary := strings.ToLower(movie.Summary)
				lowerSummary := strings.ToLower(summary)
				if strings.Contains(lowerMovieSummary, lowerSummary) == false {
					continue
				}
			}

			if year > 0 && movie.Year < year {
				continue
			}

			if imdb != "" && movie.IMDb != imdb {
				continue
			}

			if genre != "" && movieEqualGenre(genre, movie.Genres) == false {
				continue
			}

			if score > 0 && movie.Score < score {
				continue
			}

			if metascore == 1 && movie.Metascore <= 0 {
				continue
			}

			searchedMovies = append(searchedMovies, movie)
		}

		return searchedMovies
	}

	return movies
}

func movieEqualGenre(genre string, genres []string) bool {
	genre = strings.ToLower(genre)
	for _, g := range genres {
		movieGenre := strings.ToLower(g)
		if movieGenre == genre {
			return true
		}
	}
	return false
}
