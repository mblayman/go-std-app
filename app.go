package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type Movie struct {
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseYear"`
}

func main() {
	fmt.Println("Starting server...")

	router := http.NewServeMux()
	router.HandleFunc("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s", html.EscapeString(r.PathValue("user")))
	})

	router.HandleFunc("GET /movies", func(w http.ResponseWriter, r *http.Request) {
		// movies := []Movie{
		// 	{Title: "The Shawshank Redemption", ReleaseYear: 1994},
		// 	{Title: "The Usual Suspects", ReleaseYear: 1995},
		// }
		movies := map[string]Movie{
			"a": {Title: "A movie", ReleaseYear: 2024},
			"b": {Title: "B movie", ReleaseYear: 2025},
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(movies)
	})

	v2 := http.NewServeMux()
	v2.HandleFunc("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello v2 %s", html.EscapeString(r.PathValue("user")))
	})

	router.Handle("/v2/", http.StripPrefix("/v2", v2))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Fatal(server.ListenAndServe())
}
