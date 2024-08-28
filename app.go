package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

type Movie struct {
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseYear"`
}

func main() {
	fmt.Println("Starting server...")

	router := http.NewServeMux()

	// Basic routing
	router.HandleFunc("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s", html.EscapeString(r.PathValue("user")))
	})

	v2 := http.NewServeMux()
	v2.HandleFunc("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello v2 %s", html.EscapeString(r.PathValue("user")))
	})

	router.Handle("/v2/", http.StripPrefix("/v2", v2))

	// JSON
	router.HandleFunc("GET /movies", func(w http.ResponseWriter, r *http.Request) {
		movies := []Movie{
			{Title: "The Shawshank Redemption", ReleaseYear: 1994},
			{Title: "The Usual Suspects", ReleaseYear: 1995},
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(movies)
	})

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	server.ListenAndServe()
}
