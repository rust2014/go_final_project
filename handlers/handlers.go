package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rust2014/go_final_project/models"
	"net/http"
	"strconv"
	"time"
)

func HandlerTask(db *sql.DB) http.HandlerFunc { // обработчик для AddTask
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			fmt.Println("JSON deserialization error: ", err) // ошибка в консоль
			http.Error(w, `{"error": "JSON deserialization error:"}`, http.StatusBadRequest)
			return
		}
		id, err := AddTask(db, task)
		if err != nil {
			http.Error(w, `{"error": "Error when adding a task:"}`, http.StatusBadRequest) // важная ошибка для тестов
			return
		}

		w.Header().Set("Content-Type", "application/json, charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"id": id})
	}
}

func HandlerGetTasks(db *sql.DB) http.HandlerFunc { // обработчик для GET-запроса /api/tasks
	return func(w http.ResponseWriter, r *http.Request) {
		tasks := []models.Task{} // пустой слайс

		rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50") // SQL-запроса для получения задач (лимит50)
		if err != nil {
			fmt.Println("Request execution error: ", err)
			http.Error(w, `{"error": "Request execution error:"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() { // обход строк
			var task models.Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				fmt.Println("Read error: ", err)
				http.Error(w, `{"error": "Read error:"}`, http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}
		if err := rows.Err(); err != nil { // проверка обхода
			fmt.Println("Results processing error: ", err)
			http.Error(w, `{"error": "Results processing error:"}`, http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{ // ответ в формате JSON
			"tasks": tasks,
		}
		w.Header().Set("Content-Type", "application/json, charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Println("Error json: ", err)
			http.Error(w, `{"error": "Error json:"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandlerGetTask(db *sql.DB) http.HandlerFunc { // обработчик GET-запроса /api/task?id=
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error": "No identifier specified"}`, http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Incorrect identifier format"}`, http.StatusBadRequest)
			return
		}
		var task models.Task
		err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{ // Формирование ответа в формате JSON
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}
		w.Header().Set("Content-Type", "application/json, charset=UTF-8")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error": "Response encoding error"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandlerPutTask(db *sql.DB) http.HandlerFunc { //обработчик PUT-запроса /api/task (проверка как в HandlerTask)
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error": "Incorrect data format"}`, http.StatusBadRequest)
			return
		}

		if task.ID == "" || task.Date == "" || task.Title == "" {
			http.Error(w, `{"error": "No mandatory fields"}`, http.StatusBadRequest)
			return
		}

		if !IsValidJSON(task.Title) || !IsValidJSON(task.Comment) {
			http.Error(w, `{"error": "Incorrect characters in the title or comments"}`, http.StatusBadRequest)
			return
		}

		if _, err := strconv.Atoi(task.ID); err != nil {
			http.Error(w, `{"error": "Incorrect identifier format"}`, http.StatusBadRequest)
			return
		}

		if _, err := time.Parse(DefaultFormatDate, task.Date); err != nil {
			http.Error(w, `{"error": "Incorrect date format"}`, http.StatusBadRequest)
			return
		}

		if err := ValidateRepeatRule(task.Repeat); err != nil {
			http.Error(w, `{"error": "Incorrect repeat format"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
			task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, `{"error": "Error receiving a result"}`, http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json, charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{}); err != nil {
			http.Error(w, `{"error": "Response encoding error"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandlerDoneTask(db *sql.DB) http.HandlerFunc { // обработчик POST-запроса /api/task/done
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error": "No identifier specified"}`, http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Incorrect identifier format"}`, http.StatusBadRequest)
			return
		}
		err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}
		if task.Repeat == "" {
			_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", task.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) // заменить на `{"error": "Ошибка удаления задачи"}`
				return
			}
		} else {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Error calculating the next date"}`, http.StatusInternalServerError)
				return
			}
			_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, task.ID)
			if err != nil {
				fmt.Println("Task update error", err)
				http.Error(w, `{"error": "Task update error"}`, http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json, charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{}); err != nil {
			http.Error(w, `{"error": "Response encoding error"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandlerDeleteTask(db *sql.DB) http.HandlerFunc { // обработчик Delete запроса /api/task
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error": "No identifier specified"}`, http.StatusBadRequest)
			return
		}
		if _, err := strconv.Atoi(idStr); err != nil {
			http.Error(w, `{"error": "Incorrect identifier format"}`, http.StatusBadRequest)
			return
		}
		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", idStr)
		if err != nil {
			fmt.Println("Task deletion error", err)
			http.Error(w, `{"error": "Task deletion error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json, charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{}); err != nil {
			http.Error(w, `{"error": "Response encoding error"}`, http.StatusInternalServerError)
			return
		}
	}
}
