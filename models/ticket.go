package models

// User schema of the user table
type Ticket struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Priority     int64  `json:"priority"`
	Date_Created string `json:"date_created"`
	Status       string `json:"status"`
	Author       int64  `json:"author"`
}

// https://codesource.io/build-a-crud-application-in-golang-with-postgresql/
