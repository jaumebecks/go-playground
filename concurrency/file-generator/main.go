package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const MinGoroutines = 10

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

// main
// Ref: https://betterprogramming.pub/file-processing-using-concurrency-with-golang-9e08920fab63
func main() {
	db, err := sql.Open("sqlite3", "file:./file-generator-db")
	defer db.Close()
	checkErr(err)

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	checkErr(err)
	log.Printf("SQLite Version -> %s", version)

	GenerateFeedSequentially(db)
	GenerateFeedConcurrently(db)
}

func track(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s, execution time %s\n", name, time.Since(start))
	}
}

func GenerateFeedSequentially(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM `main`.`feed`")
	checkErr(err)

	masterFeed := GenerateMasterFeed(rows, err)
	_ = GenerateSpecificFeed1Sequentially(masterFeed)
}

func GenerateFeedConcurrently(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM `main`.`feed`")
	checkErr(err)

	masterFeed := GenerateMasterFeed(rows, err)
	GenerateSpecificFeed1Concurrently(masterFeed)
}

func GenerateMasterFeed(rows *sql.Rows, err error) MasterFeed {
	masterFeed := MasterFeed{}

	for rows.Next() {
		var r MasterFeedRow
		err = rows.Scan(&r.IdItem, &r.IdOffer, &r.Price, &r.Title, &r.Brand, &r.Category, &r.InPromo)
		checkErr(err)
		masterFeed.AddRow(r)
	}

	return masterFeed
}

func GenerateSpecificFeed1Sequentially(masterFeed MasterFeed) SpecificFeed1 {
	defer track("GenerateSpecificFeed1Sequentially")()
	feed := SpecificFeed1{}

	file, err := os.Create("sequential1.csv")
	defer file.Close()
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()

	for _, row := range masterFeed.Rows {
		r := createSpecificRow1(row)
		feed.AddRow(r)

		writeSpecific1Row(w, r)
	}

	return feed
}

func writeSpecific1Row(w *csv.Writer, r SpecificFeed1Row) {
	_ = w.Write([]string{
		r.Id,
		fmt.Sprintf("%f", r.Price),
		r.Title,
		r.Category,
		fmt.Sprintf("%v", r.InPromo),
	})
	time.Sleep(10 * time.Millisecond)
}

func GenerateSpecificFeed1Concurrently(masterFeed MasterFeed) {
	defer track("GenerateSpecificFeed1Concurrently")()

	var wg sync.WaitGroup

	size := len(masterFeed.Rows)
	chunk := size / MinGoroutines
	iterations := MinGoroutines
	if size%chunk != 0 {
		iterations++
	}

	csvFileFormat := "concurrent.part.%d.csv"
	for i := 0; i < iterations; i++ {
		startOffset := i * chunk
		endOffset := chunk * (i + 1)
		if chunk*(i+1) > size {
			endOffset = size
		}

		wg.Add(1)
		go func(feedPart MasterFeed, iteration int) {
			defer wg.Done()
			file, err := os.Create(fmt.Sprintf(csvFileFormat, iteration))
			defer file.Close()
			checkErr(err)
			w := csv.NewWriter(file)
			defer w.Flush()

			for _, row := range feedPart.Rows {
				r := createSpecificRow1(row)
				writeSpecific1Row(w, r)
			}
		}(MasterFeed{Rows: masterFeed.Rows[startOffset:endOffset]}, i)
	}

	wg.Wait()

	finalFeedFile, err := os.Create("concurrent1.csv")
	checkErr(err)
	for i := 0; i < iterations; i++ {
		partFileName := fmt.Sprintf(csvFileFormat, i)
		part, err := os.Open(partFileName)
		checkErr(err)
		_, err = io.Copy(finalFeedFile, part)
		part.Close()
		checkErr(err)
		err = os.Remove(partFileName)
		checkErr(err)
	}
	finalFeedFile.Close()
}

func createSpecificRow1(row MasterFeedRow) SpecificFeed1Row {
	r := SpecificFeed1Row{
		Id:       fmt.Sprintf("online:es:ES:%d:%d", row.IdItem, row.IdOffer),
		Price:    row.Price,
		Title:    fmt.Sprintf("%s - %s", row.Title, row.Brand),
		Category: row.Category,
		InPromo:  row.InPromo,
	}
	return r
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
