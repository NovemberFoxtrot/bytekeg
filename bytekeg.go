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

	for i := 0; i < 100; i++ {
		r.Db("test").Table("Table1").Insert(map[string]interface{}{"id": i, "total": 1, "correct": 1, "incorrect": 15}).Exec(session)
	}

	var response []interface{}

	query := r.Db("test").Table("Table1").Between(1, 10).OrderBy("id")

	rows, err := query.Run(session)

	if err != nil {
		panic(err)
	}

	err = rows.ScanAll(&response)

	log.Println(response)
}
