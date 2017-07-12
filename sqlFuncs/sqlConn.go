package sqlFuncs

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
var err error

func sqlConn() {
	db, err = sql.Open("mysql", "root:carnyx007@/")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("\nConnected to the Database.")
}
