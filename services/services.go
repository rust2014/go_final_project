package services

import (
	"database/sql"

	"github.com/rust2014/go_final_project/models"
)

func GetTasks(db *sql.DB) ([]models.Task, error) {
	tasks := []models.Task{}
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTask(db *sql.DB, id int) (*models.Task, error) {
	var task models.Task
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(db *sql.DB, task models.Task) (int64, error) {
	result, err := db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func DoneTask(db *sql.DB, id int, nextDate string) error {
	if nextDate == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		return err
	} else {
		_, err := db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		return err
	}
}

func DeleteTask(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
