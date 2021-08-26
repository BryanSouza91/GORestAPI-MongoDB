package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Database variable declaration
var (
	err error
)

// Parse the configuration file 'conf.json', and establish a connection to DB
func init() {
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal("error:", err)
	}
	decoder := json.NewDecoder(file)
	defer file.Close()
	err = decoder.Decode(&dao)
	if err != nil {
		log.Fatal("error:", err)
	}

	dao.Connection()
}

// Define HTTP routes
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", AllUsersEndpoint)
	mux.HandleFunc("/users/new", CreateUserEndpoint)
	mux.Handle("/users/update/", makeHandler(UpdateUserEndpoint))
	mux.Handle("/users/delete/", makeHandler(DeleteUserEndpoint))
	mux.Handle("/users/find/", makeHandler(FindUserEndpoint))
	if err = http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
