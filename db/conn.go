package db

import "os"

func SetDbConn() {
	// Set Environment Variables
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "5432")
	os.Setenv("USER", "postgres")
	os.Setenv("PASSWORD", "postgres")
	os.Setenv("NAME", "ticketdb")
}
