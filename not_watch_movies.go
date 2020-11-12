package main

func addNotWatch(imdb string) map[string]string {
	json := make(map[string]string)

	var notWatchWhere notWatch
	where := database.Where("imdb = ?", imdb).First(&notWatchWhere)

	if where.RowsAffected >= 1 {
		json["error"] = "movie already saved"
		return json
	}

	notWatch := notWatch{
		IMDb: imdb,
	}
	result := database.Create(&notWatch)

	if result.Error != nil {
		json["error"] = "an error has occurred"
		return json
	}

	json["imdb"] = imdb
	return json
}

func removeNotWatch(imdb string) map[string]string {
	json := make(map[string]string)

	var notWatchWhere notWatch
	where := database.Where("imdb = ?", imdb).First(&notWatchWhere)

	if where.RowsAffected >= 1 {
		database.Unscoped().Where("imdb = ?", imdb).Delete(notWatch{})

		json["imdb"] = imdb
		return json
	}

	json["error"] = "movie not found"
	return json
}

func listNotWatch(offset int, title string, summary string, year int64, imdb string, genre string, score float64, metascore int64, order string) []movie {
	var notWatchMovies []notWatch
	var listIMDb []string
	var movies []movie

	database.Find(&notWatchMovies)
	for _, notWatchMovie := range notWatchMovies {
		listIMDb = append(listIMDb, notWatchMovie.IMDb)
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
		query = query.Where("year > ?", year)
	}

	if imdb != "" {
		query = query.Where("imdb = ?", imdb)
	}

	if genre != "" {
		genre = "%" + genre + "%"
		query = query.Where("genres LIKE ?", genre)
	}

	if score > 0 {
		query = query.Where("score > ?", score)
	}

	if metascore > 0 {
		query = query.Where("metascore > ?", metascore)
	}

	query.
		Order("points desc").
		Limit(10).
		Offset(offset).
		Find(&movies)

	return movies
}
