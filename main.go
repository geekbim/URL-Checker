package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UrlNotFound struct {
	CreatedAt string
	Olxid     string
	Url       string
	Star      string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "8ismillah"
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

	olxDB, err := db.Query("select created_at, olxid, url from olx where created_at >= CURDATE() - INTERVAL 30 DAY and created_at < CURDATE() + INTERVAL 1 DAY  order by created_at desc")
	if err != nil {
		panic(err.Error())
	}
	urlnotfound := UrlNotFound{}
	urlnotfounds := []UrlNotFound{}
	for olxDB.Next() {
		var created_at, olxid, url string

		err = olxDB.Scan(&created_at, &olxid, &url)
		if err != nil {
			panic(err.Error())
		}
		urlnotfound.CreatedAt = created_at
		urlnotfound.Olxid = olxid
		urlnotfound.Url = url
		urlnotfound.Star = "false"
		urlnotfounds = append(urlnotfounds, urlnotfound)
	}
	defer db.Close()

	for _, urlnotfound := range urlnotfounds {
		resp, err := http.Get("https://www.olx.co.id/item/" + urlnotfound.Olxid)
		if err != nil {
			log.Fatal(err)
		}
		// Resource leak if response body isn't closed
		defer resp.Body.Close()

		// Print the HTTP Status Code and Status Name
		fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
		fmt.Println("URL:", urlnotfound)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Println("HTTP Status is in the 2xx range")
		} else {
			fmt.Println("Argh! Broken")
			// db.Exec("update olx set archived = 1 where olxid = ?",
			// 	urlnotfound.Olxid)
			// db.Exec("INSERT INTO olx_archive (deleted_at, created_at, olxid, star) "+
			// 	"VALUES( null, CURRENT_TIMESTAMP, ?, ? ) "+
			// 	"ON DUPLICATE KEY "+
			// 	"UPDATE deleted_at = null, updated_at = CURRENT_TIMESTAMP, star = ?",
			// 	urlnotfound.Olxid,
			// 	urlnotfound.Star,
			// 	urlnotfound.Star)

			db.Exec("update olx set deleted_at = CURRENT_TIMESTAMP where olxid = ?",
				urlnotfound.Olxid)

			fmt.Println(urlnotfound)
		}
	}

}
