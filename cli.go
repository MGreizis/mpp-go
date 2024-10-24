package main

import (
	"database/sql"
	"fmt"
)

// handleAddMovieCLI adds a movie to the database with the given
// IMDb ID, title, year, and rating. If the movie could not be added
// to the database, it prints an error message. Otherwise, it prints
// the movie details that were added.
func handleAddMovieCLI(db *sql.DB, imdbID, title string, year int, rating float64) {
	err := addMovie(db, imdbID, title, year, rating)
	if err != nil {
		fmt.Println("Error adding movie:", err)
		return
	}
	fmt.Printf("IMDb id: %s\nTitle: %s\nRating: %.1f\nYear: %d\nPoster: null\n", imdbID, title, rating, year)
}

// handleListMoviesCLI lists all movies in the database. It takes the
// database connection, the column to sort by, the order to sort in, and
// the year to filter by as parameters. If there is an error listing the
// movies, it prints an error message. Otherwise, it prints the title of
// each movie.
func handleListMoviesCLI(db *sql.DB, sortBy, order string, filterYear int) {
	movies, err := listMovies(db, sortBy, order, filterYear)
	if err != nil {
		fmt.Println("Error listing movies:", err)
		return
	}
	for _, movie := range movies {
		fmt.Println(movie.Title)
	}
}

// handleShowDetailsCLI shows the details of a movie with the given
// IMDb ID in the database. If the movie cannot be found, it prints
// an error message. Otherwise, it prints the movie details in a
// human-readable format.
func handleShowDetailsCLI(db *sql.DB, imdbID string) {
	movie, err := showMovieDetails(db, imdbID)
	if err != nil {
		fmt.Println("Error showing movie details:", err)
		return
	}
	fmt.Printf("IMDb id: %s\nTitle: %s\nRating: %.1f\nYear: %d\nPoster: %s\n", movie.IMDb_id, movie.Title, movie.Rating, movie.Year, movie.Poster.String)
}

// handleDeleteMovieCLI deletes the movie with the given IMDb ID from the database.
// If the movie cannot be found in the database, it prints an error message.
// Otherwise, it prints a success message.
func handleDeleteMovieCLI(db *sql.DB, imdbID string) {
	err := deleteMovie(db, imdbID)
	if err != nil {
		fmt.Println("Error deleting movie:", err)
		return
	}
	fmt.Println("Movie deleted")
}

// handleFetchPostersCLI fetches movie posters concurrently for a specified
// number of movies without posters in the database. It takes a database
// connection and a limit on the number of movies to process. If an error
// occurs during the fetching process, it prints an error message.
func handleFetchPostersCLI(db *sql.DB, limit int) {
	const workerCount = 3
	err := fetchPostersConcurrently(db, workerCount, limit)
	if err != nil {
		fmt.Println("Error fetching posters:", err)
		return
	}
}

func handleUpdatePosterCLI(db *sql.DB, imdbID, posterURL string) {
	err := updateMoviePoster(db, imdbID, posterURL)
	if err != nil {
		fmt.Println("Error updating movie poster:", err)
		return
	}
	fmt.Println("Movie poster updated")
}
