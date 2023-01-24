package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type MasterFeed struct {
	Rows []MasterFeedRow
}

func (m *MasterFeed) AddRow(r MasterFeedRow) {
	m.Rows = append(m.Rows, r)
}

type MasterFeedRow struct {
	IdItem   int
	IdOffer  int
	Price    float64
	Title    string
	Brand    string
	Category string
	InPromo  bool
}

type SpecificFeed1 struct {
	Rows []SpecificFeed1Row
}

func (f *SpecificFeed1) AddRow(r SpecificFeed1Row) {
	f.Rows = append(f.Rows, r)
}

type SpecificFeed1Row struct {
	Id       string
	Price    float64
	Title    string
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

	GenerateFeedSequentially(db)
}

// GenerateFeedSequentially TODO generate CSV file using concurrency patterns
// Ref: https://betterprogramming.pub/file-processing-using-concurrency-with-golang-9e08920fab63
func GenerateFeedSequentially(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM `main`.`feed`")
	checkErr(err)

	masterFeed := MasterFeed{}
	for rows.Next() {
		var r MasterFeedRow
		err = rows.Scan(&r.IdItem, &r.IdOffer, &r.Price, &r.Title, &r.Brand, &r.Category, &r.InPromo)
		checkErr(err)
		masterFeed.AddRow(r)
		// time.Sleep(1 * time.Microsecond)
	}

	_ = GenerateSpecificFeed1Sequentially(masterFeed)
}

func GenerateSpecificFeed1Sequentially(masterFeed MasterFeed) SpecificFeed1 {
	feed := SpecificFeed1{}
	for _, row := range masterFeed.Rows {
		r := SpecificFeed1Row{
			Id:       fmt.Sprintf("online:es:ES:%d:%d", row.IdItem, row.IdOffer),
			Price:    row.Price,
			Title:    fmt.Sprintf("%s - %s", row.Title, row.Brand),
			Category: row.Category,
			InPromo:  row.InPromo,
		}
		feed.AddRow(r)
		time.Sleep(1 * time.Microsecond)
	}

	return feed
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
