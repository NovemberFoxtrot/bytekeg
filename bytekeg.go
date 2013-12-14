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

	// objects := make([]interface{}, 0)

	objects := []interface{}{
		map[string]interface{}{"num": 0, "id": 10, "g2": 1, "g1": 1},
		map[string]interface{}{"num": 5, "id": 20, "g2": 2, "g1": 2},
		map[string]interface{}{"num": 10, "id": 30, "g2": 2, "g1": 3},
	}

	row :=	map[string]interface{}{"num": 0, "id": 100, "g2": 1, "g1": 6771}

	objects = append(objects, row)

	log.Println(objects)

	query := r.Db("test").Table("Table1").Insert(objects)

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
