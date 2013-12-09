package main

import (
	r "github.com/dancannon/gorethink"
	"log"
	"time"
)

type datum struct {
	count int
	data  []string
}

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

	r.Db("test").TableCreate("Table1").Exec(session)
	r.Db("test").Table("Table1").Insert(map[string]interface{}{"id": 6, "total": 1, "correct": 1, "incorrect": 15}).Exec(session)

	var response interface{}

	query := r.Db("test").Table("Table1").Get(6)

	r, err := query.RunRow(session)

	if err != nil {
		panic(err)
	}

	err = r.Scan(&response)

	log.Println(response)
}
