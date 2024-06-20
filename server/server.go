package server

import (
	"github.com/rust2014/go_final_project/handlers"
	"log"
	"net/http"
	"os"
)

func Run() {
	port := os.Getenv("TODO_PORT") // если переменная окружения TODO_PORT не установлена, сервер будет запущен на порту 7540
	if port == "" {
		port = "7540" // порт по умолчанию
	}
	//port := "7540"
	webDir := "./web"                              // каталог с вебом
	log.Printf("Starting server at port %s", port) // сообщение о старте + порт

	http.Handle("/", http.FileServer(http.Dir(webDir))) // обработчик файлов

	http.HandleFunc("/api/nextdate", handlers.NextDateHandler) // обработчик для вычисления следующей даты
	// пример запроса http://localhost:7540/api/nextdate?now=20240126&date=20240126&repeat=d%201
	// repeat is required -  http://localhost:7540/api/nextdate?now=20240126&date=20240126&repeat=
	//http.HandleFunc("/api/tasks", handlers.getAllTasks)       // обрабочик для получения всех задач
	//http.HandleFunc("/api/tasks", handlers.createTask)        // создание задач
	//http.HandleFunc("/api/tasks/{id}", handlers.getTaskID)    // обработчик дkя задачи по ID
	//http.HandleFunc("/api/tasks/{id}", handlers.deleteTaskID) // обработчик для удаления задачи

	err := http.ListenAndServe(":"+port, nil) // запуск сервера на нашем порте из переменной port
	if err != nil {
		log.Fatal(err) // для логирования ошибки
	}
}
