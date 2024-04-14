package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

func NewConn() (*sql.DB, error) {
	connectionString := "postgres://new_connection_pool_test:satapril2024@localhost:5432/connection_pool_test?sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error in creating new connection: %v", err)
	}
	// setting max connection as 100
	db.SetMaxOpenConns(100)

	return db, nil
}

func SimulateTooManyConn() {
	fmt.Println("simulate too many...")
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 190; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			db, err := NewConn()
			if err != nil {
				log.Printf("Error opening connection (%d): %v", i, err)
				return
			}
			defer db.Close()

			_, execErr := db.Exec("SELECT pg_sleep(0.01)")
			if execErr != nil {
				log.Printf("Error executing query (%d): %v", i, execErr)
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("Time for connection pool completed %d ms\n", time.Since(startTime).Milliseconds())
}
