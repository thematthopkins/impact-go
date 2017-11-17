package auth_test

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thematthopkins/impact-go/auth"
	"github.com/thematthopkins/impact-go/testdb"
)

func addAccessToken(db *sql.DB, userID auth.UserID, expiration time.Time) (auth.AccessToken, error) {
	sessionID, err := auth.AddSession(db, "impact.development", 1234)
	if err != nil {
		return auth.AccessToken{}, err
	}

	refreshTokenExpiration := time.Now().Add(time.Second * 10)
	accessToken, _, err := auth.AddSessionToken(db, sessionID, expiration, refreshTokenExpiration)
	if err != nil {
		return auth.AccessToken{}, err
	}

	return accessToken, nil
}

func TestInvalidClientName(t *testing.T) {
	db := testdb.Setup()
	_, err := auth.AddSession(db, "invalidClientID", 1234)
	assert.Error(t, err, "failed to find oauth client: invalidClientID");	
}

func TestValid(t *testing.T) {
	db := testdb.Setup()
	accessToken, err := addAccessToken(db, 1234, time.Now().Add(time.Second))
	assert.NoError(t, err)

	var request = httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken.Token))

	userID, err := auth.Validate(request, db)

	assert.NoError(t, err)
	assert.True(t, userID == 1234, "user id: ", userID)
}

func TestExpired(t *testing.T) {
	db := testdb.Setup()
	accessToken, err := addAccessToken(db, 2222, time.Now().Add(-time.Second))
	assert.NoError(t, err)

	var request = httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken.Token))

	_, err = auth.Validate(request, db)
	assert.True(t, errors.Cause(err) == auth.ErrSessionInvalid, "err: ", err)
}

func TestInvalid(t *testing.T) {
	db := testdb.Setup()
	var request = httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Bearer invalidToken")
	_, err := auth.Validate(request, db)
	assert.EqualError(t, errors.Cause(err), auth.ErrSessionInvalid.Error())
}

func TestMissingAuthHeader(t *testing.T) {
	db := testdb.Setup()
	var request = httptest.NewRequest("GET", "/", nil)
	_, err := auth.Validate(request, db)
	assert.EqualError(t, errors.Cause(err), auth.ErrNoAuthHeader.Error())
}