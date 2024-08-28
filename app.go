package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")

	router := http.NewServeMux()

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	server.ListenAndServe()
}
