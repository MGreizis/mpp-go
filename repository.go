package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./movies.db"

// openDB opens a database connection to the SQLite database file at
// dbPath. It returns the database connection and an error if the database
// cannot be opened.
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// movie's IMDb ID, title, year of release, and rating as parameters. It returns
// an error if the movie cannot be added to the database.
func addMovie(db *sql.DB, imdbID, title string, year int, rating float64) error {
	stmt, err := db.Prepare("INSERT INTO movies (IMDb_id, Title, Year, Rating) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(imdbID, title, year, rating)
	return err
}

// listMovies retrieves a list of movies from the database. It takes the database connection, an optional field to sort by, an optional order, and an optional year to filter by. It returns the list of movies and an error. If the database cannot be opened or the movies cannot be fetched, it returns an HTTP error. Otherwise, it returns the list of movies.
func listMovies(db *sql.DB, sortBy string, order string, filterYear int) ([]Movie, error) {
	query := "SELECT IMDb_id, Title, Year, Rating FROM movies"
	var args []interface{}

	if filterYear != 0 {
		query += " WHERE Year = ?"
		args = append(args, filterYear)
	}

	if sortBy != "" {
		if sortBy == "year" {
			query += " ORDER BY Year"
		} else if sortBy == "rating" {
			query += " ORDER BY Rating"
		}

		if order == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.IMDb_id, &movie.Title, &movie.Year, &movie.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

// showMovieDetails retrieves the details of a movie from the database based on the IMDb ID.
// It takes the database connection and the IMDb ID of the movie as parameters.
// It returns the movie details and an error. If the movie is not found, it returns an error indicating that the movie was not found.
func showMovieDetails(db *sql.DB, imdbID string) (Movie, error) {
	var movie Movie
	err := db.QueryRow("SELECT IMDb_id, Title, Rating, Year, Poster FROM movies WHERE IMDb_id = ?", imdbID).Scan(&movie.IMDb_id, &movie.Title, &movie.Rating, &movie.Year, &movie.Poster)
	if err != nil {
		if err == sql.ErrNoRows {
			return movie, fmt.Errorf("Movie not found")
		}
		return movie, err
	}
	return movie, nil
}

// deleteMovie deletes the movie with the given IMDb ID from the database.
// It prints a success message if the movie was deleted, or a message
// indicating that the movie was not found if it wasn't. If there is another
// error, it returns the error.
func deleteMovie(db *sql.DB, imdbID string) error {
	stmt, err := db.Prepare("DELETE FROM movies WHERE IMDb_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(imdbID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Movie not found")
	}

	return nil
}
