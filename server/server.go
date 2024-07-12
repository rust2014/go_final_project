package server

import (
	"github.com/rust2014/go_final_project/services"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rust2014/go_final_project/database"
	"github.com/rust2014/go_final_project/handlers"

	_ "modernc.org/sqlite"
)

func Run() {
	db, err := database.ConnectDB() // запуск бд
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer db.Close() // закрытие бд

	port := os.Getenv("TODO_PORT") // если переменная окружения TODO_PORT не установлена, сервер будет запущен на порту 7540
	if port == "" {
		port = "7540" // порт по умолчанию
	}
	webDir := "./web" // каталог с вебом
	fileServer := http.FileServer(http.Dir(webDir))

	router := chi.NewRouter()
	router.Handle("/*", fileServer) // обработчик файлов

	router.Get("/api/nextdate", services.NextDateHandler) // Правила повторения задач, обработчик для вычисления следующей даты (3)

	router.Post("/api/task", handlers.HandlerTask(db))          // добавляем задачу в бд - AddTask (4)
	router.Get("/api/task", handlers.HandlerGetTask(db))        // просмотр задачи (6)
	router.Put("/api/task", handlers.HandlerPutTask(db))        // редактирование задачи (6)
	router.Post("/api/task/done", handlers.HandlerDoneTask(db)) // завершение задачи (7)
	router.Delete("/api/task", handlers.HandlerDeleteTask(db))  // удаление задачи (7)

	router.Get("/api/tasks", handlers.HandlerGetTasks(db)) // Получаем список ближайших задач в вебе (5)

	log.Printf("Starting server at port %s", port) // сообщение о старте + порт
	err = http.ListenAndServe(":"+port, router)    // запуск сервера на нашем порте из переменной port
	if err != nil {
		log.Fatal(err) // для логирования ошибки
	}
}
