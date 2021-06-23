package db

import (
	"database/sql"
	"fmt"
	"github.com/AidanStrong/TicketSystem/models"
)

func GetAllTickets(db *sql.DB) []models.Ticket {
	var tickets []models.Ticket
	sql := `SELECT * FROM tickets`

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var record models.Ticket
		// copy values of record into variable
		err = rows.Scan(&record.Id, &record.Title, &record.Description, 
			&record.Priority, &record.Date_Created, &record.Status, 
			&record.Author)
		if err != nil {
				fmt.Printf("Scan: %v", err)
		}
			
		tickets = append(tickets, record)
	}

	return tickets
}