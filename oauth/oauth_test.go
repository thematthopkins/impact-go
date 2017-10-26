package oauth_test

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thematthopkins/impact-go/oauth"
	"github.com/thematthopkins/impact-go/testdb"
)

func addAccessToken(db *sql.DB, userID oauth.UserID, expiration time.Time) (oauth.AccessToken, error) {
	sessionID, err := oauth.AddSession(db, "impact.development", 1234)
	if err != nil {
		return oauth.AccessToken{}, err
	}

	refreshTokenExpiration := time.Now().Add(time.Second * 10)
	accessToken, _, err := oauth.AddSessionToken(db, sessionID, expiration, refreshTokenExpiration)
	if err != nil {
		return oauth.AccessToken{}, err
	}

	return accessToken, nil
}

func TestValid(t *testing.T) {
	db := testdb.Setup()
	accessToken, err := addAccessToken(db, 1234, time.Now().Add(time.Second))
	assert.NoError(t, err)

	var request = httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken.Token))

	userID, err := oauth.Validate(request, db)

	assert.NoError(t, err)
	assert.True(t, userID == 1234, "user id: ", userID)
}

func TestExpired(t *testing.T) {
	db := testdb.Setup()
	accessToken, err := addAccessToken(db, 2222, time.Now().Add(-time.Second))
	assert.NoError(t, err)

	var request = httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken.Token))

	_, err = oauth.Validate(request, db)
	assert.True(t, errors.Cause(err) == oauth.ErrSessionInvalid, "err: ", err)
}

func TestInvalid(t *testing.T) {
	db := testdb.Setup()
	var request = httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Bearer invalidToken")
	_, err := oauth.Validate(request, db)
	assert.True(t, errors.Cause(err) == oauth.ErrSessionInvalid, "err: ", err)
}
