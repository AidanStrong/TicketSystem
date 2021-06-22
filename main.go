package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux" // import package using "go get github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/AidanStrong/TicketSystem/db"
)

// database connection
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "ticketdb"
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
	json.NewEncoder(w).Encode(Tickets)
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

func createTables(db *sql.DB) {
	var sql string

	// drop tables
	sql = `DROP TABLE IF EXISTS tickets`
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}

	sql = `DROP TABLE IF EXISTS users`
	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	// create user table
	sql = `
	 	CREATE TABLE users (
	 	id SERIAL PRIMARY KEY,
	 	first_name TEXT,
	 	last_name TEXT,
		email TEXT UNIQUE NOT NULL)`
	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	// populate user table
	sql = `
	 	INSERT INTO users (first_name, last_name, email)
	 	VALUES ('Aidan', 'Strong', 'aidans@email.com'),
			   ('Alex', 'Smith', 'alecs@email.com'),
			   ('John', 'Smith', 'johns@email.com')`
	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	// create ticket table
	sql = `
	 	CREATE TABLE tickets (
	 	id				SERIAL PRIMARY KEY,
	 	title 			TEXT,
	 	description		TEXT,
	 	priority		NUMERIC,
	 	date_created	TIMESTAMP DEFAULT Now(),
	 	status 			TEXT,
	 	author INTEGER,
	 	FOREIGN KEY(author) REFERENCES users(id)
	 	)`
	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	// populate ticket table
	sql = `INSERT INTO tickets (title, description, priority, status, author)
		VALUES ('Requesting VPN access', 'Working from home, need VPN access please.', 2, 'Not assigned', 1),
			   ('Increase database storage size', 'I am on the dev team and we need the prod database storage to be increased by 10GB to 110GB.', 3, 'Not assigned', 1)
				`
	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	fmt.Print("sql done")
}

/*
	Created using tutorial: https://tutorialedge.net/golang/creating-restful-api-with-golang/
*/
func main() {
	//data base connection
	log.Println("Connecting to database...")
	db.SetDbConn()
	// Get the value of an Environment Variable
	fmt.Println(os.Getenv("HOST"))
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to database!")

	// create tables
	log.Println("Creating tables")
	createTables(db)

	// handle requests
	log.Println("Server started on port :8080...")
	handleRequests()

}
