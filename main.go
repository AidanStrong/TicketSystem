package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux" // import package using "go get github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Ticket - Our struct for all tickets, to be replaced with DB access
type Ticket struct {
	Id       int    `json:"Id"`
	Title    string `json:"Title"`
	Desc     string `json:"desc"`
	Priority int    `json:"priority"`
}

var Tickets []Ticket

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func allTickets(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Tickets)
	fmt.Println("Endpoint Hit: tickets")
}

func getTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // gets path variables from request
	id, err := strconv.Atoi(vars["id"]) // accessed using tuple and converted to int
	if err != nil {
		panic(err)
	}

	found := false
	for _, ticket := range Tickets {
		if ticket.Id == id {
			found = true
			json.NewEncoder(w).Encode(ticket)
		}
	}

	if !found {
		fmt.Fprintf(w, "No Ticket found for id %d!", id)
	}
}

func handleRequests() {
	// creates a new instance of a mux router
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/tickets", allTickets)
	router.HandleFunc("/ticket/{id}", getTicket)

	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
	Created using tutorial: https://tutorialedge.net/golang/creating-restful-api-with-golang/
*/
func main() {
	// Test Data
	Tickets = []Ticket{
		Ticket{Id: 1, Title: "Hello", Desc: "Ticket Description", Priority: 1},
		Ticket{Id: 2, Title: "Hello 2", Desc: "Ticket Description", Priority: 2},
	}

	log.Println("Server started on port :8080...")
	log.Println("change")
	handleRequests()
}
