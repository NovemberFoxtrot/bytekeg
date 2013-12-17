package main

import (
	r "github.com/dancannon/gorethink"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

type TimeConfig struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)

		if err != nil {
			return nil, err
		}

		conn.SetDeadline(time.Now().Add(rwTimeout))

		return conn, nil
	}
}

func FetchURL(theurl string) string {
	var client *http.Client

	if proxy := os.Getenv("http_proxy"); proxy != `` {
		proxyUrl, err := url.Parse(proxy)
		CheckError(err)

		transport := http.Transport{
			Dial:  TimeoutDialer(5*time.Second, 5*time.Second), // connect, read/write
			Proxy: http.ProxyURL(proxyUrl),
		}

		client = &http.Client{Transport: &transport}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest(`GET`, theurl, nil)
	CheckError(err)

	resp, err := client.Do(req)
	CheckError(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)

	return string(body)
}

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

	for i := 3000; i < 3001; i++ {
		www := FetchURL("http://www.cnn.com")
		row := map[string]interface{}{"id": i, "g1": 6771, "www": www}
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
