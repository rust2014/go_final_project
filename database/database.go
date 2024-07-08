package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite" // импорт драйвера SQLite
)

func ConnectDB() (*sql.DB, error) {
	var nameDB = "scheduler.db"
	path := os.Getenv("TODO_DBFILE")
	if path == "" {
		path = nameDB
	}

	var install bool
	if _, err := os.Stat(path); os.IsNotExist(err) { // проверка существования файла базы данных
		install = true                     // если файл базы данных не существует (true), то создаем файл через sql запросы из createTableAndIndex
		fileCreate, err := os.Create(path) // создается пустой файл бд
		if err != nil {
			return nil, err
		}
		fileCreate.Close()
	}

	db, err := sql.Open("sqlite", path) // подключение к БД (не работает с sqlite3)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	if install { // создание таблицы и индекса если надо
		createTableAndIndex(db)
	}
	return db, nil
}

func createTableAndIndex(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT '',
		title VARCHAR(256) NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat VARCHAR(128) NOT NULL DEFAULT ''
	);
	CREATE INDEX idx_date ON scheduler(date);
`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы или индекса: %v", err)
	} else {
		log.Println("Таблица и индекс успешно созданы.")
	}
}
