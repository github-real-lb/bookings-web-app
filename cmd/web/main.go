package main

import (
	"fmt"
	"net/http"

	"github.com/github-real-lb/go-web-app/pkg/handlers"
)

const addr = "localhost:8080"

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println("Starting web server on,", addr)
	http.ListenAndServe(addr, nil)
}
