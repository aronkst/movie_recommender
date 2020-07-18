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

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Points > movies[j].Points
	})

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

	sort.Slice(recommendedMovies, func(i, j int) bool {
		return recommendedMovies[i].Points > recommendedMovies[j].Points
	})

	return recommendedMovies
}

func recommendedMovies(page int64) ([]movie, []int) {
	recommendedMovies := createRecommendedMovies()
	return moviesPagination(recommendedMovies, page), moviesPages(recommendedMovies)
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
