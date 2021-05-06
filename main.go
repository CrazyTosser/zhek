package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	"log"
	"os"
)

var cfg Config

func main() {
	ctx := context.Background()
	cont, cancel := context.WithCancel(ctx)
	dbpool, err := pgxpool.Connect(cont, "postgresql://postgres:secret@localhost:5432/mir")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		println("Shutting down")
		dbpool.Close()
		cancel()
	}()

	sport := os.Getenv("SPORT")
	if sport == "" {
		sport = "2156"
		log.Printf("Defaulting to sport %s", sport)
	}

	log.Printf("Listening on sport %s", sport)
	go initScanner(dbpool, cont, &cfg)
	go ServiceUdp(cont, "0.0.0.0", sport, dbpool, &cfg)
	StartServer(dbpool, ctx)
}
