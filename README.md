# Movie Recommender

This is a movie recommendation project, using the Go programming language.

A Crawler is used to fetch information from movies on IMDb.

The rating of the recommended films works with a calculation involving the number of votes for the film on the IMDb, the film's rating on the IMDb and the film's rating on the Metascore, if this rating is on the IMDb film page.

## Run the application

To run the application use the command:

`go run .`

If you prefer, build the application with the command below:

`go build .`

Then to run the application, use the following command:

`./movie_recommender`