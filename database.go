package main

import "gorm.io/gorm"

type movie struct {
	gorm.Model
	IMDb              string  `gorm:"column:imdb;uniqueIndex"`
	Title             string  `gorm:"index"`
	Year              int64   `gorm:"index"`
	Summary           string  `gorm:"index"`
	Score             float64 `gorm:"index"`
	AmountOfVotes     int64
	Metascore         int64  `gorm:"index"`
	Points            int64  `gorm:"index"`
	Genres            string `gorm:"index"`
	RecommendedMovies string
	URLCover          string
	URLCoverSmall     string
	Cover             string
}

type notWatch struct {
	gorm.Model
	IMDb string `gorm:"column:imdb;uniqueIndex"`
}

type invalidMovie struct {
	gorm.Model
	IMDb string `gorm:"column:imdb;uniqueIndex"`
}
