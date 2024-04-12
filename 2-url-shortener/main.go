package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type AddRedirectParam struct {
	ShortURL  string `json:"short_url"`
	TargetURL string `json:"target_url"`
}

func handleAllQueries(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		addNewURL(w, r)
		return
	}
	redirectURL(w, r)
}

func addNewURL(w http.ResponseWriter, r *http.Request) {
	var redirect AddRedirectParam
	if err := json.NewDecoder(r.Body).Decode(&redirect); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tx == nil {
		http.Error(w, "transaction is nil", http.StatusInternalServerError)
		return
	}

	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Printf("Rollback transaction: %s\n", err)
		}
	}(tx)

	if _, err := db.Exec("INSERT INTO urls (short_url, target_url) VALUES ($1, $2)", redirect.ShortURL, redirect.TargetURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write([]byte(fmt.Sprintf("Added redirect from %s to %s", redirect.ShortURL, redirect.TargetURL))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	var targetURL string
	if err := db.QueryRow("SELECT target_url FROM urls WHERE short_url = $1", shortURL).Scan(&targetURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, targetURL, http.StatusFound)
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./urls.db")
	if err != nil {
		fmt.Printf("start database: %s\n", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("close database: %s\n", err)
		}
	}(db)

	http.HandleFunc("/", handleAllQueries)
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("start server: %s\n", err)
	}
}
