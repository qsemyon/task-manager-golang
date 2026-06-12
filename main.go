package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type App struct {
	db *sql.DB
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	dbConn := initDatabase(connStr)
	defer dbConn.Close()

	app := &App{db: dbConn}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", app.getTasks)
	mux.HandleFunc("POST /tasks", app.createTask)
	mux.HandleFunc("PUT /tasks", app.updateTask)
	mux.HandleFunc("DELETE /tasks", app.deleteTask)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	fmt.Printf("Сервер запущен на http://localhost:%s\n", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, mux))
}
