package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Movie struct {
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseYear"`
}

//go:embed static/*
var content embed.FS

type contextKey string

const userKey contextKey = "user"

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pull something from the request, likely a cookie.
		// Check that thing against the db.

		// For this, we will pretend that we did the
		// stuff above and have a valid user.
		ctx := context.WithValue(r.Context(), userKey, "42")

		fmt.Println("Called the auth middleware")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// adminMiddleware pretends to permit admin users. In this case, admin has an
// ID of 42.
func adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(userKey).(string)
		if userID != "42" {
			http.Error(w, "Nope", http.StatusForbidden)
			return
		}
		fmt.Println("Called the admin middleware")
		next.ServeHTTP(w, r)
	})
}

func createDb() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	insertUserSQL := `INSERT INTO users (name, age) VALUES (?, ?)`
	_, err = db.Exec(insertUserSQL, "Alice", 30)
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}

	return db
}

func main() {
	db := createDb()
	defer db.Close()

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

	// Templates
	router.HandleFunc("GET /template", func(w http.ResponseWriter, r *http.Request) {
		base := `<!DOCTYPE html>
<html><head><title>testing</title></head>
<body>{{ block "content" . }}{{end}}</body></html>`

		content := `
{{ define "content" }}
<h1>{{.Title}}</h1>
<p>{{.ReleaseYear}}</p>
{{ end }}`

		movie := Movie{Title: "The <Shawshank> Redemption", ReleaseYear: 1994}
		tmpl := template.Must(template.New("base").Parse(base))
		tmpl = template.Must(tmpl.Parse(content))
		tmpl.Execute(w, movie)
	})

	// Static files
	router.Handle("/static/", http.FileServer(http.FS(content)))

	// Middleware
	admin := http.NewServeMux()
	admin.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("A secret to everybody but 42"))
	})

	router.Handle("/admin/", http.StripPrefix("/admin", adminMiddleware(admin)))

	// Database
	router.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		queryUserSQL := `SELECT id, name, age FROM users`
		rows, err := db.Query(queryUserSQL)
		if err != nil {
			log.Fatalf("Failed to query data: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var name string
			var age int
			err = rows.Scan(&id, &name, &age)
			if err != nil {
				log.Fatalf("Failed to scan row: %v", err)
			}
			fmt.Fprintf(w, "ID: %d, Name: %s, Age: %d\n", id, name, age)
		}

	})

	server := http.Server{
		Addr:    ":8000",
		Handler: authMiddleware(router),
	}
	server.ListenAndServe()
}
