package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rust2014/go_final_project/dates"
	"github.com/rust2014/go_final_project/models"
	"github.com/rust2014/go_final_project/services"
	"github.com/rust2014/go_final_project/validation"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) { // GET-обработчик api/nextdate
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	timeNow, err := time.Parse(dates.DefaultDateFormat, now)
	if err != nil {
		http.Error(w, "invalid 'now' date format", http.StatusBadRequest) // вернет Invalid 'now' date format если не верный формат
		return
	}
	nextDate, err := dates.NextDate(timeNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // вернет "key" is required
		return
	}
	fmt.Fprintf(w, nextDate)
}

func HandlerTask(taskService *services.TaskService) http.HandlerFunc { // обработчик для AddTask
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error": "JSON deserialization error:"}`, http.StatusBadRequest)
			return
		}
		id, err := taskService.AddTask(task)
		if err != nil {
			http.Error(w, `{"error": "Error when adding a task:"}`, http.StatusBadRequest)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{"id": id})
	}
}

func HandlerGetTasks(taskService *services.TaskService) http.HandlerFunc { // обработчик для GET-запроса /api/tasks
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := taskService.GetTasks()
		if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{"tasks": tasks})
	}
}

func HandlerGetTask(taskService *services.TaskService) http.HandlerFunc { // обработчик GET-запроса /api/task?id=
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
		task, err := taskService.GetTask(id)
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

func HandlerPutTask(taskService *services.TaskService) http.HandlerFunc { //обработчик PUT-запроса /api/task (проверка как в HandlerTask)
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

		if !validation.IsValidJSON(task.Title) || !validation.IsValidJSON(task.Comment) {
			http.Error(w, `{"error": "Incorrect characters in the title or comments"}`, http.StatusBadRequest)
			return
		}

		if _, err := strconv.Atoi(task.ID); err != nil {
			http.Error(w, `{"error": "Incorrect identifier format"}`, http.StatusBadRequest)
			return
		}

		if _, err := time.Parse(dates.DefaultDateFormat, task.Date); err != nil {
			http.Error(w, `{"error": "Incorrect date format"}`, http.StatusBadRequest)
			return
		}

		if err := validation.ValidateRepeatRule(task.Repeat); err != nil {
			http.Error(w, `{"error": "Incorrect repeat format"}`, http.StatusBadRequest)
			return
		}

		err := taskService.UpdateTask(task)
		if err != nil {
			if err.Error() == "task not found" {
				http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}

		writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
	}
}

func HandlerDoneTask(taskService *services.TaskService) http.HandlerFunc { // обработчик POST-запроса /api/task/done
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
		task, err := taskService.GetTask(id)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error": "Request execution error"}`, http.StatusInternalServerError)
			return
		}
		nextDate := ""
		if task.Repeat != "" {
			nextDate, err = dates.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Error calculating the next date"}`, http.StatusInternalServerError)
				return
			}
		}
		if err := taskService.DoneTask(id, nextDate); err != nil {
			http.Error(w, `{"error": "Task update error"}`, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
	}
}

func HandlerDeleteTask(taskService *services.TaskService) http.HandlerFunc { // обработчик Delete запроса /api/task
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
		err = taskService.DeleteTask(id)
		if err != nil {
			if err.Error() == "task not found" {
				http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error": "Task deletion error"}`, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
	}
}
