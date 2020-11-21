package main

func addBlockedMovie(imdb string) map[string]string {
	var blockedMovieWhere blockedMovie

	json := make(map[string]string)

	where := database.Where("imdb = ?", imdb).First(&blockedMovieWhere)

	if where.RowsAffected >= 1 {
		json["error"] = "movie already saved"
		return json
	}

	blockedMovie := blockedMovie{
		IMDb: imdb,
	}
	result := database.Create(&blockedMovie)

	if result.Error != nil {
		json["error"] = "an error has occurred"
		return json
	}

	json["imdb"] = imdb
	return json
}

func removeBlockedMovie(imdb string) map[string]string {
	json := make(map[string]string)

	var blockedMovieWhere blockedMovie
	where := database.Where("imdb = ?", imdb).First(&blockedMovieWhere)

	if where.RowsAffected >= 1 {
		database.Unscoped().Where("imdb = ?", imdb).Delete(blockedMovie{})

		json["imdb"] = imdb
		return json
	}

	json["error"] = "movie not found"
	return json
}

func listBlockedMovies(offset int, title string, summary string, year int64, imdb string, genre string, score float64, metascore int64, order string) ([]movie, int64) {
	var blockedMovies []blockedMovie
	var listIMDb []string
	var movies []movie

	database.Find(&blockedMovies)
	for _, blockedMovie := range blockedMovies {
		listIMDb = append(listIMDb, blockedMovie.IMDb)
	}

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
