package services

import (
	"database/sql"
	"fmt"

	"github.com/rust2014/go_final_project/models"
	"github.com/rust2014/go_final_project/tests"
)

type TaskService struct {
	DB *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{DB: db}
}

func (s *TaskService) GetTasks() ([]models.Task, error) {
	tasks := []models.Task{}
	query := fmt.Sprintf("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT %d", tests.TaskLimit)
	rows, err := s.DB.Query(query)
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

func (s *TaskService) GetTask(id int) (*models.Task, error) {
	var task models.Task
	err := s.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) UpdateTask(task models.Task) error {
	result, err := s.DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (s *TaskService) DoneTask(id int, nextDate string) error {
	if nextDate == "" {
		_, err := s.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
		return err
	} else {
		_, err := s.DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		return err
	}
}

func (s *TaskService) DeleteTask(id int) error {
	result, err := s.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
