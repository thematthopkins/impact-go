package verificationreport

import (
	"github.com/thematthopkins/impact-go/auth"
	"net/http"
	"database/sql"
	"fmt"
)

func Export(w http.ResponseWriter, r *http.Request, db *sql.DB, id auth.ImpactUserID) {
	fmt.Fprintf(w, "Verification Report %v", id)
}
