package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
)

var users map[int]string
var origin = "phoneCall"

type Ticket struct {
	Id       string `json:"id"`
	Origin   string `json:"origin"`
	ClientId string `json:"clientId"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/product", ProductHandler).Methods("POST")
	users = make(map[int]string)
	err := http.ListenAndServe(":8123", r)
	if err != nil {
		panic(err)
	}

}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("number")
	CreateTicket(origin, number)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func CreateTicket(origin string, client_id string) {
	id := rand.Intn(9000000)
	users[id] = client_id
	GetTicket(id, client_id)
}

func GetTicket(id int, client_id string) {
	ticket := Ticket{
		Id:       fmt.Sprintf("%d", id),
		Origin:   origin,
		ClientId: client_id,
	}
	ticketJson, err := json.MarshalIndent(ticket, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ticketJson))
}
