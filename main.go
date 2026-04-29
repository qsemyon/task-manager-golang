package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	initDatabase(connStr)
	defer db.Close()

	http.HandleFunc("/tasks", tasksHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	fmt.Printf("Сервер запущен на http://localhost:%s\n", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}

func initDatabase(connStr string) {
	var err error

	for i := 1; i <= 15; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			fmt.Println("Успешное подключение к базе данных!")
			break
		}

		fmt.Printf("[Попытка %d] База еще не готова, ждем...\n", i)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Не удалось подключиться к БД после 15 попыток:", err)
	}

	query := `CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY, 
        title TEXT NOT NULL, 
        done BOOLEAN DEFAULT FALSE
    )`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		getTasks(w)
	case http.MethodPost:
		createTask(w, r)
	case http.MethodPut:
		updateTask(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getTasks(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, title, done FROM tasks ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	err := db.QueryRow("INSERT INTO tasks(title) VALUES($1) RETURNING id", t.Title).Scan(&t.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return
	}
	_, err := db.Exec("UPDATE tasks SET done = $1 WHERE id = $2", t.Done, t.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
