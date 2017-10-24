package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/thematthopkins/impact-go/routes"

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
	routes.Handle(http.DefaultServeMux, db)

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
