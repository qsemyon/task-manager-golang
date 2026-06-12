package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func initDatabase(connStr string) *sql.DB {
	var err error
	var localDB *sql.DB

	for i := 1; i <= 15; i++ {
		localDB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = localDB.Ping()
		}

		if err == nil {
			fmt.Println("База данных подключена")
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
	_, err = localDB.Exec(query)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}

	return localDB
}

func (a *App) GetAllTasks() ([]Task, error) {
	rows, err := a.db.Query("SELECT id, title, done FROM tasks ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (a *App) AddTask(title string) (Task, error) {
	var t Task
	t.Title = title
	err := a.db.QueryRow("INSERT INTO tasks(title) VALUES($1) RETURNING id", title).Scan(&t.ID)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (a *App) UpdateTask(id int, done bool) error {
	_, err := a.db.Exec("UPDATE tasks SET done = $1 WHERE id = $2", done, id)
	return err
}

func (a *App) DeleteTask(id string) error {
	_, err := a.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
