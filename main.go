package main

import (
	"flag"
	"os"
)

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

	switch arguments[0] {
	case "add":
		addCommand.Parse(arguments[1:])
		// :)
	case "list":
		// :)
	case "details":
		detailsCommand.Parse(arguments[1:])
		// :)
	case "delete":
		deleteCommand.Parse(arguments[1:])
		// :)
	}
}
