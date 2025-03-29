package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var origin = "phoneCall"

var (
	ticketBD = make(map[string]Ticket)
	BdMutex  sync.Mutex
)

type Ticket struct {
	Id       string `json:"id"`
	Origin   string `json:"origin"`
	ClientId string `json:"clientId"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create/ticket", CreateTicketHandler).Methods("POST")
	r.HandleFunc("/get/ticket/{id}", GetTicketHandler).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

}

func CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("clientId")
	if clientId == "" {
		http.Error(w, "ClientId is empty", http.StatusBadRequest)
		return
	}
	rand.Seed(time.Now().UnixNano())
	id := fmt.Sprintf("%07d", rand.Intn(10000000))
	ticket := Ticket{
		Id:       id,
		Origin:   origin,
		ClientId: clientId,
	}

	BdMutex.Lock()
	ticketBD[id] = ticket
	BdMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
	fmt.Printf("ticket created: %+v\n", ticket)
}

func GetTicketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	BdMutex.Lock()
	ticket, errbool := ticketBD[id]
	BdMutex.Unlock()

	if !errbool {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}
