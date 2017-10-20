package auth

import(
	"database/sql"
	"errors"
	"strings"
	"net/http"
)

// ImpactUserID validated identifier of the impact_user
type ImpactUserID int

// ErrSessionInvalid failure to authenticate user
var ErrSessionInvalid = errors.New("session invalid")

// Validate user based on http Authorization: Bearer *** token
func Validate(r *http.Request, db *sql.DB) (ImpactUserID, error) {
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
