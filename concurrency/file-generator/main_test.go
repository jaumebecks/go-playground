package main

import (
	"database/sql"
	"testing"
)

func BenchmarkGenerateFeedSequentially(b *testing.B) {
	b.ResetTimer()
	db, err := sql.Open("sqlite3", "file:./file-generator-db")
	defer db.Close()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		GenerateFeedSequentially(db)
	}
}

func BenchmarkGenerateConcurrently(b *testing.B) {
	b.ResetTimer()
	db, err := sql.Open("sqlite3", "file:./file-generator-db")
	defer db.Close()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		GenerateFeedConcurrently(db)
	}
}
