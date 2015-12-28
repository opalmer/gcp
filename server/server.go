package server

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
)

// Listen - Constructs the REST api and starts the http server
func Listen(port int, db bolt.DB) {
	router := NewRouter()

	// Run the server
	log.Print(fmt.Sprintf("Listening at http://127.0.0.1:%d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
