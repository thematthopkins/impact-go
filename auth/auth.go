package auth

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

// UserID is the validated identifier of the impact_user
type UserID int

// ErrSessionInvalid failure to authenticate user
var ErrSessionInvalid = errors.New("session invalid")

var ErrNoAuthHeader = errors.New("no authorization header")

// Validate user based on http Authorization: Bearer *** token
func Validate(r *http.Request, db *sql.DB) (UserID, error) {
	tokens, ok := r.Header["Authorization"]
	if !ok || len(tokens) != 1 {
		return 0, ErrNoAuthHeader
	}
	token := strings.TrimPrefix(tokens[0], "Bearer ")
	now := time.Now().Unix()
	var userID UserID
	err := db.QueryRow(`
		select
			oauth_sessions.owner_id
		from 
			oauth_access_tokens
			join oauth_sessions on oauth_access_tokens.session_id = oauth_sessions.id
		where
			oauth_access_tokens.id = $1
			and oauth_access_tokens.expire_time > $2
		`, token, now).Scan(&userID)

	if err == sql.ErrNoRows {
		return 0, errors.Wrapf(ErrSessionInvalid, token)
	} else if err != nil {
		panic(err)
	}

	return userID, nil
}

type ClientName string
type SessionID int64

func AddSession(db *sql.DB, oauthClientName ClientName, userID UserID) (SessionID, error) {
	var clientID string
	err := db.QueryRow(`
		select id from oauth_clients where name = $1
	`, oauthClientName).Scan(&clientID)
	if err == sql.ErrNoRows {
		return 0, errors.Wrapf(err, "failed to find oauth client", oauthClientName)
	} else if err != nil {
		panic(err)
	}

	var sessionID SessionID
	
	err = db.QueryRow(`
		insert into oauth_sessions(client_id, owner_type, owner_id, created_at, updated_at) values ($1, 'user', $2, now(), now()) returning id
	`, clientID, userID).Scan(&sessionID)

	if err != nil {
		panic(err)
	}
	return sessionID, nil
}

type AccessToken struct {
	Token      string
	Expiration time.Time
}

type RefreshToken struct {
	Token      string
	Expiration time.Time
}

func AddSessionToken(db *sql.DB, sessionID SessionID, accessTokenExpiration time.Time, refreshTokenExpiration time.Time) (AccessToken, RefreshToken, error) {
	accessToken := AccessToken{
		Token:      uuid.New().String(),
		Expiration: accessTokenExpiration,
	}
	_, err := db.Exec(`
		insert into oauth_access_tokens(id, session_id, expire_time, created_at, updated_at) values($1, $2, $3, now(), now())
	`, accessToken.Token, sessionID, accessTokenExpiration.Unix())
	if err != nil {
		panic(err)
	}

	refreshToken := RefreshToken{
		Token:      uuid.New().String(),
		Expiration: refreshTokenExpiration,
	}
	_, err = db.Exec(`
		insert into oauth_refresh_tokens(id, access_token_id, expire_time, created_at, updated_at) values($1, $2, $3, now(), now())
	`, refreshToken.Token, accessToken.Token, refreshTokenExpiration.Unix())
	if err != nil {
		panic(err)
	}

	return accessToken, refreshToken, nil
}
