package db

import (
	"database/sql"
	"fmt"
	"log"
)

// capital letter at start means public visibility
func SetDbConn() *sql.DB {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postGres!"
		dbname   = "ticketdb"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to database!")
	return db

}

func CreateTables(db *sql.DB) {
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
}
