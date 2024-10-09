package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./movies.db"

// openDB opens a connection to the database stored in the file at dbPath.
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// addMovie adds a movie to the movies table of the SQLite database at dbPath.
// It takes the IMDb ID, title, year, and rating of the movie as parameters.
func addMovie(db *sql.DB, imdbID, title string, year int, rating float64) error {
	stmt, err := db.Prepare("INSERT INTO movies (IMDb_id, Title, Year, Rating) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(imdbID, title, year, rating)
	return err
}

// listMovies executes a query against the movies table of the SQLite database
// at dbPath that lists movies according to the given parameters.
//
// It takes the following parameters:
//   - db: the SQLite database to query
//   - sortBy: the field to sort the movies by, either "year" or "rating"
//   - order: the order to sort the movies in, either "asc" or "desc"
//   - filterYear: the year to filter the movies by, or 0 to not filter
//
// It returns an error if the query cannot be executed.
func listMovies(db *sql.DB, sortBy string, order string, filterYear int) error {
	query := "SELECT Title, Year, Rating FROM movies"
	var args []interface{}

	if filterYear != 0 {
		query += " WHERE Year = ?"
		args = append(args, filterYear)
	}

	if filterYear == 0 {
		query = "SELECT Title, Year, Rating FROM movies"
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
		return err
	}
	defer rows.Close()

	// Print out the movie titles
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.Title, &movie.Year, &movie.Rating); err != nil {
			return err
		}
		fmt.Println(movie.Title)
	}
	return nil
}

// showMovieDetails shows the details of a movie with the given IMDb ID.
// It prints the IMDb id, title, rating, year, and poster of the movie.
// If the movie is not found, it prints "Movie not found" and returns nil.
// If there is another error, it returns the error.
func showMovieDetails(db *sql.DB, imdbID string) error {
	var movie Movie

	err := db.QueryRow("SELECT IMDb_id, Title, Rating, Year, Poster FROM movies WHERE IMDb_id = ?", imdbID).Scan(&movie.IMDb_id, &movie.Title, &movie.Rating, &movie.Year, &movie.Poster)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Movie not found")
			return nil
		}
		return err
	}
	fmt.Printf("IMDb id: %s\nTitle: %s\nRating: %.1f\nYear: %d\n", movie.IMDb_id, movie.Title, movie.Rating, movie.Year)
	if movie.Poster.Valid {
		fmt.Printf("Poster: %s\n", movie.Poster.String)
	} else {
		fmt.Println("Poster:") // yuck, but have to do this for tests to pass (hopefully)
	}

	return nil
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
		fmt.Println("No movie found with that IMDb ID")
	} else {
		fmt.Println("Movie deleted")
	}

	return nil
}

// !! --------------- REST API handlers ---------------

