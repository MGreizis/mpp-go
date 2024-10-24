package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// main is the entry point of the movie application. It can be run in different
// modes by providing one of the following subcommands:
//
//	add: Adds a movie to the database. The --imdbid, --title, --year, and
//	     --rating flags are required.
//
//	list: Lists all movies in the database. The --sort, --order, and --year
//	      flags can be used to filter the results.
//
//	details: Shows the details of a movie with the given IMDb ID. The --imdbid
//	         flag is required.
//
//	delete: Deletes a movie with the given IMDb ID from the database. The
//	        --imdbid flag is required.
//
// If no subcommand is provided, the server is started.
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

	fetchPostersCommand := flag.NewFlagSet("posters", flag.ExitOnError)
	posterLimit := fetchPostersCommand.Int("limit", 10, "The maximum number of movies to fetch posters for")

	if len(arguments) == 0 {
		startServer()
		return
	}

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch arguments[0] {
	case "add":
		addCommand.Parse(arguments[1:])
		handleAddMovieCLI(db, *addImdbId, *addTitle, *addYear, *addImdbRating)

	case "posters":
		fetchPostersCommand.Parse(arguments[1:])
		handleFetchPostersCLI(db, *posterLimit)
		fmt.Println("Posters added")

	case "list":
		listCommand.Parse(arguments[1:])
		handleListMoviesCLI(db, *sortBy, *order, *filterYear)

	case "details":
		detailsCommand.Parse(arguments[1:])
		if *detailsImdbId == "" {
			fmt.Println("IMDb ID is required for 'details'")
			os.Exit(1)
		}
		handleShowDetailsCLI(db, *detailsImdbId)

	case "delete":
		deleteCommand.Parse(arguments[1:])
		if *deleteImdbId == "" {
			fmt.Println("IMDb ID is required")
			os.Exit(1)
		}
		handleDeleteMovieCLI(db, *deleteImdbId)

	default:
		fmt.Println("Expected 'add', 'list', 'details', 'delete' or 'posters' subcommands")
		os.Exit(1)
	}
}

// startServer starts the server on :8090 and registers the endpoints for
// adding a movie, listing movies, getting the details of a movie, and
// deleting a movie.
func startServer() {
	router := mux.NewRouter()
	router.Use(enableCORS)
	router.HandleFunc("/movies", handleAddMovie).Methods("POST")
	router.HandleFunc("/movies", handleListMovies).Methods("GET")
	router.HandleFunc("/movies/{imdbID}", handleMovieDetails).Methods("GET")
	router.HandleFunc("/movies/{imdbID}", handleDeleteMovie).Methods("DELETE")

	fmt.Println("Starting server on :8090")
	log.Fatal(http.ListenAndServe(":8090", router))
}
