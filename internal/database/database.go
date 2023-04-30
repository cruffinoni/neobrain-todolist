package database

import (
	"fmt"
	"log"
	"time"

	"github.com/cruffinoni/neobrain-todolist/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	instance *sqlx.DB
}

const (
	maxRetries = 5
	retryDelay = 5 * time.Second
)

func connectToDatabase(config *config.Database) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempt to connect to the database (try %d/%d)", i+1, maxRetries)
		db, err = sqlx.Open(
			"mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.Username, config.Password, config.Host, config.Port),
		)
		if err == nil {
			if err = db.Ping(); err == nil {
				log.Printf("Succesfully connected to the database")
				return db, nil
			}
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to the database after %d retries: %w", maxRetries, err)
}

func NewDB(config *config.Database) (*DB, error) {
	dbConnection, err := connectToDatabase(config)
	if err != nil {
		return nil, err
	}
	return &DB{instance: dbConnection}, nil
}

func (db *DB) AddTask(task string) (int64, error) {
	result, err := db.instance.Exec(`INSERT INTO todolist.tasks (task) VALUES (?)`, task)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) DeleteTask(taskID int64) error {
	_, err := db.instance.Exec("DELETE FROM todolist.tasks WHERE id = ?", taskID)
	return err
}

func (db *DB) MarkTaskAsDone(taskID int64) error {
	_, err := db.instance.Exec("UPDATE todolist.tasks SET done = 1 WHERE id = ?", taskID)
	return err
}

type TaskFilter string

const (
	TaskFilterDone    TaskFilter = "done"
	TaskFilterNotDone TaskFilter = "notdone"
	TaskFilterAll     TaskFilter = "all"
)

func (db *DB) GetTasks(filter TaskFilter) ([]*Task, error) {
	var tasks []*Task

	query := "SELECT id, task, done FROM todolist.tasks"
	if filter == TaskFilterDone {
		query += " WHERE done = 1"
	} else if filter == TaskFilterNotDone {
		query += " WHERE done = 0"
	}

	err := db.instance.Select(&tasks, query)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (db *DB) ImportTasks(tasks []*Task) error {
	tx, err := db.instance.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO todolist.tasks (task, done) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, task := range tasks {
		_, err = stmt.Exec(task.Task, task.Done)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
