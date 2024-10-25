package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

// handleAddMovie handles HTTP POST requests to /movies. It takes a JSON
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

	fmt.Printf("Received movie: %+v\n", movie)

	db, err := openDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if movie.IMDb_id == "" || movie.Title == "" || movie.Year == 0 || movie.Rating == 0 {
		http.Error(w, "Missing required fields: IMDb ID, Title, Year, or Rating", http.StatusBadRequest)
		return
	}

	err = addMovie(db, movie.IMDb_id, movie.Title, movie.Year, movie.Rating)
	if err != nil {
		http.Error(w, "Could not add movie"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
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

	var filterYear int
	if year != "" {
		filterYear, _ = strconv.Atoi(year)
	}

	movies, err := listMovies(db, sortBy, order, filterYear)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	movie, err := showMovieDetails(db, imdbID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imdbID := vars["imdbID"]

	db, err := openDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = deleteMovie(db, imdbID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
