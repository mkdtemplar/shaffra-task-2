package main

import (
	"database/sql"
	"errors"
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
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/test?sslmode=disable")
	if err != nil {
		log.Fatal(errors.New("error opening database"))
	}

	log.Println("Connected to database:")
	
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/create", createUser)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))

	defer db.Close()
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := db.Query("SELECT name FROM users")
		if err != nil {
            fmt.Fprintf(w, "Failed to retrieve users: %v", err)
            return
        }
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
