package main

import (
	r "github.com/dancannon/gorethink"
	"log"
	"math/rand"
	"time"
)

type datum struct {
	count int
	data  []string
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var response []interface{}
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

	objects := []interface{}{}

	for i := 1100; i < 2200; i++ {
		row := map[string]interface{}{"id": i, "g1": 6771}
		objects = append(objects, row)
	}

	log.Println(objects)

	query := r.Db("test").Table("Table1").Insert(objects)

	_, err = query.Run(session)

	if err != nil {
		panic(err)
	}

	query = r.Db("test").Table("Table1").OrderBy("id")

	rows, err := query.Run(session)

	if err != nil {
		panic(err)
	}

	err = rows.ScanAll(&response)

	log.Println(response)
}
