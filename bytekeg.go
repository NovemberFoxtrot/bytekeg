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

	objects := make([]interface{}, 0)

	for i := 1; i < 10; i++ {
		total := rand.Intn(1000)

		object := map[string]interface{}{
			"id":        i,
			"total":     total,
			"correct":   total,
			"incorrect": total,
		}

		objects = append(objects, object)
	}

	// table.Insert([]interface{}{map[string]interface{}{"name": "Joe"}, map[string]interface{}{"name": "Paul"}}).RunWrite(sess)

	log.Println(objects)

	query := r.Db("test").Table("Table1").Replace(objects)

	rowz, err := query.Run(session)

	log.Println(rowz)

	if err != nil {
		panic(err)
	}

	query = r.Db("test").Table("Table1").Between(1, 10).OrderBy("id")

	rows, err := query.Run(session)

	if err != nil {
		panic(err)
	}

	err = rows.ScanAll(&response)

	log.Println(response)
}
