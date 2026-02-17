package main

import (
	"fmt",
	"log",
	"net/http"
)

func main() {
	handleRequests()
}


func testController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8081", nil))
}