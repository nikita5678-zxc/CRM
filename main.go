package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var origin = "phoneCall"

type Ticket struct {
	Id       string `json:"id"`
	Origin   string `json:"origin"`
	ClientId string `json:"clientId"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		panic("DATABASE_URL environment variable not set")
	}

	conn, err1 := pgx.Connect(context.Background(), databaseUrl)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err1)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	r := mux.NewRouter()
	r.HandleFunc("/create/ticket", func(w http.ResponseWriter, r *http.Request) {
		CreateTicketHandler(w, r, conn)
	}).Methods("POST")
	r.HandleFunc("/get/ticket/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetTicketHandler(w, r, conn)
	}).Methods("GET")

	err2 := http.ListenAndServe(":8080", r)
	if err2 != nil {
		panic(err2)
	}

}

func CreateTicketHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
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

	ticketId, err := InsertInto(conn, ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err1 := json.NewEncoder(w).Encode(ticketId)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
	fmt.Printf("ticket created: %s\n", ticketId)
}

func InsertInto(conn *pgx.Conn, ticket Ticket) (string, error) {
	_, err := conn.Exec(
		context.Background(),
		"INSERT INTO ticket (id, origin, clientId) VALUES ($1, $2, $3)",
		ticket.Id, ticket.Origin, ticket.ClientId,
	)
	if err != nil {
		return "", err
	}
	return ticket.Id, nil
}

func GetTicketHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	vars := mux.Vars(r)
	id := vars["id"]
	ticket, err := GetTicket(conn, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

func GetTicket(conn *pgx.Conn, userId string) (Ticket, error) {
	var ticket Ticket
	err := conn.QueryRow(
		context.Background(),
		"SELECT id, origin, clientId FROM ticket WHERE id = $1",
		userId,
	).Scan(&ticket.Id, &ticket.Origin, &ticket.ClientId)

	return ticket, err
}
