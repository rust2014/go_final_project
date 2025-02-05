package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/rust2014/go_final_project/dates"
	"github.com/rust2014/go_final_project/models"
	"github.com/rust2014/go_final_project/validation"
)

func (s *TaskService) AddTask(task models.Task) (int64, error) { // создание новой задачи в бд
	// Логирование входных данных для диагностики ошибок
	fmt.Printf("Task Date: %s\n", task.Date)
	fmt.Printf("Task Title: %s\n", task.Title)
	fmt.Printf("Task Comment: %s\n", task.Comment)
	fmt.Printf("Task Repeat: %s\n", task.Repeat)

	if !validation.IsValidJSON(task.Title) || !validation.IsValidJSON(task.Comment) {
		return 0, errors.New("incorrect characters in the title or comments")
	}

	if task.Title == "" {
		return 0, errors.New("no task title")
	}

	now := time.Now()
	if task.Date == "" || task.Date == "today" {
		task.Date = now.Format(dates.DefaultDateFormat)
	} else {
		date, err := time.Parse(dates.DefaultDateFormat, task.Date)
		if err != nil {
			return 0, errors.New("the date is in the wrong format")
		}
		if date.Before(now) && date.Format(dates.DefaultDateFormat) != now.Format(dates.DefaultDateFormat) { // пересчитываем дату, только если она в прошлом
			if task.Repeat == "" {
				task.Date = now.Format(dates.DefaultDateFormat)
			} else {
				if task.Repeat == "d 1" {
					task.Date = now.Format(dates.DefaultDateFormat)
				} else {
					nextDate, err := dates.NextDate(now, task.Date, task.Repeat)
					if err != nil {
						return 0, err
					}
					task.Date = nextDate
				}
			}
		}
	}

	if err := validation.ValidateRepeatRule(task.Repeat); err != nil {
		fmt.Println(err)
		return 0, err
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := s.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return 0, err
	}
	id, err := res.LastInsertId() // Получение ID
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		return 0, err
	}
	fmt.Println("Inserted task with ID:", id)
	return id, nil
}
