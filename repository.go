package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type OMDBAPIResponse struct {
	Poster string `json:"Poster"`
}

const dbPath = "./movies.db"
const OMDB_API_KEY = "34e1747c"

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

// fetchPoster retrieves the poster URL for a movie from the OMDB API based on the provided IMDb ID.
// It constructs the API request URL using the IMDb ID and the OMDB API key, and sends an HTTP GET request.
// If the request fails or the response cannot be decoded, it returns an error.
// If the poster is found, it returns the poster URL as a string.
// If no poster is found, it returns an error indicating that the poster was not found.
func fetchPoster(imdbID string) (string, error) {
	apiURL := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", imdbID, OMDB_API_KEY)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result OMDBAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Poster == "" || strings.ToLower(result.Poster) == "n/a" {
		return "", fmt.Errorf("poster not found")
	}
	return result.Poster, nil
}

// updatePosterInDB updates the poster URL for the movie with the given IMDb ID in the database.
// It takes a database connection, the IMDb ID of the movie, and the poster URL as parameters.
// It returns an error if the database cannot be updated.
func updatePosterInDB(db *sql.DB, imdbID, posterURL string) error {
	stmt, err := db.Prepare("UPDATE movies SET Poster = ? WHERE IMDb_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(posterURL, imdbID)
	return err
}

// fetchPostersConcurrently fetches movie posters concurrently for movies without posters in the database.
// It retrieves IMDb IDs of movies with missing posters from the database and utilizes multiple goroutines
// (as specified by workerCount) to fetch and update posters concurrently. Each worker fetches the poster
// URL from the OMDB API using the IMDb ID and updates the database with the poster URL if found.
// If any error occurs while fetching the poster or updating the database, it logs the error and continues with other IDs.
// Returns an error if there is a problem querying the database for IMDb IDs.
func fetchPostersConcurrently(db *sql.DB, workerCount int, limit int) error {
	rows, err := db.Query("SELECT IMDb_id FROM movies WHERE Poster IS NULL LIMIT ?", limit) // If there are DB issues check this line and replace IS NULL with ''
	if err != nil {
		return err
	}
	defer rows.Close()

	imdbIDs := []string{}
	for rows.Next() {
		var imdbID string
		if err := rows.Scan(&imdbID); err != nil {
			return err
		}
		imdbIDs = append(imdbIDs, imdbID)
	}

	var wg sync.WaitGroup
	imdbChan := make(chan string, len(imdbIDs))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for imdbID := range imdbChan {
				posterURL, err := fetchPoster(imdbID)
				if err != nil {
					continue
				}

				if err := updatePosterInDB(db, imdbID, posterURL); err != nil {
					// continue
				}
			}
		}(i + 1)
	}

	for _, imdbID := range imdbIDs {
		imdbChan <- imdbID
	}
	close(imdbChan)

	wg.Wait()
	return nil
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
	query := "SELECT IMDb_id, Title, Year, Rating, Poster FROM movies"
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
		if err := rows.Scan(&movie.IMDb_id, &movie.Title, &movie.Year, &movie.Rating, &movie.Poster); err != nil {
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
	imdbID = strings.TrimSpace(imdbID)

	var movie Movie
	query := "SELECT IMDb_id, Title, Rating, Year, Poster FROM movies WHERE IMDb_id = ?"

	err := db.QueryRow(query, imdbID).Scan(&movie.IMDb_id, &movie.Title, &movie.Rating, &movie.Year, &movie.Poster)
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
