package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type Feed struct {
	IdItem   int
	IdOffer  int
	Price    float64
	Title    string
	Brand    string
	Category string
	InPromo  bool
}

func main() {
	db, err := sql.Open("sqlite3", "file:./file-generator-db")
	defer db.Close()
	checkErr(err)

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	checkErr(err)
	fmt.Println(version)

	GenerateFeed(db)
}

func GenerateFeed(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM `main`.`feed`")
	checkErr(err)

	var feedCollection []Feed
	for rows.Next() {
		var f Feed
		err = rows.Scan(&f.IdItem, &f.IdOffer, &f.Price, &f.Title, &f.Brand, &f.Category, &f.InPromo)
		checkErr(err)
		feedCollection = append(feedCollection, f)
		time.Sleep(1 * time.Microsecond)
	}

	// TODO generate CSV file using concurrency patterns
	// Ref: https://betterprogramming.pub/file-processing-using-concurrency-with-golang-9e08920fab63
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
