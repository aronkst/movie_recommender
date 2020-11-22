# Movie Recommender

This is a movie recommender project, which uses the [Go](https://github.com/golang/go) programming language and the [Gin](https://github.com/gin-gonic/gin) framework.

A Crawler was developed to search for information about the movies on IMDb.

The rating of the recommended movies works with a calculation that involves the number of votes for the movie on the IMDb, the movie's score on the IMDb and the movie's score on the Metascore, if this score is available on the movie's IMDb page.

The application uses the [SQLite](https://sqlite.org/) database to store information about the application's movies and features.

The watched movies are stored in folders, the name of the folder being the year in which the movie was watched, or 0000 if the date when the movie was watched was not registered. The reason for storing watched movies in folders is to have control and a quick view of all watched movies, without the need to run the application for this case.

This application works in conjunction with its [front-end](https://github.com/aronkst/movie_recommender_react), which was developed in [React](https://github.com/facebook/react).

## Run the application

Clone this repository and build the application with the command below.

`go build .`

Clone the [front-end](https://github.com/aronkst/movie_recommender_react) repository, and build with the command below.

`npm run build`

Now you must join the files from the [back-end](https://github.com/aronkst/movie_recommender) and [front-end](https://github.com/aronkst/movie_recommender_react) inside a single folder, and just run the build file from Go, with the command below.

`./movie_recommender`

With that the application will raise the server, which when accessing it, will load the files from the front-end and you will have access to the application.
