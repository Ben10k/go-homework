package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/measure", measureHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