// handleAddMovie adds a movie to the database at dbPath. It takes a JSON
// object with the keys "IMDb_id", "Title", "Year", and "Rating" in the request
// body. If the object is invalid or the database cannot be opened, or if the
// movie cannot be added, it returns an HTTP error. Otherwise, it returns the
// movie with the HTTP status 201 Created.
func handleAddMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db, err := openDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = addMovie(db, movie.IMDb_id, movie.Title, movie.Year, movie.Rating)
	if err != nil {
		http.Error(w, "Could not add movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

// enableCORS is a middleware handler that enables CORS for the given next
// handler. It sets the following headers:
//
//   - Access-Control-Allow-Origin: *
//   - Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
//   - Access-Control-Allow-Headers: Content-Type
//
// It also handles preflight requests by responding with a 204 No Content
// status code when the request method is OPTIONS.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleListMovies handles the GET /movies endpoint. It returns the list of
// movies in the database, sorted by the given criteria. The query parameters
// are:
//
//   - `sort`: the column to sort by. Either "year" or "rating". If not
//     specified, defaults to "year".
//   - `order`: the order of the sort. Either "asc" or "desc". If not
//     specified, defaults to "asc".
//   - `year`: the year to filter movies by. If specified, only movies with
//     this year will be returned.
//
// If the database cannot be opened or the movies cannot be fetched, it
// returns an HTTP error. Otherwise, it returns the list of movies in JSON
// format with the HTTP status 200 OK.
func handleListMovies(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	queryParams := r.URL.Query()
	sortBy := queryParams.Get("sort")
	order := queryParams.Get("order")
	year := queryParams.Get("year")

	query := "SELECT IMDb_id, Title, Rating, Year, Poster FROM movies"
	var args []interface{}

	if year != "" {
		query += " WHERE Year = ?"
		args = append(args, year)
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
		http.Error(w, "Could not fetch movies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.IMDb_id, &movie.Title, &movie.Rating, &movie.Year, &movie.Poster); err != nil {
			// http.Error(w, "Could not fetch movies", http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		movies = append(movies, movie)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// handleMovieDetails handles HTTP GET requests to /movies/{imdbID}. It
// writes the movie with the given IMDb ID in JSON format to the
// http.ResponseWriter. If the database cannot be opened, or if the movie
// doesn't exist or cannot be fetched, it returns an HTTP error.
func handleMovieDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imdbID := vars["imdbID"]

	db, err := openDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var movie Movie
	err = db.QueryRow("SELECT IMDb_id, title, Rating, Year, Poster FROM movies WHERE IMDb_id = ?", imdbID).Scan(&movie.IMDb_id, &movie.Title, &movie.Rating, &movie.Year, &movie.Poster)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// handleDeleteMovie handles HTTP DELETE requests to /movies/{imdbID}. It
// deletes the movie with the given IMDb ID from the database. If the
// database cannot be opened, or if the movie doesn't exist or cannot be
// deleted, it returns an HTTP error. Otherwise, it returns the HTTP status
// 204 No Content.
func handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imdbID := vars["imdbID"]

	db, err := openDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM movies WHERE IMDb_id = ?", imdbID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		// http.Error(w, "Movie not found", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// main is the entry point of the program. It parses the command line
// arguments and executes the required subcommand. If no subcommand is
// provided, it starts the HTTP server.
func main() {
	arguments := os.Args[1:] // The first element is the path to the command, so we can skip that

	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	addImdbId := addCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")
	addTitle := addCommand.String("title", "Carmencita", "The movie's or series' title")
	addYear := addCommand.Int("year", 1894, "The movie's or series' year of release")
	addImdbRating := addCommand.Float64("rating", 5.7, "The movie's or series' rating on IMDb")

	detailsCommand := flag.NewFlagSet("details", flag.ExitOnError)
	detailsImdbId := detailsCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")

	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteImdbId := deleteCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")

	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	sortBy := listCommand.String("sort", "", "Sort movies by 'year' or 'rating'")
	order := listCommand.String("order", "", "Order 'asc' or 'desc'")
	filterYear := listCommand.Int("year", 0, "Filter movies by year")

	if len(arguments) == 0 {
		// Start server
		router := mux.NewRouter()
		router.Use(enableCORS)
		router.HandleFunc("/movies", handleAddMovie).Methods("POST")
		router.HandleFunc("/movies", handleListMovies).Methods("GET")
		router.HandleFunc("/movies/{imdbID}", handleMovieDetails).Methods("GET")
		router.HandleFunc("/movies/{imdbID}", handleDeleteMovie).Methods("DELETE")

		fmt.Println("Starting server on :8090")
		log.Fatal(http.ListenAndServe(":8090", router))
	}

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch arguments[0] {
	case "add":
		addCommand.Parse(arguments[1:])
		if *addImdbId == "" || *addTitle == "" || *addYear == 0 || *addImdbRating == 0 {
			fmt.Println("All fields (IMDb id, Title, Year, Rating) are required")
			os.Exit(1)
		}
		err := addMovie(db, *addImdbId, *addTitle, *addYear, *addImdbRating)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("IMDb id: %s\nTitle: %s\nRating: %.1f\nYear: %d\nPoster: null\n", *addImdbId, *addTitle, *addImdbRating, *addYear)

	case "list":
		listCommand.Parse(arguments[1:])
		err := listMovies(db, *sortBy, *order, *filterYear)
		if err != nil {
			log.Fatal(err)
		}

	case "details":
		detailsCommand.Parse(arguments[1:])
		if *detailsImdbId == "" {
			fmt.Println("IMDb ID is required")
			os.Exit(1)
		}
		err := showMovieDetails(db, *detailsImdbId)
		if err != nil {
			log.Fatal(err)
		}

	case "delete":
		deleteCommand.Parse(arguments[1:])
		if *deleteImdbId == "" {
			fmt.Println("IMDb ID is required")
			os.Exit(1)
		}
		err := deleteMovie(db, *deleteImdbId)
		if err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("Expected 'add', 'list', 'details' or 'delete' subcommands")
		os.Exit(1)
	}
}
