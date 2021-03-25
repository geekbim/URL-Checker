package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UrlNotFound struct {
	Olxid string
	Url   string
	Star  string
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

	olxDB, err := db.Query("select olxid, url from olx where created_at >= CURDATE() - INTERVAL 1 DAY and created_at < CURDATE() + INTERVAL 1 DAY  order by created_at desc")
	if err != nil {
		panic(err.Error())
	}
	urlnotfound := UrlNotFound{}
	urlnotfounds := []UrlNotFound{}
	for olxDB.Next() {
		var olxid, url string

		err = olxDB.Scan(&olxid, &url)
		if err != nil {
			panic(err.Error())
		}
		urlnotfound.Olxid = olxid
		urlnotfound.Url = url
		urlnotfound.Star = "false"
		urlnotfounds = append(urlnotfounds, urlnotfound)
	}
	defer db.Close()

	for _, urlnotfound := range urlnotfounds {
		resp, err := http.Get(urlnotfound.Url)
		if err != nil {
			log.Fatal(err)
		}

		// Print the HTTP Status Code and Status Name
		fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Println("HTTP Status is in the 2xx range")
		} else {
			fmt.Println("Argh! Broken")
			db.Exec("update olx set archived = 1 where olxid = ?",
				urlnotfound.Olxid)
			db.Exec("INSERT INTO olx_archive (deleted_at, created_at, olxid, star) "+
				"VALUES( null, CURRENT_TIMESTAMP, ?, ? ) "+
				"ON DUPLICATE KEY "+
				"UPDATE deleted_at = null, updated_at = CURRENT_TIMESTAMP, star = ?",
				urlnotfound.Olxid,
				urlnotfound.Star,
				urlnotfound.Star)

			fmt.Println(urlnotfound)
		}
	}

}
