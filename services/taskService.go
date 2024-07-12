package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rust2014/go_final_project/models"
	"time"
)

func AddTask(db *sql.DB, task models.Task) (int64, error) { // создание новой задачи в бд
	// Логирование входных данных для диагностики ошибок
	fmt.Printf("Task Date: %s\n", task.Date)
	fmt.Printf("Task Title: %s\n", task.Title)
	fmt.Printf("Task Comment: %s\n", task.Comment)
	fmt.Printf("Task Repeat: %s\n", task.Repeat)

	if !IsValidJSON(task.Title) || !IsValidJSON(task.Comment) {
		return 0, errors.New("incorrect characters in the title or comments")
	}

	if task.Title == "" {
		return 0, errors.New("no task title")
	}

	now := time.Now()
	if task.Date == "" || task.Date == "today" {
		task.Date = now.Format(DefaultFormatDate)
	} else {
		date, err := time.Parse(DefaultFormatDate, task.Date)
		if err != nil {
			return 0, errors.New("the date is in the wrong format")
		}
		if date.Before(now) && date.Format(DefaultFormatDate) != now.Format(DefaultFormatDate) { // пересчитываем дату, только если она в прошлом
			if task.Repeat == "" {
				task.Date = now.Format(DefaultFormatDate)
			} else {
				if task.Repeat == "d 1" {
					task.Date = now.Format(DefaultFormatDate)
				} else {
					nextDate, err := NextDate(now, task.Date, task.Repeat)
					if err != nil {
						return 0, err
					}
					task.Date = nextDate
				}
			}
		}
	}

	if err := ValidateRepeatRule(task.Repeat); err != nil {
		fmt.Println(err)
		return 0, err
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
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
