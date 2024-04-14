package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type ConnectionPool struct {
	mu          sync.Mutex
	connections chan *sql.DB
	maxConn     int
	minConn     int
}

func NewConnectionPool(maxConn int, minConn int) *ConnectionPool {

	return &ConnectionPool{
		connections: make(chan *sql.DB, maxConn),
		maxConn:     maxConn,
		minConn:     minConn,
	}

}

func (pool *ConnectionPool) newConn() (*sql.DB, error) {
	connectionString := "postgres://new_connection_pool_test:satapril2024@localhost:5432/connection_pool_test?sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error in creating new connection: %v", err)
	}
	return db, nil
}

func (pool *ConnectionPool) initalizeConnections() error {

	for i := 0; i < pool.minConn; i++ {
		conn, err := pool.newConn()

		if err != nil {
			return fmt.Errorf("error in initalizing connection: %v", err)
		}
		pool.connections <- conn
	}
	return nil
}

func (pool *ConnectionPool) Acquire() (*sql.DB, error) {

	if conn, ok := <-pool.connections; ok {
		return conn, nil
	}
	return nil, errors.New("no new connection available")
}

func (pool *ConnectionPool) Release(conn *sql.DB) {

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if err := conn.Ping(); err != nil {
		conn.Close()
		newConn, err := pool.newConn()
		if err != nil {
			log.Printf("Error closing connection %v\n", err)
			return
		}
		conn = newConn
	}
	pool.connections <- conn
}

func SimulateTooManyConnUsingPool() {
	fmt.Println("simulate too many using pool...")
	startTime := time.Now()
	pool := NewConnectionPool(10, 5)
	pool.initalizeConnections()
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			db, err := pool.Acquire()
			if err != nil {
				log.Printf("Error opening connection (%d): %v", i, err)
				return
			}

			_, execErr := db.Exec("SELECT pg_sleep(0.01)")
			if execErr != nil {
				log.Printf("Error executing query (%d): %v", i, execErr)
			}

			pool.Release(db)
		}(i)
	}
	wg.Wait()
	fmt.Printf("Time for connection pool completed %d ms\n", time.Since(startTime).Milliseconds())
}
