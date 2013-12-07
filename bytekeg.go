package main

import (
	r "github.com/dancannon/gorethink"
	"log"
	"time"
)

func main() {

	var session *r.Session

	session, err := r.Connect(map[string]interface{}{
		"address":     "localhost:28015",
		"database":    "test",
		"maxIdle":     10,
		"idleTimeout": time.Second * 10,
	})

	if err != nil {
		panic(err)
	}

	// Insert rows
	r.Db("test").Table("Table1").Insert(map[string]interface{}{"id": 6, "g1": 1, "g2": 1, "num": 15}).Exec(session)

	// Test query
	var response interface{}
	query := r.Db("test").Table("Table1").Get(6)
	r, err := query.RunRow(session)

	if err != nil {
		panic(err)
	}

	err = r.Scan(&response)

	log.Println(response)
}
