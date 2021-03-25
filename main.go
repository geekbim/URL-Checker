package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Url struct {
	Url string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "hwhwhwlol"
	dbName := "odb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	ErrorCheck(err)

	PingDB(db)

	return db
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func PingDB(db *sql.DB) {
	err := db.Ping()
	ErrorCheck(err)
}

func main() {

	db := dbConn()

	olxDB, err := db.Query("select url from olx where created_at >= NOW() - INTERVAL 1 DAY order by created_at desc")
	if err != nil {
		panic(err.Error())
	}
	link := Url{}
	links := []Url{}
	for olxDB.Next() {
		var url string
		err = olxDB.Scan(&url)
		if err != nil {
			panic(err.Error())
		}
		link.Url = url
		links = append(links, link)
	}
	defer db.Close()

	for _, link := range links {
		resp, err := http.Get(link.Url)
		if err != nil {
			log.Fatal(err)
		}

		// Print the HTTP Status Code and Status Name
		fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Println("HTTP Status is in the 2xx range")
		} else {
			fmt.Println("Argh! Broken")

		}
	}

}
