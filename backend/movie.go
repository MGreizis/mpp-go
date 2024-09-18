package main

import (
	"database/sql"
	"encoding/json"
)

type Movie struct {
	IMDb_id string
	Title   string
	Rating  float64
	Year    int
	Poster  NullString // Explanation below
}

// What follows here is an explanation of the custom NullString type you can
// see above. It is recommended to read and test it. But if you really want the
// TL;DR --> a NullString is a sql.NullString that looks nice when put in a
// JSON reponse.

// A sql.NullString is a string type that is nullable, i.e. when there is no
// value in the database. This type is marshalled into a JSON object like this
// when the database has stored `null`.
//
//	{
//	    String: "",
//	    Valid: false
//	}
//
// Or like this when the database has stored a value:
//
//	{
//	    String: "Your nice value here",
//	    Valid: true
//	}
//
// This is bad UX for the consumer of our APIs. Go can solve this multiple ways.
// I would recommend creating a new type that is the same as sql.NullString in
// every way but the (un)marshalling of JSON. This new NullString type shows as
// `null` in JSON when there is no value in the database, or as the actual value
// if there is one in the database.
type NullString struct {
	sql.NullString
}

// Show the string directly as value if there is one, otherwise show `null`
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// Unwrap a value into the original sql.NullString type.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}
