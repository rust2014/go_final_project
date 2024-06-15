package server

import (
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

	http.Handle("/", http.FileServer(http.Dir(webDir))) // обработчик

	err := http.ListenAndServe(":"+port, nil) // запуск сервера на нашем порте из переменной port
	if err != nil {
		log.Fatal(err) // для логирования ошибки
	}
}
