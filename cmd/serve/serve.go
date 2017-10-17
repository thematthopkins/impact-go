package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {

	phpURLStr := "https://api.impact.dev"
	phpURL, err := url.Parse(phpURLStr)
	if err != nil {
		panic(err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		panic(err)
	}

	http.Handle("/", httputil.NewSingleHostReverseProxy(phpURL))

	handleFuncAuthenticated("/report", generateReport, db)

	https := os.Getenv("ENABLE_HTTPS") != ""

	if https {
		fmt.Println("serving https")
		port := ":4443"
		err = http.ListenAndServeTLS(port, "ssl/cert.pem", "ssl/cert.pem", nil)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to bind to port %v:  %v", port, err))
		}
	} else {
		fmt.Println("serving http")
		port := ":8080"
		err = http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to bind to port %v:  %v", port, err))
		}
	}
}

func handleFuncWithPanicRecovery(path string, fn func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Panicked: ", r)
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	})
}

func handleFuncWithDb(path string, fn func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) {
	handleFuncWithPanicRecovery(path, func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	})
}

func handleFuncAuthenticated(path string, fn func(http.ResponseWriter, *http.Request, *sql.DB, ImpactUserID), db *sql.DB) {
	handleFuncWithPanicRecovery(path, func(w http.ResponseWriter, r *http.Request) {
		userID, err := authenticate(r, db)
		if err == ErrSessionInvalid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else if err != nil {
			fmt.Println("error authenticating: ", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
		} else {
			generateReport(w, r, db, userID)
		}
	})
}

type ImpactUserID int

var ErrSessionInvalid = errors.New("session invalid")

func authenticate(r *http.Request, db *sql.DB) (ImpactUserID, error) {
	tokens, ok := r.Header["Authorization"]
	if !ok || len(tokens) != 1 {
		return 0, errors.New("no authorization header")
	}
	token := strings.TrimPrefix(tokens[0], "Bearer ")
	var userID ImpactUserID
	err := db.QueryRow(`
		select
			oauth_sessions.owner_id
		from 
			oauth_access_tokens
			join oauth_sessions on oauth_access_tokens.session_id = oauth_sessions.id
		where
			oauth_access_tokens.id = $1
			and oauth_access_tokens.expire_time > EXTRACT(EPOCH FROM now())
		`, token).Scan(&userID)

	if err == sql.ErrNoRows {
		return 0, ErrSessionInvalid
	} else if err != nil {
		return 0, err
	} else {
		return userID, nil
	}
}

func generateReport(w http.ResponseWriter, r *http.Request, db *sql.DB, id ImpactUserID) {
	fmt.Fprintf(w, "Hello %v", id)
}
