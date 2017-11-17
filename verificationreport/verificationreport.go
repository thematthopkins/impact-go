package verificationreport

import (
	"database/sql"
	"fmt"
	"github.com/thematthopkins/impact-go/oauth"
	"net/http"
)

func Export(w http.ResponseWriter, r *http.Request, db *sql.DB, id oauth.UserID) {
	fmt.Fprintf(w, "Verification Report %v", id)
}
