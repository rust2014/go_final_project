package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rust2014/go_final_project/models"
	"github.com/rust2014/go_final_project/services"
)

func HandlerTask(db *sql.DB) http.HandlerFunc { // обработчик для AddTask
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error": "JSON deserialization error:"}`, http.StatusBadRequest)
			return
		}
		id, err := services.AddTask(db, task)
		if err != nil {
			http.Error(w, `{"error": "Error when adding a task:"}`, http.StatusBadRequest) // важная ошибка для тестов
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{"id": id})
	}
}

func HandlerGetTasks(db *sql.DB) http.HandlerFunc { // обработчик для GET-запроса /api/tasks
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := services.GetTasks(db)
		if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{"tasks": tasks})
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
		task, err := services.GetTask(db, id)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, task)
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

		if !services.IsValidJSON(task.Title) || !services.IsValidJSON(task.Comment) {
			http.Error(w, `{"error": "Incorrect characters in the title or comments"}`, http.StatusBadRequest)
			return
		}

		if _, err := strconv.Atoi(task.ID); err != nil {
			http.Error(w, `{"error": "Incorrect identifier format"}`, http.StatusBadRequest)
			return
		}

		if _, err := time.Parse(services.DefaultFormatDate, task.Date); err != nil {
			http.Error(w, `{"error": "Incorrect date format"}`, http.StatusBadRequest)
			return
		}

		if err := services.ValidateRepeatRule(task.Repeat); err != nil {
			http.Error(w, `{"error": "Incorrect repeat format"}`, http.StatusBadRequest)
			return
		}

		rowsAffected, err := services.UpdateTask(db, task)
		if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
	}
}

func HandlerDoneTask(db *sql.DB) http.HandlerFunc { // обработчик POST-запроса /api/task/done
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
		task, err := services.GetTask(db, id)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}
		nextDate := ""
		if task.Repeat != "" {
			nextDate, err = services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Error calculating the next date"}`, http.StatusInternalServerError)
				return
			}
		}
		if err := services.DoneTask(db, id, nextDate); err != nil {
			http.Error(w, `{"error": "Task update error"}`, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
	}
}

func HandlerDeleteTask(db *sql.DB) http.HandlerFunc { // обработчик Delete запроса /api/task
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
		rowsAffected, err := services.DeleteTask(db, id)
		if err != nil {
			http.Error(w, `{"error": "Task deletion error"}`, http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
	}
}
