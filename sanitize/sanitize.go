package sanitize

import "database/sql"

//AssessmentID of a question in the db
type AssessmentID int

type responseID int

func sanitize(assessmentID AssessmentID, db *sql.DB) {
	// rows, err := db.Query(
	// 	// Ensure the name and age parameters only match on placeholder name, not position.
	// 	"SELECT|people|age,name|name=?name,age=?age",
	// 	Named("age", 2),
	// 	Named("name", "Bob"),
	// )
	// if err != nil {
	// 	t.Fatalf("Query: %v", err)
	// }
	// type row struct {
	// 	age  int
	// 	name string
	// }
	// _, err := db.Exec(`

	// 	insert into oauth_access_tokens(id, session_id, expire_time, created_at, updated_at) values($1, $2, $3, now(), now())
	// `, accessToken.Token, sessionID, accessTokenExpiration.Unix())
	// if err != nil {
	// 	panic(err)
	// }

}
