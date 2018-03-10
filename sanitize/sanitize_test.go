package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thematthopkins/impact-go/auth"
	"github.com/thematthopkins/impact-go/testdb"
)

func TestInvalidClientName(t *testing.T) {
	db := testdb.Setup()
	_, err := auth.AddSession(db, "invalidClientID", 1234)
	assert.Error(t, err, "failed to find oauth client: invalidClientID")
}
