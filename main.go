package main

import (
	"./db"
	"./server"
	"flag"
)

func main() {
	// Command line arguments
	port := flag.Int(
		"port", 8000,
		"The local port to run the REST API on.")
	databasePath := flag.String(
		"database", ".goback.db",
		"The path to database where goback should store transient information.")
	serverMode := flag.Bool(
		"server", false, "If provided, run goback in server mode.")

	flag.Parse()

	db := db.Open(*databasePath)

	if *serverMode {
		server.Listen(*port, *db)
	}
}
