package main

import (
	"auth/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	connection := connectDB()
	if connection == nil {
		log.Panic("Can't connect to database")
	}

	app := Config{
		DB:     connection,
		Models: data.New(connection),
	}

	log.Print("Starting authentication service...")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.routes(),
	}

	// start http server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func getDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")
	counts := 0
	for {
		connection, err := getDB(dsn)
		if err != nil {
			log.Println("DB not ready")
			counts++
		} else {
			log.Println("Connected to DB")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Wait for seconds to connect again")
		time.Sleep(2 * time.Second)
		continue
	}
}
