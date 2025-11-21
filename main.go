package main

import (
	"log"
	"os"

	"tracker/pkg/db"
	"tracker/pkg/server"
)

const defaultServerPort = "7540"
const defaultDBFile = "scheduler.db"

func main() {

	// Logger creration
	logger := log.New(os.Stdout, "Tasks tracker web server: ", log.Ldate|log.Ltime)

	// Initialization of database
	err := db.Init(getDBpath())

	// Database closure setting
	defer db.Close()

	// Errors handling
	if err != nil {
		logger.Fatalf("\nDatabase initialization fatal error: %v", err)
	}

	// Start of the tracker web server
	srv := server.NewServer(logger, getServerPort())

	logger.Printf("\nServer started at %s ", srv.HttpServer.Addr)

	err = srv.HttpServer.ListenAndServe()

	// Errors handling
	if err != nil {
		logger.Fatalf("\nServer fatal error: %v", err)
	}

}

// Getting port value from TODO_PORT env variable
// If TODO_PORT is not set, returning a default port,
// set in defaultServerPort constant
func getServerPort() string {

	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = defaultServerPort
	}

	return ":" + port
}

// Getting DB path from TODO_DBFILE env variable
// If TODO_DBFILE is not set, returning a default DB path,
// set in defaultDBFile constant
func getDBpath() string {

	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = defaultDBFile
	}

	return dbFile
}
