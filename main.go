package main

import (
	"encoding/json"
	"fmt"
	"github.com/AidanStrong/TicketSystem/db"
	"github.com/gorilla/mux" // import package using "go get github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Ticket - Our struct for all tickets, to be replaced with DB access
type Ticket struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Priority int    `json:"priority"`
	Status   string `json:"status"`
}

var Tickets []Ticket

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: homePage")
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func allTickets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: GET /tickets")
	conn := db.SetDbConn()
	tickets := db.GetAllTickets(conn)
	json.NewEncoder(w).Encode(tickets)
}

func getTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // gets path variables from request
	id, err := strconv.Atoi(vars["id"]) // accessed using tuple and converted to int
	if err != nil {
		panic(err)
	}
	fmt.Printf("Endpoint Hit: GET /ticket/%d \n", id)

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

func createTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Endpoint Hit: POST /ticket \n")
	reqBody, _ := ioutil.ReadAll(r.Body) // get the body of our POST request

	var ticket Ticket
	json.Unmarshal(reqBody, &ticket) // unmarshal json into ticket struct
	Tickets = append(Tickets, ticket)

	json.NewEncoder(w).Encode(ticket)
}

func deleteTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]) // accessed using tuple and converted to int
	if err != nil {
		panic(err)
	}
	fmt.Printf("Endpoint Hit: DELETE /ticket/%d \n", id)

	for index, ticket := range Tickets {
		if ticket.Id == id {
			Tickets = append(Tickets[:index], Tickets[index+1:]...)
		}
	}

	json.NewEncoder(w).Encode(Tickets)
}

func updateTicket(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body) // get the body of our POST request

	var updatedTicket Ticket
	json.Unmarshal(reqBody, &updatedTicket)
	fmt.Printf("Endpoint Hit: PATCH /ticket/%d \n", updatedTicket.Id)

	for index, ticket := range Tickets {
		if ticket.Id == updatedTicket.Id {
			ticket.Title = updatedTicket.Title
			ticket.Desc = updatedTicket.Desc
			ticket.Priority = updatedTicket.Priority

			Tickets = append(Tickets[:index], ticket)
			json.NewEncoder(w).Encode(ticket)
		}
	}
}

func handleRequests() {
	// creates a new instance of a mux router
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/tickets", allTickets)

	// NOTE: Ordering is important here! This has to be defined before
	// the other `/article` endpoint.
	router.HandleFunc("/ticket", createTicket).Methods("POST")
	router.HandleFunc("/ticket/{id}", deleteTicket).Methods("DELETE")
	router.HandleFunc("/ticket/{id}", updateTicket).Methods("PATCH")
	router.HandleFunc("/ticket/{id}", getTicket)

	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
	Created using tutorial: https://tutorialedge.net/golang/creating-restful-api-with-golang/
*/
func main() {
	//data base connection
	log.Println("Connecting to database...")
	myDb := db.SetDbConn()

	// create tables
	log.Println("Creating tables")
	db.CreateTables(myDb)
	defer myDb.Close()
	// handle requests
	log.Println("Server started on port :8080...")
	handleRequests()

}
