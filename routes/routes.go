package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/thematthopkins/impact-go/auth"
	"github.com/thematthopkins/impact-go/verificationreport"
)

// Handle adds golang handled routes to the http multiplexer
func Handle(mux *http.ServeMux, db *sql.DB) {
	handleFuncAuthenticated("/report", verificationreport.Export, db)
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

func handleFuncAuthenticated(path string, fn func(http.ResponseWriter, *http.Request, *sql.DB, auth.UserID), db *sql.DB) {
	handleFuncWithPanicRecovery(path, func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.Validate(r, db)
		if err == auth.ErrSessionInvalid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else if err != nil {
			fmt.Println("error authenticating: ", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
		} else {
			fn(w, r, db, userID)
		}
	})
}
