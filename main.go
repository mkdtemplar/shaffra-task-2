package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var wg sync.WaitGroup

func main() {
	db, _ = sql.Open("postgres", "user=postgres dbname=test sslmode=disable")

	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/create", createUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		rows, _ := db.Query("SELECT name FROM users")
		defer rows.Close()

		for rows.Next() {
			var name string
			rows.Scan(&name)
			fmt.Fprintf(w, "User: %s\n", name)
		}
	}()

	wg.Wait()
}

func createUser(w http.ResponseWriter, r *http.Request) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		time.Sleep(5 * time.Second) // Simulate a long database operation

		username := r.URL.Query().Get("name")
		_, err := db.Exec("INSERT INTO users (name) VALUES ('" + username + "')")

		if err != nil {
			fmt.Fprintf(w, "Failed to create user: %v", err)
			return
		}

		fmt.Fprintf(w, "User %s created successfully", username)
	}()

	wg.Wait()
}