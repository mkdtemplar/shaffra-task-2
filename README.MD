In main.go database needs to close.

defer db.Close
connection string should be : db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/test?sslmode=disable")

handle error:
if err != nil {
		log.Fatal(errors.New("error opening database"))
	}

geUser function rows error should be handled:
if err != nil {
            fmt.Fprintf(w, "Failed to retrieve users: %v", err)
            return
        }