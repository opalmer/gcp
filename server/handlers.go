package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func putServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverName := vars["name"]
	fmt.Fprintln(w, serverName)
}

func getServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverName := vars["name"]
	fmt.Fprintln(w, serverName)
}
